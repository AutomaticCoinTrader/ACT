package server

import (
	"net/http"
	"log"
	"context"
	"github.com/AutomaticCoinTrader/ACT/tools/proxy/zaif/configurator"
	"github.com/gorilla/websocket"
)

type WebsocketServer struct {
	config             *configurator.ZaifProxyConfig
	server             *http.Server
	upgrader           *websocket.Upgrader
	clients            map[string]map[*websocket.Conn]bool
}

func (s *WebsocketServer) listenAndServe() {
	if err := s.server.ListenAndServe(); err != nil {
		log.Printf("ListenAndServe returns an error (%v)", err)
		if err != http.ErrServerClosed {
			log.Fatalf("HTTPServer closed with error (%v)", err)
		}
	}
}

func (s *WebsocketServer) BroadCast(currencyPair string, message []byte) {
	cs, ok := s.clients[currencyPair]
	if !ok {
		return
	}
	for ws := range cs {
		ws.WriteMessage(websocket.TextMessage, message)
	}
}

func (s *WebsocketServer) Start() {
	log.Printf("start http server")
	go s.listenAndServe()
}

func (s *WebsocketServer) Stop() {
	s.server.Shutdown(context.Background())
	log.Printf("stop http server")
}

func NewWsServer(config *configurator.ZaifProxyConfig) (*WebsocketServer) {
	wsServer := &WebsocketServer{
		config:             config,
		server:             nil,
		upgrader:           new(websocket.Upgrader),
		clients:            make(map[string]map[*websocket.Conn]bool),
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