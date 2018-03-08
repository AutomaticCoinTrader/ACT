package server

import (
	"net/http"
	"log"
	"context"
	"github.com/AutomaticCoinTrader/ACT/tools/proxy/zaif/configurator"
	"github.com/gorilla/websocket"
	"time"
	"sync"
)

type WebsocketServer struct {
	config             *configurator.ZaifProxyConfig
	server             *http.Server
	upgrader           *websocket.Upgrader
	clients            map[string]map[*websocket.Conn]bool
	clientsMutex       *sync.Mutex
	pingStopChan chan bool
	pingStopCompleteChan chan bool
}


func (w *WebsocketServer) pingLoop() {
	for {
		select {
		case _, ok := <-w.pingStopChan:
			if !ok {
				close(w.pingStopCompleteChan)
				return
			}
		case <-time.After(5 * time.Second):
			deadline := time.Now()
			deadline.Add(30 * time.Second)
			for _, cs := range w.clients {
				for ws := range cs {
					ws.WriteControl(websocket.PingMessage, []byte("ping"), deadline)
				}
			}
		}
	}
}

func (w *WebsocketServer) listenAndServe() {
	if err := w.server.ListenAndServe(); err != nil {
		log.Printf("ListenAndServe returns an error (%v)", err)
		if err != http.ErrServerClosed {
			log.Fatalf("HTTPServer closed with error (%v)", err)
		}
	}
}

func (w *WebsocketServer) BroadCast(currencyPair string, message []byte) {
	w.clientsMutex.Lock()
	defer w.clientsMutex.Unlock()
	cs, ok := w.clients[currencyPair]
	if !ok {
		return
	}
	for ws := range cs {
		ws.WriteMessage(websocket.TextMessage, message)
	}
}

func (w *WebsocketServer) Start() {
	log.Printf("start http server")
	go w.listenAndServe()
	go w.pingLoop()
}

func (w *WebsocketServer) Stop() {
	close(w.pingStopChan)
	<-w.pingStopCompleteChan
	w.server.Shutdown(context.Background())
	log.Printf("stop http server")
}

func NewWsServer(config *configurator.ZaifProxyConfig) (*WebsocketServer) {
	wsServer := &WebsocketServer{
		config:             config,
		server:             nil,
		upgrader:           new(websocket.Upgrader),
		clients:            make(map[string]map[*websocket.Conn]bool),
		clientsMutex:       new(sync.Mutex),
		pingStopChan:       make(chan bool),
		pingStopCompleteChan: make(chan bool),
	}
	for _, currencyPair := range config.CurrencyPairs {
		switch currencyPair {
		case "btc_jpy":
			http.HandleFunc("/btc_jpy", wsServer.btcJpyHandleConnection)
		case "xem_jpy":
			http.HandleFunc("/xem_jpy", wsServer.xemJpyHandleConnection)
		case "mona_jpy":
			http.HandleFunc("/mona_jpy", wsServer.monaJpyHandleConnection)
		case "bch_jpy":
			http.HandleFunc("/bch_jpy", wsServer.bchJpyHandleConnection)
		case "eth_jpy":
			http.HandleFunc("/eth_jpy", wsServer.ethJpyHandleConnection)
		case "zaif_jpy":
			http.HandleFunc("/zaif_jpy", wsServer.zaifJpyHandleConnection)
		case "pepecash_jpy":
			http.HandleFunc("/pepecash_jpy", wsServer.pepecashJpyHandleConnection)
		case "xem_btc":
			http.HandleFunc("/xem_btc", wsServer.xemBtcHandleConnection)
		case "mona_btc":
			http.HandleFunc("/mona_btc", wsServer.monaBtcHandleConnection)
		case "bch_btc":
			http.HandleFunc("/bch_btc", wsServer.bchBtcHandleConnection)
		case "eth_btc":
			http.HandleFunc("/eth_btc", wsServer.ethBtcHandleConnection)
		case "zaif_btc":
			http.HandleFunc("/zaif_btc", wsServer.zaifBtcHandleConnection)
		case "pepecash_btc":
			http.HandleFunc("/pepecash_btc", wsServer.pepecashBtcHandleConnection)
		}
	}
	wsServer.server = &http.Server{Addr: config.Server.AddrPort}
	return wsServer
}