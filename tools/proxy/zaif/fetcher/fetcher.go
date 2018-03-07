package fetcher

import (
	"sync"
	"log"
	"reflect"
	"sync/atomic"
	"time"
	"github.com/AutomaticCoinTrader/ACT/tools/proxy/zaif/configurator"
	"github.com/AutomaticCoinTrader/ACT/exchange/zaif"
)

const (
	defaultPollingConcurrency = 4
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

type fetcher struct {
	config *configurator.ZaifProxyConfig
	requester *zaif.Requester
	pollingFinish     int32
	currencyPairsInfo *currencyPairsInfo
}

func  (f *fetcher) pollingLoop(pollingRequestChan chan string, lastBidsMap map[string][][]float64, lastAsksMap map[string][][]float64, lastBidsAsksMutex *sync.Mutex) {
	log.Printf("start polling loop")
	for {
		currencyPair, ok := <- pollingRequestChan
		if !ok {
			log.Printf("finish polling loop")
			return
		}
		lastBidsAsksMutex.Lock()
		lastBids, bidsOk := lastBidsMap[currencyPair]
		lastAsks, asksOk := lastAsksMap[currencyPair]
		lastBidsAsksMutex.Unlock()
		depthResponse, _, httpResponse, err := f.requester.DepthNoRetry(currencyPair)
		if err != nil {
			if httpResponse.StatusCode == 403 {
				log.Printf("occured 403 Forbidden currency pair = %v", currencyPair)
			}
			log.Printf("can not get depth currency pair = %v", currencyPair)
			continue
		}
		if !bidsOk || !asksOk || reflect.DeepEqual(lastBids, depthResponse.Bids) == false || reflect.DeepEqual(lastAsks, depthResponse.Asks) == false {

			f.currencyPairsInfo.updateDepth(currencyPair, depthResponse.Bids, depthResponse.Asks)


			err = f.streamingCallback(currencyPair, f)

			if err != nil {
				log.Printf("streaming callback error in polling loop (%v)", err)
			}
			lastBidsAsksMutex.Lock()
			lastBidsMap[currencyPair] = depthResponse.Bids
			lastAsksMap[currencyPair] = depthResponse.Asks
			lastBidsAsksMutex.Unlock()
		}
	}
}

func  (f *fetcher) pollingRequestLoop() {
	log.Printf("start polling request loop")
	atomic.StoreInt32(&f.pollingFinish, 0)
	lastBidsMap := make(map[string][][]float64)
	lastAsksMap := make(map[string][][]float64)
	lastBidsAsksMutex := new(sync.Mutex)
	pollingRequestChan := make(chan string)
	for i := 0; i < f.config.PollingConcurrency; i++ {
		go f.pollingLoop(pollingRequestChan, lastBidsMap, lastAsksMap, lastBidsAsksMutex)
	}
FINISH:
	for {
		log.Printf("start get depth of currency Pairs (%v)", time.Now().UnixNano())
		for _, currencyPair := range f.config.CurrencyPairs {
			if atomic.LoadInt32(&f.pollingFinish) == 1{
				break FINISH
			}
			pollingRequestChan <- currencyPair
		}
	}
	close(pollingRequestChan)
	log.Printf("finish polling request loop")
}


func (f *fetcher) Start() {
	go f.pollingRequestLoop()
}

func (f *fetcher) Stop() {
	atomic.StoreInt32(&f.pollingFinish, 1)
}

func NewFetcher(config *configurator.ZaifProxyConfig) (*fetcher) {
	requesterKeys := make([]*zaif.RequesterKey, 0)
	if config.PollingConcurrency == 0 {
		config.PollingConcurrency = defaultPollingConcurrency
	}
	return &fetcher {
		requester:  zaif.NewRequester(requesterKeys, config.Retry, config.RetryWait, config.Timeout, config.ReadBufSize, config.WriteBufSize),
		config: config,
		pollingFinish: 0,
		currencyPairsInfo: &currencyPairsInfo{
			Bids:      make(map[string][][]float64),
			Asks:      make(map[string][][]float64),
			LastPrice: make(map[string]float64),
			mutex:     new(sync.Mutex),
		},
	}
}