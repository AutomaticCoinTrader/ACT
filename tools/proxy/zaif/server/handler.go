package server

import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
)

func (w *WebsocketServer) handleConnectionBase(writer http.ResponseWriter, request *http.Request, currencyPair string) {
	// Upgrade initial GET request to a websocket
	ws, err := w.upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()
	log.Printf("connected peer address = %v", ws.RemoteAddr().String())
	cs, ok := w.clients[currencyPair]
	if !ok {
		cs = make(map[*websocket.Conn]bool)
		w.clients[currencyPair] = cs
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
		}
		// nop
	}
}

func (w *WebsocketServer) btcJpyHandleConnection(writer http.ResponseWriter, request *http.Request) {
	w.handleConnectionBase(writer, request, "btc_jpy")
}

func (w *WebsocketServer) xemJpyHandleConnection(writer http.ResponseWriter, request *http.Request) {
	w.handleConnectionBase(writer, request,  "xem_jpy")
}

func (w *WebsocketServer) monaJpyHandleConnection(writer http.ResponseWriter, request *http.Request) {
	w.handleConnectionBase(writer, request,  "mona_jpy")
}

func (w *WebsocketServer) bchJpyHandleConnection(writer http.ResponseWriter, request *http.Request) {
	w.handleConnectionBase(writer, request,  "bch_jpy")
}

func (w *WebsocketServer) ethJpyHandleConnection(writer http.ResponseWriter, request *http.Request) {
	w.handleConnectionBase(writer, request,  "eth_jpy")
}

func (w *WebsocketServer) zaifJpyHandleConnection(writer http.ResponseWriter, request *http.Request) {
	w.handleConnectionBase(writer, request,  "zaif_jpy")
}

func (w *WebsocketServer) pepecashJpyHandleConnection(writer http.ResponseWriter, request *http.Request) {
	w.handleConnectionBase(writer, request,  "pepecash_jpy")
}

func (w *WebsocketServer) xemBtcHandleConnection(writer http.ResponseWriter, request *http.Request) {
	w.handleConnectionBase(writer, request,  "xem_btc")
}

func (w *WebsocketServer) monaBtcHandleConnection(writer http.ResponseWriter, request *http.Request) {
	w.handleConnectionBase(writer, request, "mona_btc")
}

func (w *WebsocketServer) bchBtcHandleConnection(writer http.ResponseWriter, request *http.Request) {
	w.handleConnectionBase(writer, request,  "bch_btc")
}

func (w *WebsocketServer) ethBtcHandleConnection(writer http.ResponseWriter, request *http.Request) {
	w.handleConnectionBase(writer, request,  "eth_btc")
}

func (w *WebsocketServer) zaifBtcHandleConnection(writer http.ResponseWriter, request *http.Request) {
	w.handleConnectionBase(writer, request, "zaif_btc")
}

func (w *WebsocketServer) pepecashBtcHandleConnection(writer http.ResponseWriter, request *http.Request) {
	w.handleConnectionBase(writer, request,  "pepecash_btc")
}

