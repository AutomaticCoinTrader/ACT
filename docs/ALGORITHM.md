# ビルトイン方式によるアルゴリズムの追加の仕方

## 準備

### 1. ACT/algorithmの下にディレクトを作ってファイルを作成する

```
mkdir ACT/algorithm/my
touch ACT/algorithm/my/my.go
```
  
## my.goの作成

### 1. パッケージ名を決める

```
    package my
```


### 2. 設定情報を保持する構造体を作る
  - tomlとyamlとjsonで解釈できるように定義をしておく

```
type Config struct {
	message string `json:"message" yaml:"message" toml:"message"`
}

type ArbitrageConfig struct {
	message string `json:"message" yaml:"message" toml:"message"`
}

```

### 3.TradeAlgorithmContextインターフェイスを備えた構造体を作る

```
type My struct {
	config         *Config
	name           string
}

func (m *My) GetName() (string) {
    // 名前を返す
	return l.name
}

func (m *My) Initialize(tradeContext exchange.TradeContext, notifier *notifier.Notifier) (error) {
    // 初期化処理
	return nil
}

func (m *My) Update(tradeContext exchange.TradeContext, notifier *notifier.Notifier) (error) {
    // トレード情報を取得するたびに呼ばれる
    // ここでトレードをやる
    
    fmt.Print(l.config.message)
    return nil
}

func (m *My) Finalize(tradeContext exchange.TradeContext, notifier *notifier.Notifier) (error) {
    // 終了処理
	return nil
}


type ArbitrageMy struct {
	config         *Config
	name           string
}

func (m *ArbitrageMy) GetName() (string) {
    // 名前を返す
	return l.name
}

func (m *ArbitrageMy) Initialize(exchanges map[string]exchange.Exchange, notifier *notifier.Notifier) (error) {
    // 初期化処理
	return nil
}

func (m *ArbitrageMy) Update(exchanges map[string]exchange.Exchange, notifier *notifier.Notifier) (error) {
    // 定期的に呼ばれる
    // ここでトレードをやる
    
    fmt.Print(l.config.message)
    return nil
}

func (m *ArbitrageMy) Finalize(exchanges map[string]exchange.Exchange, notifier *notifier.Notifier) (error) {
    // 終了処理
	return nil
}


```

### 4. (3)で作ったTradeAlgorithmContextインターフェイスを備えた構造体のポインタを返す関数を作る

```
func newMy(config interface{}) (algorithm.TradeAlgorithmContext, error) {
    myConfig := config.(*Config) // interface{} でくるので自分の構造体にキャストする
	return &My{
	    name: "my",
	    config: myConfig,
	}, nil
}

func newArbitrageMy(config interface{}) (algorithm.ArbitrageTradeAlgorithmContext, error) {
    myConfig := config.(*ArbitrageConfig) // interface{} でくるので自分の構造体にキャストする
	return &ArbitrageMy{
	    name: "my",
	    config: myConfig,
	}, nil
}
```

### 5. (4)で作った関数を登録する

```
func init() {
	algorithm.RegisterAlgorithm("my", newMy, newArbitrageMy) // "my"という名前で登録
}
```

## 設定にに組み込む

### 1. ACT/robot/config.goにmy.goを作成の(2)で作った設定情報の構造体を追加

```
type TradeConfig struct {
	Example *example.Config `json:"example" yaml:"example" toml:"example" config:"example"` 
    My      *my.Config      `json:"my"      yaml:"my"      toml:"my"      config:"my"` // myという名前で登録したのでconfigをmyにする
}

type ArbitrageTradeConfig struct {
	Example *example.ArbitrageConfig `json:"example" yaml:"example" toml:"example" config:"example"`
    My      *my.ArbitrageConfig      `json:"my"      yaml:"my"      toml:"my"      config:"my"` // myという名前で登録したのでconfigをmyにする
}

```

## 2. config.yamlの設定を追加する

```
robot:
  trade:
    my:
      message: "hello!!"
  arbitrageTrade:
    my:
      message: "hello!!"
integrator:
  zaif:
    key: "key"
    secret: "secret"
    currencyPairs:
     - src: btc
       dst: jpy
    retry: 0
    timeout: 0
    readBufSize: 0
    writeBufSize: 0
notifiler:
   Mail:
```

## 起動

### 1. 設定ファイルを指定して起動する

```
go run main.go -config ./config/config.yaml
```


# プラグイン方式によるアルゴリズムの追加の仕方

TODO
