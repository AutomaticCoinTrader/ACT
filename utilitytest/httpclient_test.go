package utilitytest

import (
	"testing"
	"github.com/AutomaticCoinTrader/ACT/utility"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"
)

func TestHttpClient(t *testing.T) {
	httpClient := utility.NewHTTPClient(3, 1, 1, nil)
	headers := make(map[string]string)
	headers["Connection"] = "keep-alive"
	request := &utility.HTTPRequest{
		Headers: headers,
		URL:   "http://www.google.com/",
	}
	// 1
	response, resBody, err := httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

	// 2
	response, resBody, err = httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

	// 3
	response, resBody, err = httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

	// 4
	response, resBody, err = httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

	// 5
	response, resBody, err = httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

	// 6
	response, resBody, err = httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

	// 7
	response, resBody, err = httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

	// 8
	response, resBody, err = httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

}

func TestHttpsClient(t *testing.T) {
	httpClient := utility.NewHTTPClient(3, 1, 1, nil)
	headers := make(map[string]string)
	headers["Connection"] = "keep-alive"
	request := &utility.HTTPRequest{
		Headers: headers,
		URL:   "https://www.google.com/",
	}
	// 1
	response, resBody, err := httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

	// 2
	response, resBody, err = httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

	// 3
	response, resBody, err = httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

	// 4
	response, resBody, err = httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

	// 5
	response, resBody, err = httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

	// 6
	response, resBody, err = httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

	// 7
	response, resBody, err = httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

	// 8
	response, resBody, err = httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

}

func TestHttpClient2(t *testing.T) {
	httpClient := utility.NewHTTPClient(3, 1, 1, nil)
	headers := make(map[string]string)
	headers["Connection"] = "keep-alive"
	request := &utility.HTTPRequest{
		Headers: headers,
		URL:   "http://www.yahoo.co.jp/",
	}
	// 1
	response, resBody, err := httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

	// 2
	response, resBody, err = httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

	// 3
	response, resBody, err = httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

	// 4
	response, resBody, err = httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

	// 5
	response, resBody, err = httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

	// 6
	response, resBody, err = httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

	// 7
	response, resBody, err = httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

	// 8
	response, resBody, err = httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

}

func TestHttpsClient2(t *testing.T) {
	httpClient := utility.NewHTTPClient(3, 1, 1, nil)
	headers := make(map[string]string)
	headers["Connection"] = "keep-alive"
	request := &utility.HTTPRequest{
		Headers: headers,
		URL:   "https://www.yahoo.co.jp/",
	}
	// 1
	response, resBody, err := httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

	// 2
	response, resBody, err = httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

	// 3
	response, resBody, err = httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

	// 4
	response, resBody, err = httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

	// 5
	response, resBody, err = httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

	// 6
	response, resBody, err = httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

	// 7
	response, resBody, err = httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

	// 8
	response, resBody, err = httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Print(response.StatusCode)
	fmt.Print(string(resBody))

}

var sampleHandler = http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Connection", "keep-alive")
	fmt.Fprintf(w, "test")
})

func TestHttpClientConcurrent(t *testing.T) {
	ts := httptest.NewServer(sampleHandler)
	defer ts.Close()
	httpClient := utility.NewHTTPClient(3, 1, 1, nil)
	finishChans := make([]chan bool, 0)
	t1 := time.Now()
	for i := 0; i < 5000; i++ {
		finishCh := make(chan bool)
		finishChans = append(finishChans, finishCh)
		go concurrencyRequest(ts, httpClient, t, i, finishCh)
	}
	for _, finishChan := range finishChans {
		<-finishChan
	}
	t2 := time.Now()
	fmt.Printf("elapsed = %v", t2.UnixNano() - t1.UnixNano())
}

func concurrencyRequest(ts *httptest.Server, httpClient *utility.HTTPClient, t *testing.T, idx int, finishChan chan bool) {
	headers := make(map[string]string)
	headers["Connection"] = "keep-alive"
	headers["X-idx"] = fmt.Sprintf("%v", idx)
	request := &utility.HTTPRequest{
		Headers: headers,
		URL:     ts.URL,
	}
	_, _, err := httpClient.DoRequest(utility.HTTPMethodGET, request, false)
	if err != nil {
		t.Fatalf("request failure")
	}
	fmt.Printf("finish = %v\n", idx)
	close(finishChan)
}
