package exchange

type OrderAction string

const (
	OrderActSell OrderAction = "sell"
	OrderActBuy  OrderAction = "buy"
	OrderActUnkown  OrderAction = "unknown"
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
	PriceAll() []float64
}

type TradesCursor interface {
	Next() (time int64, peice float64, amount float64, tradeType string, ok bool)
	Reset()
	Len() int
}

// トレードコンテキストが更新されるたびに呼ばれる
type StreamingCallback func(currencyPair string, ex Exchange) (error)

type Exchange interface {
	GetName() string
	GetCurrencyPairs() ([]string)
	Buy(currencyPair string, price float64, amount float64) (int64, error)
	Sell(currencyPair string, price float64, amount float64) (int64, error)
	Cancel(orderID int64) (error)
	GetFunds() (map[string]float64, error)
	GetLastPrice(currencyPair string) (float64, error)
	GetSellBoardCursor(currencyPair string) (BoardCursor, error)
	GetBuyBoardCursor(currencyPair string) (BoardCursor, error)
	GetTradesCursor(currencyPair string) (TradesCursor, error)
	GetOrderHistoryCursor(count int64) (OrderCursor, error)
	GetActiveOrderCursor() (OrderCursor, error)
	GetMinPriceUnit(currencyPair string) (float64)
	GetMinAmountUnit(currencyPair string) (float64)
	Initialize() (error)
	Finalize() (error)
	StartStreamings(streamingCallback StreamingCallback) (error)
	StopStreamings() (error)
}