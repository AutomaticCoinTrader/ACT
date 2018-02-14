package zaif

import (
	"github.com/pkg/errors"
	"github.com/google/go-querystring/query"
	"github.com/AutomaticCoinTrader/ACT/utility"
	"fmt"
	"math"
	"time"
	"strconv"
	"net/http"
	"github.com/AutomaticCoinTrader/ACT/exchange"
	"log"
	"net/url"
)

func (r *Requester) GetMinPriceUnit(currencyPair string) (float64) {
	switch currencyPair {
	case "btc_jpy":
		return 5
	case "xem_jpy":
		return 0.0001
	case "mona_jpy":
		return 0.1
	case "bch_jpy":
		return 5
	case "eth_jpy":
		return 5
	case "zaif_jpy":
		return 0.0001
	case "pepecash_jpy":
		return 0.0001
	case "xem_btc":
		return 0.00000001
	case "mona_btc":
		return 0.00000001
	case "bch_btc":
		return 0.0001
	case "eth_btc":
		return 0.0001
	case "zaif_btc":
		return 0.00000001
	case "pepecash_btc":
		return 0.00000001
	default:
		return -1
	}
}

func (r *Requester) getPricePrec(currencyPair string) (int) {
	switch currencyPair {
	case "btc_jpy":
		return 0
	case "xem_jpy":
		return 4
	case "mona_jpy":
		return 1
	case "bch_jpy":
		return 0
	case "eth_jpy":
		return 0
	case "zaif_jpy":
		return 4
	case "pepecash_jpy":
		return 4
	case "xem_btc":
		return 8
	case "mona_btc":
		return 8
	case "bch_btc":
		return 4
	case "eth_btc":
		return 4
	case "zaif_btc":
		return 8
	case "pepecash_btc":
		return 8
	default:
		return -1
	}
}

func (r *Requester) GetMinAmountUnit(currencyPair string) (float64) {
	switch currencyPair {
	case "btc_jpy":
		return 0.0001
	case "xem_jpy":
		return 0.1
	case "mona_jpy":
		return 1
	case "bch_jpy":
		return 0.0001
	case "eth_jpy":
		return 0.0001
	case "zaif_jpy":
		return 0.1
	case "pepecash_jpy":
		return 0.0001
	case "xem_btc":
		return 1
	case "mona_btc":
		return 1
	case "bch_btc":
		return 0.0001
	case "eth_btc":
		return 0.0001
	case "zaif_btc":
		return 1
	case "pepecash_btc":
		return 1
	default:
		return -1
	}
}

func (r *Requester) getAmountPrec(currencyPair string) (int) {
	switch currencyPair {
	case "btc_jpy":
		return 4
	case "xem_jpy":
		return 1
	case "mona_jpy":
		return 0
	case "bch_jpy":
		return 4
	case "eth_jpy":
		return 4
	case "zaif_jpy":
		return 1
	case "pepecash_jpy":
		return 4
	case "xem_btc":
		return 0
	case "mona_btc":
		return 0
	case "bch_btc":
		return 4
	case "eth_btc":
		return 4
	case "zaif_btc":
		return 0
	case "pepecash_btc":
		return 0
	default:
		return -1
	}
}


func (r *Requester) GetMinFee(currency string) (float64) {
	switch currency {
	case "btc":
		return 0.00001
	case "xem":
		return 2
	case "mona":
		return 0
	case "bch":
		return 0
	case "eth":
		return 0
	case "zaif":
		return 0
	case "pepecash":
		return 0
	default:
		return -1
	}
}

type TradeCommonResponse struct {
	Success int    `json:"success"`
	Error   string `json:"error"`
}

func (t TradeCommonResponse) needRetry() (bool) {
	if t.Success == 0 {
		log.Printf(" error message (%v)", t.Error)
		return true
	}
	return false
}

// TradeGetInfoResponse is response of get information
type TradeGetInfoResponse struct {
	Return struct {
		Deposit struct {
			Btc      float64 `json:"btc"`
			Jpy      float64 `json:"jpy"`
			Bch      float64 `json:"bch"`
			Eth      float64 `json:"eth"`
			Mona     float64 `json:"mona"`
			Xem      float64 `json:"xem"`
			Zaif     float64 `json:"zaif"`
			Pepecash float64 `json:"pepecash"`
		} `json:"deposit"`
		Funds struct {
			Btc      float64 `json:"btc"`
			Jpy      float64 `json:"jpy"`
			Bch      float64 `json:"bch"`
			Eth      float64 `json:"eth"`
			Mona     float64 `json:"mona"`
			Xem      float64 `json:"xem"`
			Zaif     float64 `json:"zaif"`
			Pepecash float64 `json:"pepecash"`
		} `json:"funds"`
		OpenOrders int `json:"open_orders"`
		Rights struct {
			IDInfo       int64 `json:"id_info"`
			Info         int64 `json:"info"`
			PersonalInfo int64 `json:"personal_info"`
			Trade        int64 `json:"trade"`
			Withdraw     int64 `json:"withdraw"`
		} `json:"rights"`
		ServerTime int64 `json:"server_time"`
		TradeCount int64 `json:"trade_count"`
	} `json:"return"`
	TradeCommonResponse
}

// GetInfo is get informarion
func (r *Requester) GetInfo() (*TradeGetInfoResponse, *utility.HTTPRequest, *http.Response, error) {
	request := r.makeTradeRequest("get_info", "")
	newRes, response, err := r.unmarshal(func(request *utility.HTTPRequest) (interface{}, *http.Response, []byte, error) {
		res, resBody, err := r.httpClient.DoRequest(utility.HTTPMethdoPOST, request, false)
		if err != nil {
			return nil, res, resBody, errors.Wrap(err, fmt.Sprintf("can not get info (url = %v)", request.URL))
		}
		newRes := new(TradeGetInfoResponse)
		return newRes, res, resBody, err
	}, request)
	return newRes.(*TradeGetInfoResponse), request, response, err
}

// TradeGetInfo2 is response of informarion2
type TradeGetInfo2Response struct {
	Return struct {
		Deposit struct {
			Btc      float64 `json:"btc"`
			Bch      float64 `json:"bch"`
			Eth      float64 `json:"eth"`
			Jpy      float64 `json:"jpy"`
			Mona     float64 `json:"mona"`
			Xem      float64 `json:"xem"`
			Zaif     float64 `json:"zaif"`
			Pepecash float64 `json:"pepecash"`
		} `json:"deposit"`
		Funds struct {
			Btc      float64 `json:"btc"`
			Bch      float64 `json:"bch"`
			Eth      float64 `json:"eth"`
			Jpy      float64 `json:"jpy"`
			Mona     float64 `json:"mona"`
			Xem      float64 `json:"xem"`
			Zaif     float64 `json:"zaif"`
			Pepecash float64 `json:"pepecash"`
		} `json:"funds"`
		OpenOrders int `json:"open_orders"`
		Rights struct {
			Info         int64 `json:"info"`
			PersonalInfo int64 `json:"personal_info"`
			Trade        int64 `json:"trade"`
			Withdraw     int64 `json:"withdraw"`
		} `json:"rights"`
		ServerTime int64 `json:"server_time"`
	} `json:"return"`
	TradeCommonResponse
}

// GetInfo is get informarion2
func (r *Requester) GetInfo2() (*TradeGetInfo2Response, *utility.HTTPRequest, *http.Response, error) {
	request := r.makeTradeRequest("get_info2", "")
	newRes, response, err := r.unmarshal(func(request *utility.HTTPRequest) (interface{}, *http.Response, []byte, error) {
		res, resBody, err := r.httpClient.DoRequest(utility.HTTPMethdoPOST, request, false)
		if err != nil {
			return nil, res, resBody, errors.Wrap(err, fmt.Sprintf("can not get info2 (url = %v)", request.URL))
		}
		newRes := new(TradeGetInfo2Response)
		return newRes, res, resBody, err
	}, request)
	return newRes.(*TradeGetInfo2Response), request, response, err
}

// TradeGetPersonalInfo is response of get id information
type TradeGetPersonalInfoResponse struct {
	Return struct {
		IconPath        string `json:"icon_path"`
		RankingNickname string `json:"ranking_nickname"`
	} `json:"return"`
	TradeCommonResponse
}

// GetPersonalInfo is get personal information
func (r *Requester) GetPersonalInfo() (*TradeGetPersonalInfoResponse, *utility.HTTPRequest, *http.Response, error) {
	for {
		request := r.makeTradeRequest("get_personal_info", "")
		newRes, response, err := r.unmarshal(func(request *utility.HTTPRequest) (interface{}, *http.Response, []byte, error) {
			res, resBody, err := r.httpClient.DoRequest(utility.HTTPMethdoPOST, request, false)
			if err != nil {
				return nil, res, resBody, errors.Wrap(err, fmt.Sprintf("can not get personal info (url = %v)", request.URL))
			}
			newRes := new(TradeGetPersonalInfoResponse)
			return newRes, res, resBody, err
		}, request)
		if newRes.(*TradeGetPersonalInfoResponse).needRetry() {
			continue
		}
		return newRes.(*TradeGetPersonalInfoResponse), request, response, err
	}
}

// TradeGetIDInfo is response of get id information
type TradeGetIDInfoResponse struct {
	Return struct {
		User struct {
			ID        int64  `json:"id"`
			Email     string `json:"email"`
			Name      string `json:"name"`
			Kana      string `json:"kana"`
			Certified bool   `json:"certified"`
		} `json:"user"`
	} `json:"return"`
	TradeCommonResponse
}

// GetPersonalInfo is get id information
func (r *Requester) GetIDInfo() (*TradeGetIDInfoResponse, *utility.HTTPRequest, *http.Response, error) {
	for {
		request := r.makeTradeRequest("get_id_info", "")
		newRes, response, err := r.unmarshal(func(request *utility.HTTPRequest) (interface{}, *http.Response, []byte, error) {
			res, resBody, err := r.httpClient.DoRequest(utility.HTTPMethdoPOST, request, false)
			if err != nil {
				return nil, res, resBody, errors.Wrap(err, fmt.Sprintf("can not get id info (url = %v)", request.URL))
			}
			newRes := new(TradeGetIDInfoResponse)
			return newRes, res, resBody, err
		}, request)
		if newRes.(*TradeGetIDInfoResponse).needRetry() {
			continue
		}
		return newRes.(*TradeGetIDInfoResponse), request, response, err
	}
}

// TradeHistoryParams is parameter of trade history
type TradeHistoryParams struct {
	From         int64  `url:"from,omitempty"`
	Count        int64  `url:"count,omitempty"`
	FromID       int64  `url:"from_id,omitempty"`
	EndID        int64  `url:"end_id"`
	Order        string `url:"order,omitempty"`
	Since        int64  `url:"since,omitempty"`
	End          int64  `url:"end"`
	CurrencyPair string `url:"currency_pair,omitempty"`
	IsToken      bool   `url:"is_token,omitempty"`
}

// TradeHistoryParams is create TradeHistoryParams
func (r *Requester) NewTradeHistoryParams() (*TradeHistoryParams) {
	return &TradeHistoryParams{}
}

// TradeHistoryResponse is response of trade history
type TradeHistoryResponse struct {
	Return map[string]TradeHistoryRecordResponse `json:"return"`
	TradeCommonResponse
}

// TradeHistoryRecordResponse is response trade history record
type TradeHistoryRecordResponse struct {
	Action       string  `json:"action"`
	Amount       float64 `json:"amount"`
	Bonus        float64 `json:"bonus"`
	CurrencyPair string  `json:"currency_pair"`
	Fee          float64 `json:"fee"`
	Price        float64 `json:"price"`
	Timestamp    string  `json:"timestamp"`
	YourAction   string  `json:"your_action"`
}

// GetUnixTimestamp is get unix timestamp
func (t TradeHistoryRecordResponse) GetUnixTimestamp() (int64, error) {
	return strconv.ParseInt(t.Timestamp, 10, 64)
}

// TradeHistory is get trade history
func (r *Requester) TradeHistory(tradeHistoryParams *TradeHistoryParams) (*TradeHistoryResponse, *utility.HTTPRequest, *http.Response, error) {
	params, err := query.Values(tradeHistoryParams)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, fmt.Sprintf("can not create request parameter of trade history (params = %v)", tradeHistoryParams))
	}
	for {
		request := r.makeTradeRequest("trade_history", params.Encode())
		newRes, response, err := r.unmarshal(func(request *utility.HTTPRequest) (interface{}, *http.Response, []byte, error) {
			res, resBody, err := r.httpClient.DoRequest(utility.HTTPMethdoPOST, request, false)
			if err != nil {
				return nil, res, resBody, errors.Wrap(err, fmt.Sprintf("can not get trade history (url = %v, params = %v)", request.URL, params.Encode()))
			}
			newRes := new(TradeHistoryResponse)
			return newRes, res, resBody, err
		}, request)
		if newRes.(*TradeHistoryResponse).needRetry() {
			continue
		}
		return newRes.(*TradeHistoryResponse), request, response, err
	}
}

// TradeActiveOrdersParams is parameter of active order
type TradeActiveOrderParams struct {
	CurrencyPair string `url:"currency_pair,omitempty"`
	IsToken      bool   `url:"is_token,omitempty"`
	IsTokenBoth  bool   `url:"is_token_both,omitempty"`
}

// TradeActiveOrderParams is create TradeActiveOrderParams
func (r *Requester) NewTradeActiveOrderParams() (*TradeActiveOrderParams) {
	return &TradeActiveOrderParams{}
}

// TradeActiveOrderResponse is response of active order
type TradeActiveOrderResponse struct {
	Return map[string]TradeActiveOrderRecordResponse `json:"return"`
	TradeCommonResponse
}

// TradeActiveOrderRecordResponse is response of active order
type TradeActiveOrderRecordResponse struct {
	Action       string  `json:"action"`
	Amount       float64 `json:"amount"`
	CurrencyPair string  `json:"currency_pair"`
	Price        float64 `json:"price"`
	Timestamp    string  `json:"timestamp"`
}

// TradeActiveOrderBothResponse is response of active order both
type TradeActiveOrderBothResponse struct {
	Return struct {
		ActiveOrders      map[string]TradeActiveOrderRecordResponse `json:"active_orders"`
		TokenActiveOrders map[string]TradeActiveOrderRecordResponse `json:"token_active_orders"`
	} `json:"return"`
	TradeCommonResponse
}

// TradeActiveOrder is get trade active order
func (r *Requester) TradeActiveOrder(tradeActiveOrderParams *TradeActiveOrderParams) (*TradeActiveOrderResponse, *utility.HTTPRequest, *http.Response, error) {
	tradeActiveOrderParams.IsTokenBoth = false
	params, err := query.Values(tradeActiveOrderParams)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, fmt.Sprintf("can not create request parameter of active order (params = %v)", tradeActiveOrderParams))
	}
	for {
		request := r.makeTradeRequest("active_orders", params.Encode())
		newRes, response, err := r.unmarshal(func(request *utility.HTTPRequest) (interface{}, *http.Response, []byte, error) {
			res, resBody, err := r.httpClient.DoRequest(utility.HTTPMethdoPOST, request, false)
			if err != nil {
				return nil, res, resBody, errors.Wrap(err, fmt.Sprintf("can not get active order (url = %v, params = %v)", request.URL, params.Encode()))
			}
			newRes := new(TradeActiveOrderResponse)
			return newRes, res, resBody, err
		}, request)
		if newRes.(*TradeActiveOrderResponse).needRetry() {
			continue
		}
		return newRes.(*TradeActiveOrderResponse), request, response, err
	}
}

// TradeActiveOrderBoth is get trade active order
func (r *Requester) TradeActiveOrderBoth(tradeActiveOrderParams *TradeActiveOrderParams) (*TradeActiveOrderBothResponse, *utility.HTTPRequest, *http.Response, error) {
	tradeActiveOrderParams.IsTokenBoth = true
	params, err := query.Values(tradeActiveOrderParams)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, fmt.Sprintf("can not create request parameter of active order (params = %v)", tradeActiveOrderParams))
	}
	for {
		request := r.makeTradeRequest("active_orders", params.Encode())
		newRes, response, err := r.unmarshal(func(request *utility.HTTPRequest) (interface{}, *http.Response, []byte, error) {
			res, resBody, err := r.httpClient.DoRequest(utility.HTTPMethdoPOST, request, false)
			if err != nil {
				return nil, res, resBody, errors.Wrap(err, fmt.Sprintf("can not get active order with both (url = %v, params = %v)", request.URL, params.Encode()))
			}
			newRes := new(TradeActiveOrderBothResponse)
			return newRes, res, resBody, err
		}, request)
		if newRes.(*TradeActiveOrderBothResponse).needRetry() {
			continue
		}
		return newRes.(*TradeActiveOrderBothResponse), request, response, err
	}
}

// TradeParams is parameter of trade
type TradeParams struct {
	CurrencyPair string  `url:"currency_pair"`
	Action       string  `url:"action"`
	Price        float64 `url:"price"`
	Amount       float64 `url:"amount"`
	Limit        float64 `url:"limit,omitempty"`
}

func (t *TradeParams) fixupPriceAndAmount(r *Requester) {
	priceUnit := r.GetMinPriceUnit(t.CurrencyPair)
	amountUnit := r.GetMinAmountUnit(t.CurrencyPair)
	fixedPrice := math.Floor((float64(int64(t.Price/priceUnit))*priceUnit)*100000000) / 100000000
	if fixedPrice != t.Price {
		if t.Action == "bid" {
			t.Price = math.Floor((fixedPrice+priceUnit)*100000000) / 100000000
		} else if t.Action == "ask" {
			t.Price = fixedPrice
		}
	}
	fixedAmount := math.Floor((float64(int64(t.Amount/amountUnit))*amountUnit)*10000) / 10000
	if fixedAmount != t.Amount {
		t.Amount = fixedAmount
	}
}

// NewTradeParams is create TradeParams
func (r *Requester) NewTradeParams() (*TradeParams) {
	return &TradeParams{}
}

// TradeResponse is response of trade
type TradeResponse struct {
	Return struct {
		Funds struct {
			Btc      float64 `json:"btc"`
			Bch      float64 `json:"bch"`
			Eth      float64 `json:"eth"`
			Mona     float64 `json:"mona"`
			Xem      float64 `json:"xem"`
			Jpy      float64 `json:"jpy"`
			Zaif     float64 `json:"zaif"`
			Pepecash float64 `json:"pepecash"`
		} `json:"funds"`
		OrderID  int64   `json:"order_id"`
		Received float64 `json:"received"`
		Remains  float64 `json:"remains"`
	} `json:"return"`
	TradeCommonResponse
}

func (r *Requester) tradeBase(tradeParams *TradeParams, retryCallback exchange.RetryCallback, retryCallbackData interface{}) (*TradeResponse, *utility.HTTPRequest, *http.Response, error) {
	for {
		tradeParams.fixupPriceAndAmount(r)
		params := make(url.Values)
		params.Add("currency_pair", tradeParams.CurrencyPair)
		params.Add("action", tradeParams.Action)
		params.Add("price", strconv.FormatFloat(tradeParams.Price, 'f', r.getPricePrec(tradeParams.CurrencyPair), 64))
		params.Add("amount", strconv.FormatFloat(tradeParams.Amount, 'f', r.getAmountPrec(tradeParams.CurrencyPair), 64))
		if tradeParams.Limit != 0 {
			params.Add("limit", strconv.FormatFloat(tradeParams.Limit, 'f', r.getPricePrec(tradeParams.CurrencyPair), 64))
		}
		request := r.makeTradeRequest("trade", params.Encode())
		log.Printf("try trade action = %v, currency pair = %v, price = %v, amount = %v, params = %v", tradeParams.Action, tradeParams.CurrencyPair, tradeParams.Price, tradeParams.Amount, params.Encode())
		newRes, response, err := r.unmarshal(func(request *utility.HTTPRequest) (interface{}, *http.Response, []byte, error) {
			res, resBody, err := r.httpClient.DoRequest(utility.HTTPMethdoPOST, request, true)
			if err != nil {
				return nil, res, resBody, errors.Wrap(err, fmt.Sprintf("can not trade (url = %v, params = %v)", request.URL, params.Encode()))
			}
			newRes := new(TradeResponse)
			return newRes, res, resBody, err
		}, request)
		if err != nil || newRes.(*TradeResponse).needRetry() {
			if (err != nil) {
				log.Printf("error occured currency pair = %v (%v)", tradeParams.CurrencyPair, err)
			}
			retry := retryCallback(&tradeParams.Price, &tradeParams.Amount, retryCallbackData)
			if !retry {
				return newRes.(*TradeResponse), request, response, err
			}
			log.Printf("retry trade action = v, currency pair = %v", tradeParams.Action, tradeParams.CurrencyPair)
			continue
		}
		log.Printf("trade done action = %v, currency pair = %v, price = %v, amount = %v", tradeParams.Action, tradeParams.CurrencyPair, tradeParams.Price, tradeParams.Amount)
		return newRes.(*TradeResponse), request, response, err
	}
}

// TradeBuy is buy trade
func (r *Requester) TradeBuy(tradeParams *TradeParams, retryCallback exchange.RetryCallback, retryCallbackData interface{}) (*TradeResponse, *utility.HTTPRequest, *http.Response, error) {
	tradeParams.Action = "bid"
	return r.tradeBase(tradeParams, retryCallback, retryCallbackData)
}

// TradeSell is sell trade
func (r *Requester) TradeSell(tradeParams *TradeParams, retryCallback exchange.RetryCallback, retryCallbackData interface{}) (*TradeResponse, *utility.HTTPRequest, *http.Response, error) {
	tradeParams.Action = "ask"
	return r.tradeBase(tradeParams, retryCallback, retryCallbackData)
}

// TradeCancelOrderParams is parameter of cancel order
type TradeCancelOrderParams struct {
	OrderId int64 `url:"order_id"`
	IsToken bool  `url:"is_token,omitempty"`
}

// NewTradeCancelOrderParams is create TradeCancelOrderParams
func (r *Requester) NewTradeCancelOrderParams() (*TradeCancelOrderParams) {
	return &TradeCancelOrderParams{}
}

// TradeCancelOrderResponse is response of calcel order
type TradeCancelOrderResponse struct {
	Return struct {
		Funds struct {
			Btc      float64 `json:"btc"`
			Bch      float64 `json:"bch"`
			Eth      float64 `json:"eth"`
			Mona     float64 `json:"mona"`
			Xem      float64 `json:"xem"`
			Jpy      float64 `json:"jpy"`
			Zaif     float64 `json:"zaif"`
			Pepecash float64 `json:"pepecash"`
		} `json:"funds"`
		OrderID int64 `json:"order_id"`
	} `json:"return"`
	TradeCommonResponse
}

// TradeCancelOrder is cancel order
func (r *Requester) TradeCancelOrder(tradeCancelOrderParams *TradeCancelOrderParams) (*TradeCancelOrderResponse, *utility.HTTPRequest, *http.Response, error) {
	params, err := query.Values(tradeCancelOrderParams)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, fmt.Sprintf("can not create request parameter of cancel order (params = %v)", tradeCancelOrderParams))
	}
	for {
		request := r.makeTradeRequest("cancel_order", params.Encode())
		newRes, response, err := r.unmarshal(func(request *utility.HTTPRequest) (interface{}, *http.Response, []byte, error) {
			res, resBody, err := r.httpClient.DoRequest(utility.HTTPMethdoPOST, request, false)
			if err != nil {
				return nil, res, resBody, errors.Wrap(err, fmt.Sprintf("can not cancel order (url = %v, params = %v)", request.URL, params.Encode()))
			}
			newRes := new(TradeCancelOrderResponse)
			return newRes, res, resBody, err
		}, request)
		if newRes.(*TradeCancelOrderResponse).needRetry() {
			continue
		}
		return newRes.(*TradeCancelOrderResponse), request, response, err
	}
}

// TradeWithdrawParams is parameter of withdraw
type TradeWithdrawParams struct {
	Currency string  `url:"currency"`
	Address  string  `url:"address"`
	Message  string  `url:"message,omitempty"`
	Amount   float64 `url:"amount"`
	OptFee   float64 `url:"opt_fee,omitempty"`
}

func (t *TradeWithdrawParams) fixupFee(r *Requester) {
	switch t.Currency {
	case "btc":
		if t.OptFee < r.GetMinFee(t.Currency) {
			t.OptFee = r.GetMinFee(t.Currency)
		}
	case "xem":
		if t.OptFee < r.GetMinFee(t.Currency) {
			t.OptFee = r.GetMinFee(t.Currency)
		}
	default:
		t.OptFee = r.GetMinFee(t.Currency)
	}
}

// NewTradeWithdrawParams is create TradeWithdrawParams
func (r *Requester) NewTradeWithdrawParams() (*TradeWithdrawParams) {
	return &TradeWithdrawParams{}
}

// TradeWithdrawResponse is response of Withdraw
type TradeWithdrawResponse struct {
	Return struct {
		Funds struct {
			Btc      float64 `json:"btc"`
			Bch      float64 `json:"bch"`
			Eth      float64 `json:"eth"`
			Mona     float64 `json:"mona"`
			Xem      float64 `json:"xem"`
			Jpy      float64 `json:"jpy"`
			Zaif     float64 `json:"zaif"`
			Pepecash float64 `json:"pepecash"`
		} `json:"funds"`
		Fee  float64 `json:"fee"`
		TxID string  `json:"txid"`
	} `json:"return"`
	TradeCommonResponse
}

// TradeCancelOrder is cancel order
func (r *Requester) TradeWithdraw(tradeWithdrawParams *TradeWithdrawParams) (*TradeWithdrawResponse, *utility.HTTPRequest, *http.Response, error) {
	tradeWithdrawParams.fixupFee(r)
	params, err := query.Values(tradeWithdrawParams)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, fmt.Sprintf("can not create request parameter of Withdraw (params = %v)", tradeWithdrawParams))
	}
	for {
		request := r.makeTradeRequest("withdraw", params.Encode())
		newRes, response, err := r.unmarshal(func(request *utility.HTTPRequest) (interface{}, *http.Response, []byte, error) {
			res, resBody, err := r.httpClient.DoRequest(utility.HTTPMethdoPOST, request, false)
			if err != nil {
				return nil, res, resBody, errors.Wrap(err, fmt.Sprintf("can not Withdraw (url = %v, params = %v)", request.URL, params.Encode()))
			}
			newRes := new(TradeWithdrawResponse)
			return newRes, res, resBody, err
		}, request)
		if newRes.(*TradeWithdrawResponse).needRetry() {
			continue
		}
		return newRes.(*TradeWithdrawResponse), request, response, err
	}
}

// TradeDepositHistoryParams is parameter of trade history
type TradeDepositHistoryParams struct {
	Currency string `url:"currency"`
	From     int64  `url:"from,omitempty"`
	Count    int64  `url:"count,omitempty"`
	FromID   int64  `url:"from_id,omitempty"`
	EndID    int64  `url:"end_id"`
	Order    string `url:"order,omitempty"`
	Since    int64  `url:"since,omitempty"`
	End      int64  `url:"end"`
}

// NewTradeDepositHistoryParams is create TradeDepositHistoryParams
func (r *Requester) NewTradeDepositHistoryParams() (*TradeDepositHistoryParams) {
	return &TradeDepositHistoryParams{
		EndID: math.MaxInt64,
		End:   time.Now().Unix(),
	}
}

// TradeDepositHistoryResponse is response of deposit history
type TradeDepositHistoryResponse struct {
	Return map[string]TradeDepositHistoryRecordResponse `json:"return"`
	TradeCommonResponse
}

// TradeDepositHistoryRecordResponse is response of deposit history record
type TradeDepositHistoryRecordResponse struct {
	Address   string  `json:"address"`
	Amount    float64 `json:"amount"`
	Timestamp string  `json:"timestamp"`
	TxID      string  `json:"txid"`
}

// TradeDepositHistory is deposit history
func (r *Requester) TradeDepositHistory(tradeDepositHistoryParams *TradeDepositHistoryParams) (*TradeDepositHistoryResponse, *utility.HTTPRequest, *http.Response, error) {
	params, err := query.Values(tradeDepositHistoryParams)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, fmt.Sprintf("can not create request parameter of deposit history (params = %v)", tradeDepositHistoryParams))
	}
	for {
		request := r.makeTradeRequest("deposit_history", params.Encode())
		newRes, response, err := r.unmarshal(func(request *utility.HTTPRequest) (interface{}, *http.Response, []byte, error) {
			res, resBody, err := r.httpClient.DoRequest(utility.HTTPMethdoPOST, request, false)
			if err != nil {
				return nil, res, resBody, errors.Wrap(err, fmt.Sprintf("can not get deposit history (url = %v, params = %v)", request.URL, params.Encode()))
			}
			newRes := new(TradeDepositHistoryResponse)
			return newRes, res, resBody, err
		}, request)
		if newRes.(*TradeDepositHistoryResponse).needRetry() {
			continue
		}
		return newRes.(*TradeDepositHistoryResponse), request, response, err
	}
}

// TradeWithdrawHistoryParams is parameter of trade history
type TradeWithdrawHistoryParams struct {
	Currency string `url:"currency"`
	From     int64  `url:"from,omitempty"`
	Count    int64  `url:"count,omitempty"`
	FromID   int64  `url:"from_id,omitempty"`
	EndID    int64  `url:"end_id"`
	Order    string `url:"order,omitempty"`
	Since    int64  `url:"since,omitempty"`
	End      int64  `url:"end"`
}

// NewTradeWithdrawHistoryParams is create TradeWithdrawHistoryParams
func (r *Requester) NewTradeWithdrawHistoryParams() (*TradeWithdrawHistoryParams) {
	return &TradeWithdrawHistoryParams{
		EndID: math.MaxInt64,
		End:   time.Now().Unix(),
	}
}

// TradeWithdrawHistoryResponse is response of withdraw history
type TradeWithdrawHistoryResponse struct {
	Return map[string]TradeWithdrawHistoryRecordResponse `json:"return"`
	TradeCommonResponse
}

// TradeWithdrawHistoryRecordResponse is response withdraw history record
type TradeWithdrawHistoryRecordResponse struct {
	Address   string  `json:"address"`
	Amount    float64 `json:"amount"`
	Timestamp string  `json:"timestamp"`
	TxID      string  `json:"txid"`
}

// TradeWithdrawHistory is withdraw history
func (r *Requester) TradeWithdrawHistory(tradeWithdrawHistoryParams *TradeWithdrawHistoryParams) (*TradeWithdrawHistoryResponse, *utility.HTTPRequest, *http.Response, error) {
	params, err := query.Values(tradeWithdrawHistoryParams)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, fmt.Sprintf("can not create request parameter of withdraw history (params = %v)", tradeWithdrawHistoryParams))
	}
	for {
		request := r.makeTradeRequest("withdraw_history", params.Encode())
		newRes, response, err := r.unmarshal(func(request *utility.HTTPRequest) (interface{}, *http.Response, []byte, error) {
			res, resBody, err := r.httpClient.DoRequest(utility.HTTPMethdoPOST, request, false)
			if err != nil {
				return nil, res, resBody, errors.Wrap(err, fmt.Sprintf("can not get withdraw history (url = %v, params = %v)", request.URL, params.Encode()))
			}
			newRes := new(TradeWithdrawHistoryResponse)
			return newRes, res, resBody, err
		}, request)
		if newRes.(*TradeWithdrawHistoryResponse).needRetry() {
			continue
		}
		return newRes.(*TradeWithdrawHistoryResponse), request, response, err
	}
}
