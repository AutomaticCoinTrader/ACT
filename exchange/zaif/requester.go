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
	"sync"
	"log"
	"net"
)

type RequesterKey struct {
	Key    string
	Secret string
}

type Requester struct {
	httpClients           []*utility.HTTPClient
	httpClientsIndex      int
	httpClientsMutex      *sync.Mutex
	wsClients             map[string]*utility.WSClient
	proxyWsClients        map[string]*utility.WSClient
	readBufSize           int
	writeBufSize          int
	retry                 int
	retryWait             int
	keys                  []*RequesterKey
	keyIndex              int
	keysMutex             *sync.Mutex
	publicApiHistory      []int64
	publicApiHistoryMutex *sync.Mutex
	tradeApiHistory       []int64
	tradeApiHistoryMutex  *sync.Mutex
}

type urlBuilder int

const (
	Public urlBuilder = iota
	Trade
)

const (
	restrictionWait    = 1000
	publicApiGurdCount = 100
	tradeApiGurdCount  = 50
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

func (r *Requester) MakePublicRequest(resource string, params string) (*utility.HTTPRequest) {
	r.publicApiHistoryMutex.Lock()
	for {
		// waitを入れる処理
		// ４０３ Forbiddenを緩和する
		lastIdx := -1
		var lastTs int64 = 0
		now := time.Now()
		for idx, ts := range r.publicApiHistory {
			if ts > now.UnixNano()-time.Second.Nanoseconds() {
				// １秒以内のものになったらbreak
				lastIdx = idx
				lastTs = ts
				break
			}
		}
		if lastIdx == -1 && len(r.publicApiHistory) > 0 {
			// 全部ふるいから消す
			r.publicApiHistory = make([]int64, 0, publicApiGurdCount)
		} else if lastIdx > 0 {
			// 古いやつだけ消す
			r.publicApiHistory = r.publicApiHistory[lastIdx:]
		}
		if len(r.publicApiHistory) < publicApiGurdCount {
			r.publicApiHistory = append(r.publicApiHistory, now.UnixNano())
			break
		}
		if lastTs > 0 {
			time.Sleep(time.Duration(lastTs-(now.UnixNano()-time.Second.Nanoseconds())) * time.Nanosecond)
		} else {
			log.Print("called trade api with no wait")
		}
	}
	r.publicApiHistoryMutex.Unlock()
	u := Public.buildURL(resource)
	if params != "" {
		u += "?" + params
	}
	headers := make(map[string]string)
	headers["Connection"] = "close"
	return &utility.HTTPRequest{
		Headers: headers,
		URL:     u,
	}
}

func (r *Requester) makeTradeRequest(method string, params string) (*utility.HTTPRequest) {
	r.tradeApiHistoryMutex.Lock()
	for {
		// waitを入れる処理
		// ４０３ Forbiddenを緩和する
		lastIdx := -1
		var lastTs int64 = 0
		now := time.Now()
		for idx, ts := range r.tradeApiHistory {
			if ts > now.UnixNano()-time.Second.Nanoseconds() {
				// １秒以内のものになったらbreak
				lastIdx = idx
				lastTs = ts
				break
			}
		}
		if lastIdx == -1 && len(r.tradeApiHistory) > 0 {
			// 全部ふるいから消す
			r.tradeApiHistory = make([]int64, 0, tradeApiGurdCount)
		} else if lastIdx > 0 {
			// 古いやつだけ消す
			r.tradeApiHistory = r.tradeApiHistory[lastIdx:]
		}
		if len(r.tradeApiHistory) < tradeApiGurdCount {
			r.tradeApiHistory = append(r.tradeApiHistory, now.UnixNano())
			break
		}
		if lastTs > 0 {
			time.Sleep(time.Duration(lastTs-(now.UnixNano()-time.Second.Nanoseconds())) * time.Nanosecond)
		} else {
			log.Print("called trade api with no wait")
		}
	}
	r.tradeApiHistoryMutex.Unlock()
	u := Trade.getURL()
	values := url.Values{}
	values.Set("nonce", r.getNonce())
	values.Set("method", method)
	body := values.Encode()
	if params != "" {
		body += "&" + params
	}
	r.keysMutex.Lock()
	key := r.keys[r.keyIndex].Key
	secret := r.keys[r.keyIndex].Secret
	r.keyIndex = (r.keyIndex + 1) % len(r.keys)
	r.keysMutex.Unlock()
	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write([]byte(body))
	sign := hex.EncodeToString(mac.Sum(nil))
	headers := make(map[string]string)
	headers["Connection"] = "keep-alive"
	headers["Conent-Type"] = "application/x-www-form-urlencoded"
	headers["Key"] = key
	headers["Sign"] = sign
	//log.Printf("key = %v, sign = %v", key, sign)
	return &utility.HTTPRequest{
		URL:     u,
		Headers: headers,
		Body:    body,
	}
}

func (r *Requester) getNonce() (string) {
	now := time.Now()
	s := atomic.AddInt64(&seq, 1)
	return strconv.FormatInt(now.Unix(), 10) + "." + fmt.Sprintf("%06d", now.Nanosecond()/1000) + fmt.Sprintf("%03d", s%1000)
}

func (r *Requester) unmarshal(requestFunc requestFunc, request *utility.HTTPRequest) (interface{}, *http.Response, error) {
	newRes, res, resBody, err := requestFunc(request)
	err = json.Unmarshal(resBody, newRes)
	if err != nil {
		return newRes, res, errors.Wrap(err, fmt.Sprintf("can not unmarshal response (url = %v, method = %v)", request.URL, request.RequestMethod))
	}
	return newRes, res, err
}

func (r *Requester) getHttpClient() (*utility.HTTPClient) {
	r.httpClientsMutex.Lock()
	defer r.httpClientsMutex.Unlock()
	httpClient := r.httpClients[r.httpClientsIndex]
	r.httpClientsIndex += 1
	if r.httpClientsIndex >= len(r.httpClients) {
		r.httpClientsIndex = 0
	}
	return httpClient
}

// NewRequester is create requester
func NewRequester(keys []*RequesterKey, bindAddresses []string, retry int, retryWait, timeout int, readBufSize int, writeBufSize int) (*Requester, error) {
	requester := &Requester{
		httpClients:           make([]*utility.HTTPClient, 0),
		httpClientsIndex:      0,
		httpClientsMutex:      new(sync.Mutex),
		wsClients:             make(map[string]*utility.WSClient),
		proxyWsClients:        make(map[string]*utility.WSClient),
		readBufSize:           readBufSize,
		writeBufSize:          writeBufSize,
		retry:                 retry,
		retryWait:             retryWait,
		keys:                  keys,
		keyIndex:              0,
		keysMutex:             new(sync.Mutex),
		publicApiHistory:      make([]int64, 0, publicApiGurdCount),
		publicApiHistoryMutex: new(sync.Mutex),
		tradeApiHistory:       make([]int64, 0, tradeApiGurdCount),
		tradeApiHistoryMutex:  new(sync.Mutex),
	}
	if bindAddresses == nil || len(bindAddresses) == 0 {
		requester.httpClients = append(requester.httpClients, utility.NewHTTPClient(retry, retryWait, timeout, nil), )
	} else {
		for _, addr := range bindAddresses {
			localAddr, err := net.ResolveIPAddr("ip", addr)
			if err != nil {
				return nil, errors.Wrap(err, "can not resolve ip address")
			}
			requester.httpClients = append(requester.httpClients, utility.NewHTTPClient(retry, retryWait, timeout, localAddr))
		}
	}
	return requester, nil
}
