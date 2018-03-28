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
	"github.com/viki-org/dnscache"
	"strings"
	"sync"
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

type client struct {
	httpClient    *http.Client
	requestMutex  *sync.Mutex
}

type clientCache struct {
	client    *client
	tlsClient *client
}

type HTTPClient struct {
	retry             int
	retryWait         int
	timeout           int
	localAddr         *net.TCPAddr
	resolver          *dnscache.Resolver
	resolverIdx       int
	resolverIdxMutex  *sync.Mutex
	clientsCache      map[string]*clientCache
	clientsCacheMutex *sync.Mutex
}

func (c *HTTPClient) newHTTPTransport(localAddr *net.TCPAddr) (transport *http.Transport) {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: func(network string, address string) (net.Conn, error) {
			ipv6 := false
			if strings.LastIndex(address, ".") == -1 {
				ipv6 = true
			}
			separator := strings.LastIndex(address, ":")
			fmt.Println(address[:separator])
			ips, _ := c.resolver.Fetch(address[:separator])
			c.resolverIdxMutex.Lock()
			c.resolverIdx += 1
			if len(ips) <= c.resolverIdx {
				c.resolverIdx = 0
			}
			resolverIds := c.resolverIdx
			c.resolverIdxMutex.Unlock()
			ip := ips[resolverIds]
			ipStr := ip.String()
			if strings.LastIndex(ipStr, ".") == -1 {
				ipv6 = true
			}
			if ipv6 {
				ipStr = "[" + ipStr + "]"
			}
			return (&net.Dialer{
				LocalAddr: localAddr,
				Timeout:   300 * time.Second,
				KeepAlive: 300 * time.Second,
			}).Dial("tcp", ipStr+address[separator:])
		},
		TLSHandshakeTimeout:   300 * time.Second,
		ExpectContinueTimeout: 300 * time.Second,
		MaxIdleConns:          1000,
		MaxIdleConnsPerHost:   50,
	}
}

func (c *HTTPClient) newClient(scheme string, host string, timeout int) *client {
	c.clientsCacheMutex.Lock()
	defer c.clientsCacheMutex.Unlock()
	clients, ok := c.clientsCache[host]
	if ok {
		if scheme == "https" {
			if clients.tlsClient != nil {
				return clients.tlsClient
			}
		} else {
			if clients.client != nil {
				return clients.client
			}
		}
	} else {
		c.clientsCache[host] = &clientCache{}
	}
	transport := c.newHTTPTransport(c.localAddr)
	if scheme == "https" {
		transport.TLSClientConfig = &tls.Config{ServerName: host}
	}
	newHttpClient := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeout) * time.Second,
	}
	newClinet := &client{
		httpClient: newHttpClient,
		requestMutex: new(sync.Mutex),
	}
	if scheme == "https" {
		c.clientsCache[host].tlsClient = newClinet
	} else {
		c.clientsCache[host].client = newClinet
	}
	return newClinet
}


func (c *HTTPClient) methodFuncBase(method string, request *HTTPRequest) (*http.Response, []byte, error) {
	client := c.newClient(request.ParsedURL.Scheme, request.ParsedURL.Host, c.timeout)
	req, err := http.NewRequest(method, request.URL, bytes.NewBufferString(request.Body))
	if err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("can not create request (url = %v, method = %v, request body = %v,)", request.URL, request.RequestMethod, request.Body))
	}
	for k, v := range request.Headers {
		req.Header.Set(k, v)
	}
	client.requestMutex.Lock()
	res, err := client.httpClient.Do(req)
	if err != nil {
		client.requestMutex.Unlock()
		return res, nil, errors.Wrap(err, fmt.Sprintf("request of HTTPMethodGET is failure (url = %v, method = %v, request body = %v)", request.URL, request.RequestMethod, request.Body))
	}
	client.requestMutex.Unlock()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		res.Body.Close()
		return res, resBody, errors.Wrap(err, fmt.Sprintf("can not read response (url = %v, method = %v, request body = %v)", request.URL, request.RequestMethod, request.Body))
	}
	res.Body.Close()

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
			if c.retryWait != 0 {
				time.Sleep(time.Duration(c.retryWait) * time.Millisecond)
			}
			continue
		} else if noRetry && err != nil {
			log.Printf("request is failure (url = %v, method = %v, reason = %v)", request.URL, request.RequestMethod, err)
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
		return nil, nil, errors.Wrapf(err, "unsupported request method (url = %v, method = %v)", request.URL, requestMethod)
	}
	res, resBody, err := c.retryRequest(request, noRetry)
	if err != nil {
		return res, resBody, errors.Wrapf(err, "request is failure (url = %v, method = %v)", request.URL, request.RequestMethod)
	}
	return res, resBody, err
}

func NewHTTPClient(retry int, retryWait int, timeout int, localAddr *net.IPAddr) *HTTPClient {
	if retry == 0 {
		retry = 200000
	}
	if timeout == 0 {
		timeout = 300
	}
	newHTTPClient := &HTTPClient{
		retry:             retry,
		retryWait:         retryWait,
		timeout:           timeout,
		resolver:          dnscache.New(time.Second * 10),
		resolverIdx:       0,
		resolverIdxMutex:  new(sync.Mutex),
		clientsCache:      make(map[string]*clientCache),
		clientsCacheMutex: new(sync.Mutex),
	}
	if localAddr == nil {
		newHTTPClient.localAddr = nil
	} else {
		newHTTPClient.localAddr = &net.TCPAddr{
			IP: localAddr.IP,
		}
	}
	return newHTTPClient
}

type WSCallback func(conn *websocket.Conn, data interface{}) error

type WSClient struct {
	readBufSize           int
	writeBufSize          int
	retry                 int
	retryWait             int
	conn                  *websocket.Conn
	connChan              chan error
	started               bool
	finished              uint32
	connectLoopFinishChan chan bool
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

func (w *WSClient) pingLoop(pingStopChan chan bool, pingStopCompleteChan chan bool) {
	for {
		select {
		case _, ok := <-pingStopChan:
			if !ok {
				close(pingStopCompleteChan)
				return
			}
		case <-time.After(5 * time.Second):
			deadline := time.Now()
			deadline.Add(30 * time.Second)
			w.conn.WriteControl(websocket.PingMessage, []byte("ping"), deadline)
		}
	}
}

func (w *WSClient) startPing() (chan bool, chan bool) {
	pingStopChan := make(chan bool)
	pingStopCompleteChan := make(chan bool)
	go w.pingLoop(pingStopChan, pingStopCompleteChan)
	return pingStopChan, pingStopCompleteChan
}

func (w *WSClient) stopPing(pingStopChan chan bool, pingStopCompleteChan chan bool) {
	close(pingStopChan)
	<-pingStopCompleteChan
	return
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
				time.Sleep(1 * time.Second)
				continue
			}
			if response.StatusCode < 200 && response.StatusCode <= 300 {
				log.Printf("can not dial, retry ... (url = %v, header = %v, reason = %v)", requestURL, requestHeaders, err)
				time.Sleep(1 * time.Second)
				continue
			}
			log.Printf("can not dial (URL = %v, header = %v)", requestURL, requestHeaders)
			continue
		}
		if !w.started {
			w.connChan <- nil
			w.started = true
		}
		w.conn = conn
		pingStopChan, pingStopCompleteChan := w.startPing()
		err = w.messageLoop(callback, callbackData)
		if err != nil {
			log.Printf("error occuered in message loop: %v", err)
		}
		w.stopPing(pingStopChan, pingStopCompleteChan)
		conn.Close()
		i = 0
		if atomic.LoadUint32(&w.finished) == 1 {
			log.Printf("connect loop finished")
			close(w.connectLoopFinishChan)
			return
		}
	}
	log.Printf("give up retry (url = %v, header = %v)", requestURL, requestHeaders)
}

func (w *WSClient) Start(callback WSCallback, callbackData interface{}, requestURL string, requestHeaders map[string]string) error {
	atomic.StoreUint32(&w.finished, 0)
	go w.connect(callback, callbackData, requestURL, requestHeaders)
	select {
	case err := <-w.connChan:
		if err != nil {
			close(w.connChan)
			return errors.Wrap(err, fmt.Sprintf("can not connect (URL = %v, header = %v)", requestURL, requestHeaders))
		}
	}
	return nil
}

func (w *WSClient) Stop() {
	atomic.StoreUint32(&w.finished, 1)
	close(w.connChan)
	select {
	case <-time.After(1 * time.Second):
	case <-w.connectLoopFinishChan:
	}
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
		readBufSize:           readBufSize,
		writeBufSize:          writeBufSize,
		retry:                 retry,
		retryWait:             retryWait,
		connChan:              make(chan error),
		connectLoopFinishChan: make(chan bool),
	}
}
