package notifier

import (
	"errors"
	"fmt"

	"github.com/AutomaticCoinTrader/ACT/utility"
)

const (
	EventName = "act_notify"
)

type IFTTTNotifierConfig struct {
	Key string `json:"key" yaml:"key" toml:"key"`
}

type IFTTTNotifier struct {
	Key        string
	httpClient *utility.HTTPClient
}

// Notify : Send an IFTTT notification via Web Request
func (n *IFTTTNotifier) Notify(msg string) error {

	req := &utility.HTTPRequest{
		URL: fmt.Sprintf("https://maker.ifttt.com/trigger/%s/with/key/%s", EventName, n.Key),
	}

	resp, _, err := n.httpClient.DoRequest(utility.HTTPMethodGET, req, false)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("Unexpected status code %d in IFTTTNotifier", resp.StatusCode))
	}

	return nil
}

func NewIFTTTNotifier(config *IFTTTNotifierConfig) (*IFTTTNotifier, error) {
	return &IFTTTNotifier{
		Key:        config.Key,
		httpClient: utility.NewHTTPClient(10, 1000, 60, nil),
	}, nil
}
