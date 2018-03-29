package fetcher

import (
	"sync"
	"log"
	"sync/atomic"
	"time"
	"path"
	"github.com/AutomaticCoinTrader/ACT/utility"
	"github.com/AutomaticCoinTrader/ACT/tools/proxy/zaif/server"
	"github.com/AutomaticCoinTrader/ACT/tools/proxy/zaif/configurator"
	"github.com/AutomaticCoinTrader/ACT/exchange/zaif"
	"net"
	"github.com/pkg/errors"
)



type currencyPairsInfo struct {
	Bids      map[string][][]float64
	Asks      map[string][][]float64
	LastPrice map[string]float64
	mutex     *sync.Mutex
}

func (c *currencyPairsInfo) updateDepth(currencyPair string, currencyPairsBids [][]float64, currencyPairsAsks [][]float64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Bids[currencyPair] = currencyPairsBids
	c.Asks[currencyPair] = currencyPairsAsks
}

type Fetcher struct {
	config              *configurator.ZaifProxyConfig
	requester           *zaif.Requester
	httpClients         []*utility.HTTPClient
	httpClientsIdx      int
	httpClientsMutex    *sync.Mutex
	pollingFinish       int32
	currencyPairsInfo   *currencyPairsInfo
	websocketServer     *server.WebsocketServer
	pausePollingRequest int32
}

func (f *Fetcher) pollingLoop(pollingRequestChan chan string, lastBidsMap map[string][][]float64, lastAsksMap map[string][][]float64, lastBidsAsksMutex *sync.Mutex) {
	log.Printf("start polling loop")
	for {
		currencyPair, ok := <-pollingRequestChan
		if !ok {
			log.Printf("finish polling loop")
			return
		}
		// select httpClient
		f.httpClientsMutex.Lock()
		httpClient := f.httpClients[f.httpClientsIdx]
		f.httpClientsIdx += 1
		if f.httpClientsIdx >= len(f.httpClients) {
			f.httpClientsIdx = 0
		}
		f.httpClientsMutex.Unlock()
		request := f.requester.MakePublicRequest(path.Join("depth", currencyPair), "")
		res, resBody, err := httpClient.DoRequest(utility.HTTPMethodGET, request, true)
		if err != nil {
			log.Printf("can not get depcth (url = %v)", request.URL)
			if res != nil && res.StatusCode == 403 {
				log.Printf("occured 403 Forbidden currency pair = %v", currencyPair)
				atomic.StoreInt32(&f.pausePollingRequest, 1)
			}
			continue
		}
		f.websocketServer.BroadCast(currencyPair, resBody)
	}
}

func (f *Fetcher) pollingRequestLoop() {
	log.Printf("start polling request loop")
	atomic.StoreInt32(&f.pollingFinish, 0)
	lastBidsMap := make(map[string][][]float64)
	lastAsksMap := make(map[string][][]float64)
	lastBidsAsksMutex := new(sync.Mutex)
	pollingRequestChan := make(chan string)
	for i := 0; i < len(f.config.CurrencyPairs) * 2; i++ {
		go f.pollingLoop(pollingRequestChan, lastBidsMap, lastAsksMap, lastBidsAsksMutex)
	}
	fetchCount := uint64(0)
FINISH:
	for {
		if fetchCount % 20 == 0 {
			log.Printf("start get depth of currency Pairs (fetch count = %v, time = %v)", fetchCount, time.Now().UnixNano())
		}
		for _, currencyPair := range f.config.CurrencyPairs {
			if atomic.LoadInt32(&f.pollingFinish) == 1 {
				break FINISH
			}
			if atomic.LoadInt32(&f.pausePollingRequest) == 1 {
				time.Sleep(time.Duration(f.config.PauseWait) * time.Second)
				atomic.StoreInt32(&f.pausePollingRequest, 0)
			}
			pollingRequestChan <- currencyPair
			time.Sleep(time.Duration(f.config.PollingWait) * time.Millisecond)
		}
		fetchCount += 1
	}
	close(pollingRequestChan)
	log.Printf("finish polling request loop")
}

func (f *Fetcher) Start() {
	go f.pollingRequestLoop()
	f.websocketServer.Start()
}

func (f *Fetcher) Stop() {
	f.websocketServer.Stop()
	atomic.StoreInt32(&f.pollingFinish, 1)
}

func NewFetcher(config *configurator.ZaifProxyConfig) (*Fetcher, error) {
	dummyRequesterKeys := make([]*zaif.RequesterKey, 0)

	httpClients := make([]*utility.HTTPClient, 0)
	for _, clientBindAddress := range config.ClientBindAddresses {
		localAddr, err := net.ResolveIPAddr("ip", clientBindAddress)
		if err != nil {
			return nil, errors.Wrap(err, "can not resolve ip address")
		}
		httpClient := utility.NewHTTPClient(config.Retry, config.RetryWait, config.Timeout, localAddr)
		httpClients = append(httpClients, httpClient)
	}
	return &Fetcher{
		requester:     zaif.NewRequester(dummyRequesterKeys, config.Retry, config.RetryWait, config.Timeout, config.ReadBufSize, config.WriteBufSize),

		httpClients:   httpClients,
		httpClientsIdx: 0,
		httpClientsMutex: new(sync.Mutex),

		config:        config,
		pollingFinish: 0,
		currencyPairsInfo: &currencyPairsInfo{
			Bids:      make(map[string][][]float64),
			Asks:      make(map[string][][]float64),
			LastPrice: make(map[string]float64),
			mutex:     new(sync.Mutex),
		},
		websocketServer: server.NewWsServer(config),
		pausePollingRequest: 0,
	}, nil
}
