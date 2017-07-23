package coincheck

import (
	"github.com/AutomaticCoinTrader/ACT/utility"
	"time"
	"crypto/hmac"
	"crypto/sha256"
)

const (
	ENDPOINT = "coincheck.com"
	ENDPOINT_WS = "ws-api.coincheck.com"
)

type CoincheckRequester struct {
	httpClient *utility.HTTPClient
	apiKey string
	apiSecret string
}

func (cr *CoincheckRequester) sign(req *utility.HTTPRequest) error {
	nonce := time.Now().Unix()
	message := string(nonce) + req.URL + req.Body
	m := hmac.New(sha256.New, []byte(cr.apiSecret))
	m.Write([]byte(message))
	signature := string(m.Sum(nil))

	req.Headers["ACCESS-KEY"] = cr.apiKey
	req.Headers["ACCESS-NONCE"] = string(nonce)
	req.Headers["ACCESS-SIGNATURE"] = signature

	return nil
}

func (cr *CoincheckRequester) GetBoard() (error) {

	headers := make(map[string]string)

	req := &utility.HTTPRequest{
		URL: "https://" + ENDPOINT + "/api/order_books",
		Headers: headers,
		Body: "",
	}

	cr.sign(req)
	resp, body, err := cr.httpClient.DoRequest(utility.HTTPMethodGET, req)
	if err != nil {
		panic("failed")
	}

	println(resp, body)

	return nil
}

func (cr *CoincheckRequester) GetBalance() error {
	return nil
}

func NewRequester() *CoincheckRequester {
	return &CoincheckRequester {
		httpClient: nil,
		apiKey: "",
		apiSecret: "",
	}
}