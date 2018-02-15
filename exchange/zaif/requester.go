package zaif

import (
	"github.com/pkg/errors"
	"github.com/AutomaticCoinTrader/ACT/utility"
	"net/url"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"time"
	"strconv"
	"fmt"
	"net/http"
	"encoding/json"
	"sync/atomic"
)

type Requester struct {
	httpClient   *utility.HTTPClient
	wsClients    map[string]*utility.WSClient
	readBufSize  int
	writeBufSize int
	retry        int
	retryWait    int
	key          string
	secret       string
}

type urlBuilder int

const (
	Public urlBuilder = iota
	Trade
)

var seq int64

type requestFunc func(request *utility.HTTPRequest) (interface{}, *http.Response, []byte, error)

func (b urlBuilder) buildURL(resource string) (string) {
	switch b {
	case Public:
		return "https://api.zaif.jp/api/1" + "/" + resource
	case Trade:
		return "https://api.zaif.jp/tapi" + "/" + resource
	default:
		panic("not reached")
	}
}

func (b urlBuilder) getURL() (string) {
	switch b {
	case Public:
		return "https://api.zaif.jp/api/1"
	case Trade:
		return "https://api.zaif.jp/tapi"
	default:
		panic("not reaced")
	}
}

func (r *Requester) makePublicRequest(resource string, params string) (*utility.HTTPRequest) {
	u := Public.buildURL(resource)
	if params != "" {
		u += "?" + params
	}
	return &utility.HTTPRequest{
		URL: u,
	}
}

func (r *Requester) makeTradeRequest(method string, params string) (*utility.HTTPRequest) {
	u := Trade.getURL()
	values := url.Values{}
	values.Set("nonce", r.getNonce())
	values.Set("method", method)
	body := values.Encode()
	if params != "" {
		body += "&" + params
	}
	mac := hmac.New(sha512.New, []byte(r.secret))
	mac.Write([]byte(body))
	sign := hex.EncodeToString(mac.Sum(nil))
	headers := make(map[string]string)
	headers["Conent-Type"] = "application/x-www-form-urlencoded"
	headers["Key"] = r.key
	headers["Sign"] = sign
	//log.Printf("key = %v, sign = %v", r.key, sign)
	return &utility.HTTPRequest{
		URL:     u,
		Headers: headers,
		Body:    body,
	}
}

func (r *Requester) getNonce() (string) {
	now := time.Now()
	s:= atomic.AddInt64(&seq, 1)
	return strconv.FormatInt(now.Unix(), 10) + "." + fmt.Sprintf("%06d", now.Nanosecond() / 1000) + fmt.Sprintf("%03d", s % 1000)
}

func (r *Requester) unmarshal(requestFunc requestFunc, request *utility.HTTPRequest) (interface{}, *http.Response, error) {
	newRes, res, resBody, err := requestFunc(request)
	err = json.Unmarshal(resBody, newRes)
	if err != nil {
		return newRes, res, errors.Wrap(err, fmt.Sprintf("can not unmarshal response (url = %v, method = %v)", request.URL, request.RequestMethod))
	}
	return newRes, res, err
}

// NewRequester is create requester
func NewRequester(key string, secret string, retry int, retryWait, timeout int, readBufSize int, writeBufSize int) (*Requester) {
	return &Requester{
		httpClient:   utility.NewHTTPClient(retry, retryWait, timeout),
		wsClients:    make(map[string]*utility.WSClient),
		readBufSize:  readBufSize,
		writeBufSize: writeBufSize,
		retry:        retry,
		retryWait:    retryWait,
		key:          key,
		secret:       secret,
	}
}
