package coincheck

import (
	"crypto/hmac"
	"crypto/sha256"
	"time"

	"github.com/AutomaticCoinTrader/ACT/utility"
)

const (
	// Endpoint ...
	Endpoint = "coincheck.com"
	// WebsocketEndpoint ...
	WebsocketEndpoint = "ws-api.coincheck.com"
)

// CoincheckRequester ...
type CoincheckRequester struct {
	httpClient      *utility.HTTPClient
	websocketClient []*utility.WSClient
	apiKey          string
	apiSecret       string
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

// NewCoincheckRequester ...
func NewCoincheckRequester(key string, secret string) *CoincheckRequester {
	return &CoincheckRequester{
		httpClient:      nil,
		websocketClient: make([]*utility.WSClient, 0),
		apiKey:          key,
		apiSecret:       secret,
	}
}
