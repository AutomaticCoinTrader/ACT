package exchange

type OrderAction string

const (
	OrderActSell   OrderAction = "sell"
	OrderActBuy    OrderAction = "buy"
	OrderActUnkown OrderAction = "unknown"
)

type OrderCursor interface {
	Next() (orderID int64, currencyPair string, action OrderAction, price float64, amount float64, timestamp int64, ok bool)
	Reset()
	Len() int
}

type BoardCursor interface {
	Next() (price float64, amount float64, ok bool)
	Reset()
	Len() int
	All() [][]float64
}

type TradesCursor interface {
	Next() (time int64, peice float64, amount float64, tradeType string, ok bool)
	Reset()
	Len() int
}

// トレードコンテキストが更新されるたびに呼ばれる
type StreamingCallback func(currencyPair string, ex Exchange) (error)
type RetryCallback func(price *float64, amount *float64, errMsg string, retryCallbackData interface{}) (bool)

type Exchange interface {
	GetName() string
	GetCurrencyPairs() ([]string)
	Buy(currencyPair string, price float64, amount float64, retryCallback RetryCallback, retryCallbackData interface{}) (int64, float64, float64, error)
	Sell(currencyPair string, price float64, amount float64, retryCallback RetryCallback, retryCallbackData interface{}) (int64, float64, float64, error)
	Cancel(orderID int64, currencyPair string) (error)
	GetFunds() (map[string]float64, error)
	GetLastPrice(currencyPair string) (float64, error)
	GetSellBoardCursor(currencyPair string) (BoardCursor, error)
	GetBuyBoardCursor(currencyPair string) (BoardCursor, error)
	GetSellBuyBoardCursor(currencyPair string) (BoardCursor, BoardCursor, error)
	GetTradesCursor(currencyPair string) (TradesCursor, error)
	GetOrderHistoryCursor(count int64) (OrderCursor, error)
	GetActiveOrderCursor() (OrderCursor, error)
	GetMinPriceUnit(currencyPair string) (float64)
	GetMinAmountUnit(currencyPair string) (float64)
	GetTradeFeeRate(currencyPair string) (float64)
        FixPrice(currencyPair string, price float64) (float64)
        FixAmount(currencyPair string, amount float64) (float64)
	Initialize(streamingCallback StreamingCallback) (error)
	Finalize() (error)
	StartStreamings() (error)
	StopStreamings() (error)
}
