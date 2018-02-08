package coincheck
//
//import (
//	"encoding/json"
//	"fmt"
//	"log"
//
//	"github.com/AutomaticCoinTrader/ACT/utility"
//	"github.com/gorilla/websocket"
//	"github.com/pkg/errors"
//)
//
//// Public
//// - GET /api/ticker
//// - GET /api/trades
//// - GET /api/order_books
//// - GET /api/exchange/orders/rate
//// - GET /api/rate/<pair>
////
//// Private
//// - POST /api/exchange/orders
//// - GET /api/exchange/orders/opens
//// - DELETE /api/exchange/orders/<id>
//// - GET /api/exchange/orders/transactions
//// - GET /api/exchange/orders/transactions_pagination
//// - GET /api/exchange/leverage/positions
//// - GET /api/accounts/balance
//// - GET /api/accounts/leverage_balance
//// - POST /api/send_money
//// - GET /api/send_money
//// - GET /api/deposit_money
//// - POST /api/deposit_money/[id]/fast
//// - GET /api/accounts
//// - GET /api/bank_accounts
//// - POST /api/bank_accounts
//// - DELETE /api/bank_accounts/[id]
//// ...
////
//// Streaming
//// - {"type": "subscribe","channel": "[pair]-trades"}
//// - {"type": "subscribe","channel": "[pair]-orderbook"}
//
//// GetBalance ...
//func (cr *CoincheckRequester) GetBalance() error {
//	return nil
//}
//
//// GetBoard ...
//func (cr *CoincheckRequester) GetBoard() error {
//
//	headers := make(map[string]string)
//
//	req := &utility.HTTPRequest{
//		URL:     "https://" + Endpoint + "/api/order_books",
//		Headers: headers,
//		Body:    "",
//	}
//
//	cr.sign(req)
//	resp, body, err := cr.httpClient.DoRequest(utility.HTTPMethodGET, req)
//	if err != nil {
//		panic("failed")
//	}
//
//	println(resp, body)
//
//	return nil
//}
//
//type StreamingCallback func(pair string, response []interface{}, callbackdata interface{}) error
//
//type coincheckStreamingCallbackData struct {
//	currencyPair             string
//	tradeHistoryCallback     StreamingCallback
//	tradeHistoryCallbackData interface{}
//	boardUpdateCallback      StreamingCallback
//	boardUpdateCallbackData  interface{}
//}
//
//func (cr *CoincheckRequester) streamingCallback(conn *websocket.Conn, userCallbackData interface{}) error {
//	streaminCallbackData := userCallbackData.(*coincheckStreamingCallbackData)
//	messageType, message, err := conn.ReadMessage()
//	if err != nil {
//		return errors.Wrap(err, "can not read message of streaming")
//	}
//	if messageType != websocket.TextMessage {
//		log.Printf("unsupported message type (message type = %v, message = %v)", messageType, message)
//		return nil
//	}
//	newRes := make([]interface{}, 5)
//	err = json.Unmarshal(message, &newRes)
//	if err != nil {
//		return errors.Wrap(err, fmt.Sprintf("can not unmarshal message of streaming"))
//	}
//	err = streaminCallbackData.tradeHistoryCallback(streaminCallbackData.currencyPair, newRes, streaminCallbackData.tradeHistoryCallbackData)
//	if err != nil {
//		return errors.Wrap(err, fmt.Sprintf("call back error of streaming"))
//	}
//	return nil
//}
//
//// StreamingStart ...
//func (cr *CoincheckRequester) StreamingStart(_ string, callback StreamingCallback, callbackData interface{}) error {
//	requestURL := fmt.Sprintf("wss://%s", WebsocketEndpoint)
//	streaminCallbackData := &coincheckStreamingCallbackData{
//		currencyPair:             "btc_jpy",
//		tradeHistoryCallback:     callback,
//		tradeHistoryCallbackData: callbackData,
//	}
//	newClient := utility.NewWSClient(128*1024, 128*1024, 5)
//	err := newClient.Start(cr.streamingCallback, streaminCallbackData, requestURL, nil)
//	if err != nil {
//		return errors.Wrap(err, fmt.Sprintf("can not start streaming (url = %v)", requestURL))
//	}
//	newClient.Send(struct {
//		APIType string `json:"type"`
//		Channel string `json:"channel"`
//	}{
//		APIType: "subscribe", Channel: "btc_jpy-trades",
//	})
//	cr.websocketClient = append(cr.websocketClient, newClient)
//	return nil
//}
//
//// StreamingStop ...
//func (cr *CoincheckRequester) StreamingStop(_ string) {
//	log.Println("CoincheckRequester StreamingStop")
//	for _, ws := range cr.websocketClient {
//		ws.Stop()
//	}
//}
