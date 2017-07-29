package coincheck

import (
	"github.com/AutomaticCoinTrader/ACT/utility"
	"time"
	"crypto/hmac"
	"crypto/sha256"
	"log"
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

func (cr *CoincheckRequester) StreamingStart() error {
	log.Println("Coincheck StreamingStart")
	return nil
}

func (cr *CoincheckRequester) StreamingStop() error {
	log.Println("Coincheck StreamingStop")
	return nil
}

func NewRequester() *CoincheckRequester {
	return &CoincheckRequester {
		httpClient: nil,
		apiKey: "",
		apiSecret: "",
	}
}