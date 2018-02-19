package utility

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
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
	URL                 string
	ParsedURL           *url.URL
	RequestMethod       RequestMethod
	RequestMethodString string
	Headers             map[string]string
	Body                string
}

type HTTPClient struct {
	retry     int
	retryWait int
	timeout   int
}

func (c *HTTPClient) newHTTPTransport() (transport *http.Transport) {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   300 * time.Second,
			KeepAlive: 300 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   300 * time.Second,
		ExpectContinueTimeout: 300 * time.Second,
	}
}

func (c *HTTPClient) newHTTPClient(scheme string, host string, timeout int) *http.Client {
	transport := c.newHTTPTransport()
	if scheme == "https" {
		transport.TLSClientConfig = &tls.Config{ServerName: host}
	}
	return &http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeout) * time.Second,
	}
}

func (c *HTTPClient) methodFuncBase(method string, request *HTTPRequest) (*http.Response, []byte, error) {
	httpClient := c.newHTTPClient(request.ParsedURL.Scheme, request.ParsedURL.Host, c.timeout)
	req, err := http.NewRequest(method, request.URL, bytes.NewBufferString(request.Body))
	if err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("can not create request (url = %v, method = %v, request body = %v,)", request.URL, request.RequestMethod, request.Body))
	}
	for k, v := range request.Headers {
		req.Header.Set(k, v)
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return res, nil, errors.Wrap(err, fmt.Sprintf("request of HTTPMethodGET is failure (url = %v, method = %v, request body = %v)", request.URL, request.RequestMethod, request.Body))
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res, resBody, errors.Wrap(err, fmt.Sprintf("can not read response (url = %v, method = %v, request body = %v)", request.URL, request.RequestMethod, request.Body))
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return res, resBody, errors.Errorf("unexpected status code (url = %v, method = %v, request body = %v, status = %v, body = %v)", request.URL, request.RequestMethod, request.Body, res.StatusCode, string(resBody))
	}
	// log.Printf("http ok (url = %v, method = %v, request body = %v, status = %v, response body = %v)", request.URL, request.RequestMethod , request.Body, res.StatusCode, string(resBody))
	return res, resBody, nil
}

func (c *HTTPClient) retryRequest(request *HTTPRequest, noRetry bool) (*http.Response, []byte, error) {
	for i := 0; i <= c.retry; i++ {
		res, resBody, err := c.methodFuncBase(request.RequestMethodString, request)
		if !noRetry && err != nil {
			log.Printf("request is failure, retry... (url = %v, method = %v, reason = %v)", request.URL, request.RequestMethod, err)
			if (c.retryWait != 0) {
				time.Sleep(time.Duration(c.retryWait) * time.Millisecond)
			}
			continue
		}
		return res, resBody, err
	}
	return nil, nil, errors.Errorf("give up retry (url = %v, method = %v)", request.URL, request.RequestMethod)
}

func (c *HTTPClient) DoRequest(requestMethod RequestMethod, request *HTTPRequest, noRetry bool) (*http.Response, []byte, error) {
	u, err := url.Parse(request.URL)
	if err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("can not parse url (url = %v, method = %v)", request.URL, requestMethod))
	}
	request.ParsedURL = u
	request.RequestMethod = requestMethod
	switch requestMethod {
	case HTTPMethodHEAD:
		request.RequestMethodString = "HEAD"
	case HTTPMethodGET:
		request.RequestMethodString = "GET"
	case HTTPMethodPUT:
		request.RequestMethodString = "PUT"
	case HTTPMethdoPOST:
		request.RequestMethodString = "POST"
	case HTTPMethodDELETE:
		request.RequestMethodString = "DELETE"
	case HTTPMethodPATCH:
		request.RequestMethodString = "PATCH"
	default:
		return nil, nil, errors.Wrap(err, fmt.Sprintf("unsupported request method (url = %v, method = %v)", request.URL, requestMethod))
	}
	res, resBody, err := c.retryRequest(request, noRetry)
	if err != nil {
		return res, resBody, errors.Wrap(err, fmt.Sprintf("request is failure (url = %v, method = %v)", request.URL, request.RequestMethod))
	}
	return res, resBody, err
}

func NewHTTPClient(retry int, retryWait int, timeout int) *HTTPClient {
	if retry == 0 {
		retry = 200000
	}
	if timeout == 0 {
		timeout = 300
	}
	return &HTTPClient{
		retry:     retry,
		retryWait: retryWait,
		timeout:   timeout,
	}
}

type WSCallback func(conn *websocket.Conn, data interface{}) error

type WSClient struct {
	readBufSize  int
	writeBufSize int
	retry        int
	retryWait    int
	conn         *websocket.Conn
	connChan     chan error
	started      bool
	finished     uint32
}

func (w *WSClient) messageLoop(callback WSCallback, callbackData interface{}) (error) {
	for {
		err := callback(w.conn, callbackData)
		if err != nil {
			return errors.Wrapf(err, "callback error (reason = %v)", err)
		}
		if atomic.LoadUint32(&w.finished) == 1 {
			log.Printf("message loop finished")
			return nil
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

func (w *WSClient) startPing() chan bool {
	pingChan := make(chan bool)
	go w.pingLoop(pingChan)
	return pingChan
}

func (w *WSClient) stopPing(pingChan chan bool) {
	close(pingChan)
}

func (w *WSClient) connect(callback WSCallback, callbackData interface{}, requestURL string, requestHeaders map[string]string) {
	u, err := url.Parse(requestURL)
	if err != nil {
		w.connChan <- errors.Wrap(err, fmt.Sprintf("can not parse url (url = %v, header = %v)", requestURL, requestHeaders))
		return
	}
	header := http.Header{}
	for k, v := range requestHeaders {
		header.Set(k, v)
	}
	dialer := &websocket.Dialer{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{ServerName: u.Host},
	}
	i := 0
	for {
		i++
		if w.retry != 0 {
			if i > w.retry {
				break
			}
		}
		conn, response, err := dialer.Dial(requestURL, header)
		if err != nil {
			if response == nil {
				log.Printf("can not dial, retry ... (url = %v, header = %v, reason = %v)", requestURL, requestHeaders, err)
				time.Sleep(500 * time.Millisecond)
				continue
			}
			if response.StatusCode < 200 && response.StatusCode <= 300 {
				log.Printf("can not dial, retry ... (url = %v, header = %v, reason = %v)", requestURL, requestHeaders, err)
				time.Sleep(500 * time.Millisecond)
				continue
			}
			w.connChan <- errors.Wrap(err, fmt.Sprintf("can not dial (URL = %v, header = %v)", requestURL, requestHeaders))
			continue
		}
		if !w.started {
			w.connChan <- nil
			w.started = true
		}
		w.conn = conn
		pingChan := w.startPing()
		err = w.messageLoop(callback, callbackData)
		if err != nil {
			log.Printf("error occuered in message loop: %v", err)
		}
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

func (w *WSClient) Start(callback WSCallback, callbackData interface{}, requestURL string, requestHeaders map[string]string) error {
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

// Send ...
// TODO: cleanup
func (w *WSClient) Send(v interface{}) error {
	return w.conn.WriteJSON(v)
}

func NewWSClient(readBufSize int, writeBufSize int, retry int, retryWait int) *WSClient {
	if readBufSize == 0 {
		readBufSize = 1024 * 1024 * 2
	}
	if writeBufSize == 0 {
		writeBufSize = 1024 * 1024 * 2
	}
	return &WSClient{
		readBufSize:  readBufSize,
		writeBufSize: writeBufSize,
		retry:        retry,
		retryWait:    retryWait,
		connChan:     make(chan error),
	}
}
