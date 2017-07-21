# インターフェイス


## TradeAlgorithmContextインターフェイス

- 各取引所、各品目(例 btc_jpy)毎のトレードにおいて各アルゴリズムが実装するインターフェイス
  
```
type TradeAlgorithmContext interface {
	GetName() (string)
	Initialize(tradeContext exchange.TradeContext, notifier *notifier.Notifier) (error)
	Update(tradeContext exchange.TradeContext, notifier *notifier.Notifier) (error)
	Finalize(tradeContext exchange.TradeContext, notifier *notifier.Notifier) (error)
}
```

- GetName
  - アルゴリズム名を返す
- Initialize
  - トレード情報を継続的に取得する処理(streaming)が開始される前に呼ばれる
  - ここでアルゴリズムの初期化処理を行う
- Update
  - トレード情報を取得するたびに呼ばれる
  - 実際のトレード行う
- Finalize
  - トレード情報を継続的な取得(streaming)を停止した場合に呼ばれる
  - ここでアルゴリズムの終了処理を行う

(今後必要に応じてメソッドを追加していく予定)  

## ArbitrageTradeAlgorithmContextインターフェイス
 
- 取引所を跨いだトレードにおいて各アルゴリズムが実装が実装するインターフェイス

```
type ArbitrageTradeAlgorithm interface {
	GetName() (string)
	Initialize(exchanges map[string]exchange.Exchange, notifier *notifier.Notifier) (error)
	Update(exchanges map[string]exchange.Exchange, notifier *notifier.Notifier) (error)
	Finalize(exchanges map[string]exchange.Exchange, notifier *notifier.Notifier) (error)
}
```

- GetName
  - アルゴリズム名を返す
- Initialize
  - 定期的なポーリングが開始される前に呼ばれる
  - ここでアルゴリズムの初期化処理を行う
- Update
  - 定期的なポーリングのたびに呼ばれる
  - 実際のトレード行う
- Finalize
  - 定期的なポーリングを停止した場合に呼ばれる
  - ここでアルゴリズムの終了処理を行う

(今後必要に応じてメソッドを追加する予定)  
    
## Exchangeインターフェイス
  
- 各取引所が実装するインターフェイス
  
```
type Exchange interface {
	GetName() string
	Initialize(streamingCallback StreamingCallback, userCallbackData interface{}) (error)
	Finalize() (error)
	GetTradeContext(srcCurrency string, dstCurrency string) (TradeContext, bool)
	GetTradeContextCursor() (TradeContextCursor)
	StartStreaming(tradeContext TradeContext) (error)
	StopStreaming(tradeContext TradeContext) (error)
}
```

- GetName
  - 取引所名を返す
- Initialize
  - 初期化処理を行う
- Finalize
  - 終了処理を行う
- GetTradeContext
  - トレードコンテキストを取得する
  - srcとdstには通貨の名前を入れる(例. btc/jpy -> srcCurrency:btc, dstCurrency:jpy )
- GetTradeContextCursor
  - トレードコンテキストをイテレーションしながら取得するオブジェクト
  - TradeContextCursorはNextメソッドを呼び出すことでトレードコンテキストが取得できる
- StartStreaming
  - トレード情報の継続的な取得を開始する
- StartStreaming
  - トレード情報の継続的な取得を停止する

(今後必要に応じてメソッドを追加する予定)  

## TradeContextインターフェイス
  
- 各取引所の各品目(例. btc_jpt)毎に作成されるトレードコンテキストが実装するインターフェイス

```
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
	GetBuyBoardCursor() (BoardCursor, error)
	GetSellBoardCursor() (BoardCursor, error)
	GetTradeHistoryCursor() (TradeHistoryCursor, error)
	GetActiveOrderCursor() (OrderCursor, error)
	GetMinPriceUnit() (float64)
	GetMinAmountUnit() (float64)
}
```

- GetID
  - トレードコンテキストを識別するIDを返す
- GetExhangeName
  - 取引所の名前を返す
  - ExchangeインターフェイスのGetNameと同等
- Buy
  - srcCurrencyで指定されているものを買う
- Sell
  - srcCurrencyで指定されているものを売る
- Cancel
  - アクティブなオーダーをキャンセルする
  - アクティブなオーダー情報はGetActiveOrderCursorで取得可能
- GetSrcCurrencyFund
  - srcCurrencyで指定されているものの資産情報を返す
- GetDstCurrencyFund
  - dstCurrencyで指定されているものの資産情報を返す
- GetSrcCurrencyName
  - srcCurrencyの名前を返す (例. btc)
- GetDstCurrencyName
  - dstCurrencyの名前を返す (例. jpy)
- GetPrice
  - srcCurrencyで指定されているもののdstcurrency換算の現在の価格を返す
- GetBuyBoardCursor
  - 買いの板情報を返すイテレーションオブジェクト
  - 取引所によってどの程度の量を返すか違うため、データ量は取引所依存になる
- GetSellBoardCursor
  - 売りの板情報を返すイテレーションオブジェクト
  - 取引所によってどの程度の量を返すか違うため、データ量は取引所依存になる
- GetTradeHistoryCursor
  - 約定履歴を返すイテレーションオブジェクト
  - 取引所によってどの程度の量を返すか違うため、データ量は取引所依存になる
- GetMinPriceUnit
  - 最小の価格の単位を返す (例. 5円)
- GetMinAmountUnit
  - 最小の量の単位を返す (例. 0.01)
  
(今後必要に応じてメソッドを追加する予定)  

  
## カーソルインターフェイス
  - イテレーションを目的としたインターフェイス

```
type OrderCursor interface {
	Next() (orderID int64, action OrderAction, price float64, amount float64, ok bool)
	Reset()
	Len() int
}
```

```
type BoardCursor interface {
	Next() (price float64, amount float64, ok bool)
	Reset()
	Len() int
}
```

```
type TradeHistoryCursor interface {
	Next() (time int64, peice float64, amount float64, tradeType string, ok bool)
	Reset()
	Len() int
}
```

```
type TradeContextCursor interface {
	Next() (tradeContext TradeContext, ok bool)
	Reset()
	Len() int
}
```

- Next
  - 現在の値を取得して、カーソルを次に進める
  - okがfalseの場合は終わりに達して値がなかったことを表す
- Reset
  - カーソル位置を先頭に戻す
- Len
  - データ数を返す
