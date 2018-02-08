package exchange

type OrderAction string

const (
	OrderActSell OrderAction = "sell"
	OrderActBuy  OrderAction = "buy"
)

type OrderCursor interface {
	Next() (orderID int64, action OrderAction, price float64, amount float64, ok bool)
	Reset()
	Len() int
}

type BoardCursor interface {
	Next() (price float64, amount float64, ok bool)
	Reset()
	Len() int
}

type TradesCursor interface {
	Next() (time int64, peice float64, amount float64, tradeType string, ok bool)
	Reset()
	Len() int
}

type TradeContextCursor interface {
	Next() (tradeContext TradeContext, ok bool)
	Reset()
	Len() int
}

type TradeContext interface {
	GetExchangeName() string
	Buy(currencyPair string, price float64, amount float64) (int64, error)
	Sell(currencyPair string, price float64, amount float64) (int64, error)
	Cancel(orderID int64) (error)
	GetFunds() (map[string]float64, error)
	GetLastPrice(currencyPair string) (float64, error)
	GetBuyBoardCursor(currencyPair string) (BoardCursor, error)
	GetSellBoardCursor(currencyPair string) (BoardCursor, error)
	GetTradesCursor(currencyPair string) (TradesCursor, error)
	GetOrderHistoryCursor() (OrderCursor, error)
	GetActiveOrderCursor() (OrderCursor, error)
	GetMinPriceUnit(currencyPair string) (float64)
	GetMinAmountUnit(currencyPair string) (float64)
}

// トレードコンテキストが更新されるたびに呼ばれる
type StreamingCallback func(currencyPair string, tradeContext TradeContext) (error)

type Exchange interface {
	GetName() string
	Initialize(streamingCallback StreamingCallback) (error)
	Finalize() (error)
	GetTradeContext() (TradeContext)
	StartStreamings(tradeContext TradeContext) (error)
	StopStreamings(tradeContext TradeContext) (error)
}