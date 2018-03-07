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
)

type RequesterKey struct {
	Key    string
	Secret string
}

type Requester struct {
	httpClient   *utility.HTTPClient
	wsClients    map[string]*utility.WSClient
	readBufSize  int
	writeBufSize int
	retry        int
	retryWait    int
	keys         []*RequesterKey
	keyIndex     int
	keysMutex    *sync.Mutex
	lastPulblicApiNanoTs int64
	lastPulblicApiNanoTsMutex *sync.Mutex
	lastTradeApiHistory  []int64
	lastTradeApiHistoryMutex *sync.Mutex
}

type urlBuilder int

const (
	Public urlBuilder = iota
	Trade
)

const (
	restrictionWait = 1000
	insufficientWait = 1000
	publicApiGurdTime = 12 // 本来は10
	tradeApiGurdCount = 40 // 本来は50
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
	// waitを入れる処理
	r.lastPulblicApiNanoTsMutex.Lock()
	for {
		// waitを入れる処理
		// ４０３ Forbiddenを緩和する
		now := time.Now()
		if now.UnixNano() >= r.lastPulblicApiNanoTs+(publicApiGurdTime*time.Millisecond.Nanoseconds()) {
			r.lastPulblicApiNanoTs = now.UnixNano()
			break
		}
		time.Sleep(time.Duration(r.lastPulblicApiNanoTs + (publicApiGurdTime*time.Millisecond.Nanoseconds()) - now.UnixNano()) * time.Nanosecond)
	}
	r.lastPulblicApiNanoTsMutex.Unlock()
	u := Public.buildURL(resource)
	if params != "" {
		u += "?" + params
	}
	return &utility.HTTPRequest{
		URL: u,
	}
}

func (r *Requester) makeTradeRequest(method string, params string) (*utility.HTTPRequest) {
	r.lastTradeApiHistoryMutex.Lock()
	for {
		// waitを入れる処理
		// ４０３ Forbiddenを緩和する
		lastIdx := -1
		var lastTs int64 = 0
		now := time.Now()
		for idx, ts := range r.lastTradeApiHistory {
			if ts > now.UnixNano() -  time.Second.Nanoseconds() {
				// １秒以内のものになったらbreak
				lastIdx = idx
				lastTs = ts
				break;
			}
		}
		if lastIdx == -1 && len(r.lastTradeApiHistory) > 0 {
			// 全部ふるいから消す
			r.lastTradeApiHistory = make([]int64, 0, tradeApiGurdCount)
		} else if lastIdx > 0 {
			// 古いやつだけ消す
			r.lastTradeApiHistory = r.lastTradeApiHistory[lastIdx:]
		}
		if len(r.lastTradeApiHistory) < tradeApiGurdCount {
			r.lastTradeApiHistory = append(r.lastTradeApiHistory, now.UnixNano())
			break
		}
		if lastTs > 0 {
			time.Sleep(time.Duration(lastTs-(now.UnixNano()-time.Second.Nanoseconds())) * time.Nanosecond)
		} else {
			log.Print("called trade api with no wait")
		}
	}
	r.lastTradeApiHistoryMutex.Unlock()
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
func NewRequester(keys []*RequesterKey, retry int, retryWait, timeout int, readBufSize int, writeBufSize int) (*Requester) {
	return &Requester{
		httpClient:   utility.NewHTTPClient(retry, retryWait, timeout),
		wsClients:    make(map[string]*utility.WSClient),
		readBufSize:  readBufSize,
		writeBufSize: writeBufSize,
		retry:        retry,
		retryWait:    retryWait,
		keys:         keys,
		keyIndex:     0,
		keysMutex:    new(sync.Mutex),
		lastPulblicApiNanoTs: 0,
		lastPulblicApiNanoTsMutex: new(sync.Mutex),
		lastTradeApiHistory: make([]int64, 0, tradeApiGurdCount),
		lastTradeApiHistoryMutex: new(sync.Mutex),
	}
}
