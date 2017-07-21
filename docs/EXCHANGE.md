# exchangeの追加の仕方

## 準備

### 1. exchangeの中に新しい取引所のディレクトリを掘る

```
mkdir ACT/exchange/myexchange
```

### 2. exchange.goとrequester.goを作る
  - exchange.goとrequester.go以外のファイルは自由。なくても良い。

```
touch ACT/exchange/myexchange/requester.go
touch ACT/exchange/myexchange/exchange.go
```

## 実装

### 1. ディレクトリ名と同じ名前のパッケージ名にする

```
package myexchange
```

### 2. requester.goを作る

```
vi ACT/exchange/myexchange/requester.go
```

  - requester.goは実装した全てのAPIを呼すことができるもの
    - newをしてオブジェクトを生成してオブジェクトの各メソッドを呼び出すことでapiを利用できるようにする
    - streaming系の処理がある場合は、callback等を用いる
    - 実装イメージは以下を参照
  
```
    r, err := myexchange.newRequester(...) // requesterを作る。このオブジェクトで全てのリクエストをこなせる。
    if err != nil {
       ...
    }
    p := myexchange.getTradeParam() // リクエストパラメータ構造体を取得する (関数化しなくても良い)
    p.action = myexchange.ActionBuy // パラメータのセット
    res, err := r.Trade(p) // トレードリクエストを実行
    if err != nil {
       ...
    }
    fmt.println(string(json.Marshal(res))) // レスポンスの出力
```

### 2. exchange.goを作る

```
vi ACT/exchange/myexchange/exchange.go
```
 - exchange.goはrequester.goをラップしてexchangeインターフェイスを実装した取引所ごとの構造体を実装する
   - exchange/excahnge.goのインターフェイスを全て実装すること 
     - APIが提供されてなくて実装できない場合は要相談
    
### 3. コンフィグ構造体を定義する

```
type ExchangeConfigCurrencyPair struct {
	Src string `json:"src" yaml:"src" toml:"src"`
	Dst string `json:"dst" yaml:"dst" toml:"dst"`
}

type ExchangeConfig struct {
	Key           string                         `json:"key"          yaml:"key"          toml:"key"`
	Secret        string                         `json:"secret"       yaml:"secret"       toml:"secret"`
	Retry         int                            `json:"retry"        yaml:"retry"        toml:"retry"`
	Timeout       int                            `json:"timeout"      yaml:"timeout"      toml:"timeout"`
	CurrencyPairs []*ExchangeConfigCurrencyPair  `json:"currencyPairs" yaml:"currencyPairs" toml:"currencyPairs"`
}
```

### 4. (2)で作ったexchangeインターフェイスを持つ構造体のオブジェクトを返す関数を作る

```
func newMyExchange(config interface{}) (exchange.Exchange, error)  {
	myConfig := config.(*ExchangeConfig)
	return &MyExchange{
		config:        myConfig,
		name :         exchangeName,
		requester:     NewRequester(myConfig.Key, myConfig.Secret, myConfig.Retry, myConfig.Timeout),
		tradeContexts: make([]*TradeContext, 0),
		funds : &ExchageFunds{
			funds: make(map[string]float64),
			mutex: new(sync.Mutex),
		},
	}, nil
}
```

### 5. (4)で作った関数を登録する

```
func init() {
	exchange.RegisterExchange("my", newMyExchange)
}
```

## コンフィグに追加

### 1. integratorのconfig.goに1の（3）で作ったコンフィグ構造体を追加する

```
type exchangesConfig struct {
	Zaif *zaif.ExchangeConfig `json:"zaif" yaml:"zaif" toml:"zaif" config:"zaif"`
	My *my.ExchangeConfig     `json:"may"  yaml:"may"  toml:"my"   config:"my"` // myでRegisterExchangeしたのでconfigはmyにする
}
```

## その他 

 - utilityのhttpclientを利用すると良い。