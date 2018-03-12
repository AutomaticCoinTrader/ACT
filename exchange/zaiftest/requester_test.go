package zaiftest

import (
	"testing"
	"path"
	"log"
	"github.com/AutomaticCoinTrader/ACT/integrator"
	"github.com/AutomaticCoinTrader/ACT/configurator"
	"github.com/AutomaticCoinTrader/ACT/utility"
	"github.com/AutomaticCoinTrader/ACT/exchange/zaif"
	"net/http"
	"encoding/json"
	"fmt"
	"time"
)

func dump(t *testing.T, res interface{}, httpReq *utility.HTTPRequest, httpRes *http.Response, ) {
	fmt.Printf("========================\n")
	fmt.Printf("url: %v, method = %v\n", httpReq.URL, httpReq.RequestMethodString)
	bytes, err := json.Marshal(httpReq.Headers)
	if err != nil {
		t.Fatalf("can not marshal json")
	}
	fmt.Printf("headers: %v\n", string(bytes))
	fmt.Printf("body: %v\n", httpReq.Body)
	fmt.Printf("status; %v\n", httpRes.Status)
	fmt.Printf("length: %v\n", httpRes.ContentLength)
	bytes, err = json.Marshal(res)
	if err != nil {
		t.Fatalf("can not marshal json")
	}
	fmt.Printf("response: %v\n", string(bytes))
}

func callback(currencyPair string, streamingResponse *zaif.StreamingResponse, streamingCallbackData interface{}) (error) {
	fmt.Printf("--------------------------\n");
	fmt.Printf("currencyPair: %v\n", currencyPair)
	t := streamingCallbackData.(*testing.T)
	bytes, err := json.Marshal(streamingResponse)
	if err != nil {
		t.Fatalf("can not marshal json")
	}
	fmt.Printf("response: %v\n", string(bytes))
	return nil
}

func retryCallback(price *float64, amount *float64, errMsg string, retryCallbackData interface{}) (bool) {
	return true
}

func TestRequester(t *testing.T) {
	cf, err := configurator.NewConfigurator(path.Join("../../config", "act"))
	if err != nil {
		log.Printf("can not create configurator (config dir = %v, reason = %v)", "../../config", err)
		t.Fatalf("can not create configurator")
	}
	newConfig := new(integrator.Config)
	err = cf.Load(newConfig)
	if err != nil {
		log.Printf("can not load config (config dir = %v, reason = %v)", "../../config", err)
		t.Fatalf("can not load config")
	}
	fmt.Printf("keys = %v\n", newConfig.Exchanges.Zaif.Keys)
	requesterKeys := make([]*zaif.RequesterKey, 0, len(newConfig.Exchanges.Zaif.Keys))
	for _, key := range newConfig.Exchanges.Zaif.Keys {
		requesterKeys = append(requesterKeys, &zaif.RequesterKey{Key : key.Key, Secret:key.Secret})
	}
	r := zaif.NewRequester(requesterKeys, 10, 500, 600, 10*1024*1024, 10*1024*1024)

	res1, httpreq, httpres, err := r.Currencies("all")
	if (err != nil) {
		t.Fatalf("Currencies failure")
	}
	dump(t, res1, httpreq, httpres)

	res2, httpreq, httpres, err := r.CurrencyPairs("all")
	if (err != nil) {
		t.Fatalf("CurrencyPairs failure")
	}
	dump(t, res2, httpreq, httpres)

	res3, httpreq, httpres, err := r.LastPrice("btc_jpy")
	if (err != nil) {
		t.Fatalf("LastPrice failure")
	}
	dump(t, res3, httpreq, httpres)

	res4, httpreq, httpres, err := r.Ticker("btc_jpy")
	if (err != nil) {
		t.Fatalf("Ticker failure")
	}
	dump(t, res4, httpreq, httpres)

	res5, httpreq, httpres, err := r.Trades("btc_jpy")
	if (err != nil) {
		t.Fatalf("Trades failure")
	}
	dump(t, res5, httpreq, httpres)

	res6, httpreq, httpres, err := r.Depth("btc_jpy")
	if (err != nil) {
		t.Fatalf("Depth failure")
	}
	dump(t, res6, httpreq, httpres)

	res7, httpreq, httpres, err := r.GetInfo()
	if (err != nil) {
		t.Fatalf("GetInfo failure")
	}
	dump(t, res7, httpreq, httpres)

	res8, httpreq, httpres, err := r.GetInfo2()
	if (err != nil) {
		t.Fatalf("GetInfo2 failure")
	}
	dump(t, res8, httpreq, httpres)

	res9, httpreq, httpres, err := r.GetPersonalInfo()
	if (err != nil) {
		t.Fatalf("GetPersonalInfo failure")
	}
	dump(t, res9, httpreq, httpres)

	res10, httpreq, httpres, err := r.GetIDInfo()
	if (err != nil) {
		t.Fatalf("GetIDInfo failure")
	}
	dump(t, res10, httpreq, httpres)

	param1 := r.NewTradeHistoryParams()
	param1.Count = 10
	res11, httpreq, httpres, err := r.TradeHistory(param1)
	if (err != nil) {
		t.Fatalf("TradeHistory failure")
	}
	dump(t, res11, httpreq, httpres)

	param2 := r.NewTradeActiveOrderParams()
	res12, httpreq, httpres, err := r.TradeActiveOrder(param2)
	if (err != nil) {
		t.Fatalf("TradeActiveOrder failure")
	}
	dump(t, res12, httpreq, httpres)

	param3 := r.NewTradeActiveOrderParams()
	res13, httpreq, httpres, err := r.TradeActiveOrderBoth(param3)
	if (err != nil) {
		t.Fatalf("TradeActiveOrderBoth failure")
	}
	dump(t, res13, httpreq, httpres)

	param4 := r.NewTradeParams()
	param4.CurrencyPair = "zaif_jpy"
	param4.Price = 0.1
	param4.Amount = 1
	res14, httpreq, httpres, err := r.TradeBuy(param4, retryCallback, nil)
	if (err != nil) {
		t.Fatalf("TradeBuy failure")
	}
	dump(t, res14, httpreq, httpres)

	param5 := r.NewTradeCancelOrderParams()
	param5.OrderId = res14.Return.OrderID
	param5.IsToken = true
	res15, httpreq, httpres, err := r.TradeCancelOrder(param5)
	if (err != nil) {
		t.Fatalf("TradeCancelOrder failure")
	}
	dump(t, res15, httpreq, httpres)

	param6 := r.NewTradeParams()
	param6.CurrencyPair = "xem_btc"
	param6.Price = 1.3333e-07
	param6.Amount = 1
	res16, httpreq, httpres, err := r.TradeBuy(param6, retryCallback, nil)
	if (err != nil) {
		t.Fatalf("TradeBuy failure")
	}
	dump(t, res16, httpreq, httpres)

	param7 := r.NewTradeCancelOrderParams()
	param7.OrderId = res16.Return.OrderID
	param7.IsToken = true
	res17, httpreq, httpres, err := r.TradeCancelOrder(param7)
	if (err != nil) {
		t.Fatalf("TradeCancelOrder failure")
	}
	dump(t, res17, httpreq, httpres)


	r.StreamingStart("btc_jpy", callback, t)

	time.Sleep(time.Duration(60 * time.Second))

	r.StreamingStop("btc_jpy")
}
