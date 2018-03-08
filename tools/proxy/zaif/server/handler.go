package server

import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
)

func (s *WebsocketServer) handleConnectionBase(writer http.ResponseWriter, request *http.Request, currencyPair string) {
	// Upgrade initial GET request to a websocket
	ws, err := s.upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()
	cs, ok := s.clients[currencyPair]
	if !ok {
		cs = make(map[*websocket.Conn]bool)
		s.clients[currencyPair] = cs
	}
	cs[ws] = true
	for {
		// Read in a new message as JSON and map it to a Message object
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("can not read message (error = %v)", err)
			delete(cs, ws)
			break
		}
		if messageType != websocket.TextMessage {
			log.Printf("unsupported message type (message type = %v, message = %v)", messageType, message)
			delete(cs, ws)
			break
		}
		// nop
	}
}

func (s *WebsocketServer) btcJpyHandleConnection(writer http.ResponseWriter, request *http.Request) {
	s.handleConnectionBase(writer, request, "btc_jpy")
}

func (s *WebsocketServer) xemJpyHandleConnection(writer http.ResponseWriter, request *http.Request) {
	s.handleConnectionBase(writer, request,  "xem_jpy")
}

func (s *WebsocketServer) monaJpyHandleConnection(writer http.ResponseWriter, request *http.Request) {
	s.handleConnectionBase(writer, request,  "mona_jpy")
}

func (s *WebsocketServer) bchJpyHandleConnection(writer http.ResponseWriter, request *http.Request) {
	s.handleConnectionBase(writer, request,  "bch_jpy")
}

func (s *WebsocketServer) ethJpyHandleConnection(writer http.ResponseWriter, request *http.Request) {
	s.handleConnectionBase(writer, request,  "eth_jpy")
}

func (s *WebsocketServer) zaifJpyHandleConnection(writer http.ResponseWriter, request *http.Request) {
	s.handleConnectionBase(writer, request,  "zaif_jpy")
}

func (s *WebsocketServer) pepecashJpyHandleConnection(writer http.ResponseWriter, request *http.Request) {
	s.handleConnectionBase(writer, request,  "pepecash_jpy")
}

func (s *WebsocketServer) xemBtcHandleConnection(writer http.ResponseWriter, request *http.Request) {
	s.handleConnectionBase(writer, request,  "xem_btc")
}

func (s *WebsocketServer) monaBtcHandleConnection(writer http.ResponseWriter, request *http.Request) {
	s.handleConnectionBase(writer, request, "mona_btc")
}

func (s *WebsocketServer) bchBtcHandleConnection(writer http.ResponseWriter, request *http.Request) {
	s.handleConnectionBase(writer, request,  "bch_btc")
}

func (s *WebsocketServer) ethBtcHandleConnection(writer http.ResponseWriter, request *http.Request) {
	s.handleConnectionBase(writer, request,  "eth_btc")
}

func (s *WebsocketServer) zaifBtcHandleConnection(writer http.ResponseWriter, request *http.Request) {
	s.handleConnectionBase(writer, request, "zaif_btc")
}

func (s *WebsocketServer) pepecashBtcHandleConnection(writer http.ResponseWriter, request *http.Request) {
	s.handleConnectionBase(writer, request,  "pepecash_btc")
}

