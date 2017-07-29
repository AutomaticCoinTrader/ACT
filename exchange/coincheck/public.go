package coincheck

import "github.com/AutomaticCoinTrader/ACT/utility"

func (cr *CoincheckRequester) GetBalance() error {
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


