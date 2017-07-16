package utility

import (
	"time"
	"github.com/pkg/errors"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"fmt"
	"io/ioutil"
	"net"
	"crypto/tls"
	"log"
	"bytes"
	"sync/atomic"
)

type RequestMethod int

const (
	HTTPMethodHEAD   RequestMethod = iota
	HTTPMethodGET
	HTTPMethodPUT
	HTTPMethdoPOST
	HTTPMethodDELETE
	HTTPMethodPATCH
)

// HTTPRequest is http request
type HTTPRequest struct {
	URL		      string
	ParsedURL     *url.URL
	RequestMethod RequestMethod
	Headers  	  map[string]string
	Body     	  string
}

type HTTPClient struct {
	retry 	int
	timeout int
}

type httpMethodFunc func(request *HTTPRequest) (*http.Response, []byte, bool, error)

func (c *HTTPClient) newHTTPTransport() (transport *http.Transport) {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   60 * time.Second,
			KeepAlive: 60 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   30 * time.Second,
		ExpectContinueTimeout: 30 * time.Second,
	}
}

func (c *HTTPClient) newHTTPClient(scheme string, host string, timeout int) (*http.Client) {
	transport := c.newHTTPTransport()
	if scheme == "https" {
		transport.TLSClientConfig = &tls.Config{ServerName: host}
	}
	return &http.Client{
		Transport: transport,
		Timeout: time.Duration(timeout) * time.Second,
	}
}

func (c *HTTPClient) methodFuncBase(method string, request *HTTPRequest) (*http.Response, []byte, bool, error) {
	httpClient := c.newHTTPClient(request.ParsedURL.Scheme, request.ParsedURL.Host, c.timeout)
	req, err := http.NewRequest(method, request.URL, bytes.NewBufferString(request.Body))
	if err != nil {
		return nil, nil, false, errors.Wrap(err, fmt.Sprintf("can not create request (url = %v, method = %v, request body = %v,)", request.URL, request.RequestMethod, request.Body))
	}
	for k,v := range request.Headers {
		req.Header.Set(k, v)
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return res, nil, true, errors.Wrap(err, fmt.Sprintf("request of HTTPMethodGET is failure (url = %v, method = %v, request body = %v)", request.URL, request.RequestMethod, request.Body))
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res, resBody, true, errors.Wrap(err, fmt.Sprintf("can not read response (url = %v, method = %v, request body = %v)", request.URL, request.RequestMethod, request.Body))
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		retryable := false
		if res.StatusCode >= 500 && res.StatusCode < 600 {
			retryable = true
		}
		return res, resBody, retryable, errors.Errorf("unexpected status code (url = %v, method = %v, request body = %v, status = %v, body = %v)", request.URL, request.RequestMethod, request.Body, res.StatusCode, string(resBody))
	}
	// log.Printf("http ok (url = %v, method = %v, request body = %v, status = %v, response body = %v)", request.URL, request.RequestMethod , request.Body, res.StatusCode, string(resBody))
	return res, resBody, false, nil
}

func (c *HTTPClient) retryRequest(methodFunc httpMethodFunc, request *HTTPRequest) (*http.Response, []byte, error)  {
	for i := 0; i <= c.retry; i++ {
		res, resBody, retryable, err := methodFunc(request)
		if retryable {
			log.Printf("request is failure, retry... (url = %v, method = %v, reason = %v)", request.URL, request.RequestMethod, err)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		return res, resBody, err
	}
	return nil, nil, errors.Errorf("give up retry (url = %v, method = %v)", request.URL, request.RequestMethod)
}

func (c *HTTPClient) DoRequest(requestMethod RequestMethod, request *HTTPRequest) (*http.Response, []byte, error) {
	u, err :=  url.Parse(request.URL)
	if err != nil {
		return  nil, nil, errors.Wrap(err,fmt.Sprintf("can not parse url (url = %v, method = %v)", request.URL, requestMethod))
	}
	request.ParsedURL = u
	var methodFunc httpMethodFunc
	switch requestMethod {
	case HTTPMethodHEAD:
		methodFunc = func (request *HTTPRequest) (*http.Response, []byte, bool, error) {
			return c.methodFuncBase("HEAD", request)
		}
	case HTTPMethodGET:
		methodFunc = func (request *HTTPRequest) (*http.Response, []byte, bool, error) {
			return c.methodFuncBase("GET", request)
		}
	case HTTPMethodPUT:
		methodFunc = func (request *HTTPRequest) (*http.Response, []byte, bool, error) {
			return c.methodFuncBase("PUT", request)
		}
	case HTTPMethdoPOST:
		methodFunc = func (request *HTTPRequest) (*http.Response, []byte, bool, error) {
			return c.methodFuncBase("POST", request)
		}
	case HTTPMethodDELETE:
		methodFunc = func (request *HTTPRequest) (*http.Response, []byte, bool, error) {
			return c.methodFuncBase("DELETE", request)
		}
	case HTTPMethodPATCH:
		methodFunc = func (request *HTTPRequest) (*http.Response, []byte, bool, error) {
			return c.methodFuncBase("PATCH", request)
		}
	default:
		return nil, nil, errors.Wrap(err, fmt.Sprintf("unsupported request method (url = %v, method = %v)", request.URL, requestMethod))
	}
	request.RequestMethod = requestMethod
	res, resBody, err := c.retryRequest(methodFunc, request)
	if err != nil {
		return res, resBody, errors.Wrap(err, fmt.Sprintf("request is failure (url = %v, method = %v)", request.URL, request.RequestMethod))
	}
	return res, resBody, err
}

func NewHTTPClient(retry int, timeout int) (*HTTPClient) {
	if retry == 0 {
		retry = 60
	}
	if timeout == 0 {
		timeout = 60
	}
	return &HTTPClient{
		retry: retry,
		timeout: timeout,
	}
}

type WSCallback func(conn *websocket.Conn, data interface{}) (error)

type WSClient struct {
	readBufSize 	int
	writeBufSize	int
	retry       	int
	conn          	*websocket.Conn
	connChan  		chan error
	started			bool
	finished        uint32
}

func (w *WSClient) messageLoop(callback WSCallback, callbackData interface{}) {
	for {
		err := callback(w.conn, callbackData)
		if err != nil {
			log.Printf("callback error (reason = %v)", err)
		}
		if atomic.LoadUint32(&w.finished) == 1{
			log.Printf("message loop finished")
			return
		}
	}
}

func (w *WSClient) pingLoop(pingChan chan bool) {
	for {
		select {
		case <-pingChan:
			return
		case <-time.After(5 * time.Second):
			deadline := time.Now()
			deadline.Add(30 * time.Second)
			w.conn.WriteControl(websocket.PingMessage, []byte("ping"), deadline)
		}
	}
}

func (w *WSClient) startPing() (chan bool) {
	pingChan := make(chan bool)
	go w.pingLoop(pingChan)
	return pingChan
}

func (w *WSClient) stopPing(pingChan chan bool) {
	close(pingChan)
}

func (w *WSClient) connect(callback WSCallback, callbackData interface{}, requestURL string, requestHeaders map[string]string) {
	u, err :=  url.Parse(requestURL)
	if err != nil {
		w.connChan <- errors.Wrap(err,fmt.Sprintf("can not parse url (url = %v, header = %v)", requestURL, requestHeaders))
		return
	}
	header := http.Header{}
	for k, v := range requestHeaders {
		header.Set(k, v)
	}
	dialer := &websocket.Dialer{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{ServerName: u.Host},
	}
	i := 0
	for {
		i++
		if i > w.retry {
			break
		}
		conn, response, err := dialer.Dial(requestURL, header)
		if err != nil {
			if response == nil {
				log.Printf("can not dial, retry ... (url = %v, header = %v, reason = %v)", requestURL, requestHeaders, err)
				time.Sleep(500 * time.Millisecond)
				continue
			}
			if response.StatusCode >= 500 && response.StatusCode < 600 {
				log.Printf("can not dial, retry ... (url = %v, header = %v, reason = %v)", requestURL, requestHeaders, err)
				time.Sleep(500 * time.Millisecond)
				continue
			}
			w.connChan <- errors.Wrap(err, fmt.Sprintf("can not dial (URL = %v, header = %v)", requestURL, requestHeaders))
			return
		}
		if !w.started {
			w.connChan <- nil
			w.started = true
		}
		w.conn = conn
		pingChan := w.startPing()
		w.messageLoop(callback, callbackData)
		w.stopPing(pingChan)
		conn.Close()
		i = 0
		if atomic.LoadUint32(&w.finished) == 1 {
			log.Printf("connect loop finished")
			return
		}
	}
	w.connChan <- errors.Errorf("give up retry (url = %v, header = %v)", requestURL, requestHeaders)
}

func (w *WSClient) Start(callback WSCallback, callbackData interface{}, requestURL string, requestHeaders map[string]string) (error){
	atomic.StoreUint32(&w.finished, 0)
	go w.connect(callback, callbackData, requestURL, requestHeaders)
	select {
	case err := <-w.connChan:
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("can not connect (URL = %v, header = %v)", requestURL, requestHeaders))
		}
	}
	close(w.connChan)
	return nil
}

func (w *WSClient) Stop() {
	atomic.StoreUint32(&w.finished, 1)
}

func NewWSClient(readBufSize int, writeBufSize int, retry int) (*WSClient) {
	if readBufSize == 0 {
		readBufSize = 1024 * 1024 * 2
	}
	if writeBufSize == 0 {
		writeBufSize = 1024 * 1024 * 2
	}
	if retry == 0 {
		retry = 60
	}
	return &WSClient{
		readBufSize: readBufSize,
		writeBufSize: writeBufSize,
		retry: retry,
		connChan: make(chan error),
	}
}
