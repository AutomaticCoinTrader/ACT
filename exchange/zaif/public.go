package zaif

import (
	"github.com/pkg/errors"
	"github.com/gorilla/websocket"
	"github.com/AutomaticCoinTrader/ACT/utility"
	"path"
	"fmt"
	"net/http"
	"encoding/json"
	"log"
	"time"
)

// PublicCurrenciesResponse is response of currencies
type PublicCurrenciesResponse []PublicCurrencyResponse

// PublicCurrency is response of currency
type PublicCurrencyResponse struct {
	Name    string `json:"name"`
	IsToken bool   `json:"is_token"`
}

// GetCurrencies is get currencies
func (r *Requester) Currencies(currency string) (*PublicCurrenciesResponse, *utility.HTTPRequest, *http.Response, error) {
	for {
		request := r.MakePublicRequest(path.Join("currencies", currency), "")
		newRes, response, err := r.unmarshal(func(request *utility.HTTPRequest) (interface{}, *http.Response, []byte, error) {
			httpClient := r.getHttpClient()
			res, resBody, err := httpClient.DoRequest(utility.HTTPMethodGET, request, true)
			if err != nil {
				return nil, res, resBody, errors.Wrap(err, fmt.Sprintf("can not get currencies (url = %v)", request.URL))
			}
			newRes := new(PublicCurrenciesResponse)
			return newRes, res, resBody, err
		}, request)
		if err != nil {
			time.Sleep(time.Duration(r.retryWait) * time.Millisecond)
			log.Printf("retry currencies (currency = %v, err: %v)", currency, err)
			continue
		}
		return newRes.(*PublicCurrenciesResponse), request, response, err
	}
}

// PublicCurrencyPairsResponse is response of currency pairs
type PublicCurrencyPairsResponse []PublicCurrencyPairResponse

// Currency is response of CurrencyPair
type PublicCurrencyPairResponse struct {
	AuxUnitMin   float64 `json:"aux_unit_min"`
	AuxUnitStep  float64 `json:"aux_unit_step"`
	CurrencyPair string  `json:"currency_pair"`
	Description  string  `json:"description"`
	EventNumber  int64   `json:"event_number"`
	IsToken      bool    `json:"is_token"`
	ItemUnitMin  float64 `json:"item_unit_min"`
	ItemUnitStep float64 `json:"item_unit_step"`
	Name         string  `json:"name"`
	Title        string  `json:"title"`
}

// CurrencyPairs is get currency pairs
func (r *Requester) CurrencyPairs(currencyPair string) (*PublicCurrencyPairsResponse, *utility.HTTPRequest, *http.Response, error) {
	for {
		request := r.MakePublicRequest(path.Join("currency_pairs", currencyPair), "")
		newRes, response, err := r.unmarshal(func(request *utility.HTTPRequest) (interface{}, *http.Response, []byte, error) {
			httpClient := r.getHttpClient()
			res, resBody, err := httpClient.DoRequest(utility.HTTPMethodGET, request, true)
			if err != nil {
				return nil, res, resBody, errors.Wrap(err, fmt.Sprintf("can not get currency pairs (url = %v)", request.URL))
			}
			newRes := new(PublicCurrencyPairsResponse)
			return newRes, res, resBody, err
		}, request)
		if err != nil {
			time.Sleep(time.Duration(r.retryWait) * time.Millisecond)
			log.Printf("retry currency pairs (currency pair = %v, err: %v)", currencyPair, err)
			continue
		}
		return newRes.(*PublicCurrencyPairsResponse), request, response, err
	}
}

// PublicLastPriceResponse is response of last price
type PublicLastPriceResponse struct {
	LastPrice float64 `json:"last_price"`
}

// LastPricee is get last place
func (r *Requester) LastPrice(currencyPair string) (*PublicLastPriceResponse, *utility.HTTPRequest, *http.Response, error) {
	for {
		request := r.MakePublicRequest(path.Join("last_price", currencyPair), "")
		newRes, response, err := r.unmarshal(func(request *utility.HTTPRequest) (interface{}, *http.Response, []byte, error) {
			httpClient := r.getHttpClient()
			res, resBody, err := httpClient.DoRequest(utility.HTTPMethodGET, request, true)
			if err != nil {
				return nil, res, resBody, errors.Wrap(err, fmt.Sprintf("can not get last price (url = %v)", request.URL))
			}
			newRes := new(PublicLastPriceResponse)
			return newRes, res, resBody, err
		}, request)
		if err != nil {
			time.Sleep(time.Duration(r.retryWait) * time.Millisecond)
			log.Printf("retry last price (currency pair = %v, err: %v)", currencyPair, err)
			continue
		}
		return newRes.(*PublicLastPriceResponse), request, response, err
	}
}

// PublicTickerResponse is response of ticker
type PublicTickerResponse struct {
	Ask    float64 `json:"ask"`
	Bid    float64 `json:"bid"`
	High   float64 `json:"high"`
	Last   float64 `json:"last"`
	Low    float64 `json:"low"`
	Volume float64 `json:"volume"`
	Vwap   float64 `json:"vwap"`
}

// Ticker is get ticker
func (r *Requester) Ticker(currencyPair string) (*PublicTickerResponse, *utility.HTTPRequest, *http.Response, error) {
	for {
		request := r.MakePublicRequest(path.Join("ticker", currencyPair), "")
		newRes, response, err := r.unmarshal(func(request *utility.HTTPRequest) (interface{}, *http.Response, []byte, error) {
			httpClient := r.getHttpClient()
			res, resBody, err := httpClient.DoRequest(utility.HTTPMethodGET, request, true)
			if err != nil {
				return nil, res, resBody, errors.Wrap(err, fmt.Sprintf("can not get ticker (url = %v)", request.URL))
			}
			newRes := new(PublicTickerResponse)
			return newRes, res, resBody, err
		}, request)
		if err != nil {
			time.Sleep(time.Duration(r.retryWait) * time.Millisecond)
			log.Printf("retry ticker (currency pair = %v, err: %v)", currencyPair, err)
			continue
		}
		return newRes.(*PublicTickerResponse), request, response, err
	}
}

// PublicTradesResponse is response of trades
type PublicTradesResponse []PublicTradeResponse

// PublicTradeResponse is response of trade
type PublicTradeResponse struct {
	Amount       float64 `json:"amount"`
	CurrencyPair string  `json:"currency_pair"`
	Date         int64   `json:"date"`
	Price        float64 `json:"price"`
	Tid          int64   `json:"tid"`
	TradeType    string  `json:"trade_type"`
}

// Trades is get trades
func (r *Requester) Trades(currencyPair string) (*PublicTradesResponse, *utility.HTTPRequest, *http.Response, error) {
	for {
		request := r.MakePublicRequest(path.Join("trades", currencyPair), "")
		newRes, response, err := r.unmarshal(func(request *utility.HTTPRequest) (interface{}, *http.Response, []byte, error) {
			httpClient := r.getHttpClient()
			res, resBody, err := httpClient.DoRequest(utility.HTTPMethodGET, request, true)
			if err != nil {
				return nil, res, resBody, errors.Wrap(err, fmt.Sprintf("can not get trades (url = %v)", request.URL))
			}
			newRes := new(PublicTradesResponse)
			return newRes, res, resBody, err
		}, request)
		if err != nil {
			time.Sleep(time.Duration(r.retryWait) * time.Millisecond)
			log.Printf("retry terades (currency pair = %v, err: %v)", currencyPair, err)
			continue
		}
		return newRes.(*PublicTradesResponse), request, response, err
	}
}

// PublicDepthReaponse is response of depth
type PublicDepthReaponse struct {
	Asks [][]float64 `json:"asks"`
	Bids [][]float64 `json:"bids"`
}

// DepthNoRetry is get depth with no retry
func (r *Requester) DepthNoRetry(currencyPair string) (*PublicDepthReaponse, *utility.HTTPRequest, *http.Response, error) {
	request := r.MakePublicRequest(path.Join("depth", currencyPair), "")
	newRes, response, err := r.unmarshal(func(request *utility.HTTPRequest) (interface{}, *http.Response, []byte, error) {
		httpClient := r.getHttpClient()
		res, resBody, err := httpClient.DoRequest(utility.HTTPMethodGET, request, true)
		if err != nil {
			return nil, res, resBody, errors.Wrap(err, fmt.Sprintf("can not get depth (url = %v)", request.URL))
		}
		newRes := new(PublicDepthReaponse)
		return newRes, res, resBody, err
	}, request)
	if err != nil {
		return nil, request, response, err
	} else {
		return newRes.(*PublicDepthReaponse), request, response, nil
	}
}

// Depth is get depth
func (r *Requester) Depth(currencyPair string) (*PublicDepthReaponse, *utility.HTTPRequest, *http.Response, error) {
	for {
		newRes, request, response, err := r.DepthNoRetry(currencyPair)
		if err != nil {
			time.Sleep(time.Duration(r.retryWait) * time.Millisecond)
			log.Printf("retry depth (currency pair = %v, err: %v)", currencyPair, err)
			continue
		}
		return newRes, request, response, err
	}
}

type StreamingCallback func(currencyPair string, streamingResponse *StreamingResponse, streamingCallbackData interface{}) (error)

type StreamingResponse struct {
	Asks         [][]float64 `json:"asks"`
	Bids         [][]float64 `json:"bids"`
	CurrencyPair string      `json:"currency_pair"`
	LastPrice struct {
		Action string  `json:"action"`
		Price  float64 `json:"price"`
	} `json:"last_price"`
	Timestamp string                     `json:"timestamp"`
	Trades    []*StreamingTradesResponse `json:"trades"`
}

type StreamingTradesResponse struct {
	Amount       float64 `json:"amount"`
	CurrentyPair string  `json:"currenty_pair"`
	Date         int64   `json:"date"`
	Price        float64 `json:"price"`
	Tid          int64   `json:"tid"`
	TradeType    string  `json:"trade_type"`
}

type streaminCallbackData struct {
	currencyPair string
	callback     StreamingCallback
	callbackData interface{}
}

func (r Requester) streamingCallback(conn *websocket.Conn, userCallbackData interface{}) (error) {
	streaminCallbackData := userCallbackData.(*streaminCallbackData)
	messageType, message, err := conn.ReadMessage()
	if err != nil {
		return errors.Wrap(err, "can not read message of streaming")
	}
	if messageType != websocket.TextMessage {
		log.Printf("unsupported message type (message type = %v, message = %v)", messageType, message)
		return nil
	}
	newRes := new(StreamingResponse)
	err = json.Unmarshal(message, newRes)
	if err != nil {
		log.Printf("can not unmarshal message of streaming (%v)", message)
		return nil
	}
	err = streaminCallbackData.callback(streaminCallbackData.currencyPair, newRes, streaminCallbackData.callbackData)
	if err != nil {
		log.Printf("call back error of streaming (%v)", err)
		return nil
	}
	return nil
}


func (r Requester) StreamingStart(currencyPair string, callback StreamingCallback, callbackData interface{}) (error) {
	_, ok := r.wsClients[currencyPair]
	if ok {
		return errors.Errorf("already exists streaming (currency pair = %v)", currencyPair)
	}
	log.Printf("start streaming (currency pair = %v)", currencyPair)
	requestURL := "wss://ws.zaif.jp/stream?currency_pair=" + currencyPair
	streaminCallbackData := &streaminCallbackData{
		currencyPair: currencyPair,
		callback:     callback,
		callbackData: callbackData,
	}
	newClient := utility.NewWSClient(r.readBufSize, r.writeBufSize, r.retry, r.retryWait)
	err := newClient.Start(r.streamingCallback, streaminCallbackData, requestURL, nil)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("can not start streaming (url = %v)", requestURL))
	}
	r.wsClients[currencyPair] = newClient

	return nil
}

func (r Requester) StreamingStop(currencyPair string) {
	client, ok := r.wsClients[currencyPair]
	if !ok {
		log.Printf("not found streaming (currency pair = %v)", currencyPair)
		return
	}
	client.Stop()
}

type ProxyStreamingCallback func(currencyPair string, proxyStreamingResponse *PublicDepthReaponse, streamingCallbackData interface{}) (error)

type proxyStreaminCallbackData struct {
	currencyPair string
	callback     ProxyStreamingCallback
	callbackData interface{}
}

func (r Requester) proxyStreamingCallback(conn *websocket.Conn, userCallbackData interface{}) (error) {
	proxyStreaminCallbackData := userCallbackData.(*proxyStreaminCallbackData)
	messageType, message, err := conn.ReadMessage()
	if err != nil {
		return errors.Wrap(err, "can not read message of proxy streaming")
	}
	if messageType != websocket.TextMessage {
		log.Printf("unsupported message type (message type = %v, message = %v)", messageType, message)
		return nil
	}
	newRes := new(PublicDepthReaponse)
	err = json.Unmarshal(message, newRes)
	if err != nil {
		log.Printf("can not unmarshal message of proxy streaming (%v)", message)
		return nil
	}
	err = proxyStreaminCallbackData.callback(proxyStreaminCallbackData.currencyPair, newRes, proxyStreaminCallbackData.callbackData)
	if err != nil {
		log.Printf("call back error of proxy streaming (%v)", err)
		return nil
	}
	return nil
}

func (r Requester) ProxyStreamingStart(addrPort string, currencyPair string, callback ProxyStreamingCallback, callbackData interface{}) (error) {
	clientId := addrPort + "@" + currencyPair
	_, ok := r.proxyWsClients[clientId]
	if ok {
		return errors.Errorf("already exists proxy streaming (currency pair = %v)", currencyPair)
	}
	log.Printf("start proxy streaming (currency pair = %v)", currencyPair)
	requestURL := fmt.Sprintf("ws://%v/%v", addrPort, currencyPair)
	proxyStreaminCallbackData := &proxyStreaminCallbackData{
		currencyPair: currencyPair,
		callback:     callback,
		callbackData: callbackData,
	}
	newClient := utility.NewWSClient(r.readBufSize, r.writeBufSize, r.retry, r.retryWait)
	err := newClient.Start(r.proxyStreamingCallback, proxyStreaminCallbackData, requestURL, nil)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("can not start proxy streaming (url = %v)", requestURL))
	}
	r.proxyWsClients[clientId] = newClient

	return nil
}

func (r Requester) ProxyStreamingStop(addrPort string, currencyPair string) {
	clientId := addrPort + "@" + currencyPair
	client, ok := r.proxyWsClients[clientId]
	if !ok {
		log.Printf("not found proxy streaming (currency pair = %v)", currencyPair)
		return
	}
	client.Stop()
}