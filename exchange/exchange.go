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

type HistoryCursor interface {
	Next() (price float64, amount float64, ok bool)
	Reset()
	Len() int
}

type TradeContextCursor interface {
	Next() (tradeContext TradeContext, ok bool)
	Reset()
	Len() int
}

type TradeContext interface {
	GetID() string
	GetExchangeName() string
	Buy(price float64, amount float64) (error)
	Sell(price float64, amount float64) (error)
	Cancel(orderID int64) (error)
	GetSrcCurrencyFund() (float64, error)
	GetDstCurrencyFund() (float64, error)
	GetSrcCurrencyName() (string)
	GetDstCurrencyName() (string)
	GetPrice() (float64, error)
	GetBuyHistoryCursor() (HistoryCursor, error)
	GetSellHistoryCursor() (HistoryCursor, error)
	GetActiveOrderCursor() (OrderCursor, error)
	GetMinPriceUnit() (float64)
	GetMinAmountUnit() (float64)
}

// トレードコンテキストが更新されるたびに呼ばれる
type StreamingCallback func(tradeContext TradeContext, userCallbackData interface{}) (error)

type Exchange interface {
	GetName() string
	Initialize(streamingCallback StreamingCallback, userCallbackData interface{}) (error)
	Finalize() (error)
	GetTradeContext(srcCurrency string, dstCurrency string) (TradeContext, bool)
	GetTradeContextCursor() (TradeContextCursor)
	StartStreaming(tradeContext TradeContext) (error)
	StopStreaming(tradeContext TradeContext) (error)
}

func MakeTradeID(exchangeName string, currencyPair string) string {
	return exchangeName + ":" + currencyPair
}
