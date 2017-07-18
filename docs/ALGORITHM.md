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
type TradeConfig struct {
	message string `json:"message" yaml:"message" toml:"message"`
}

type ArbitrageTradeConfig struct {
	message string `json:"message" yaml:"message" toml:"message"`
}

type Config struct {
	Trade          *TradeConfig          `json:"trade"          yaml:"trade"          toml:"trade"`
    ArbitrageTrade *ArbitrageTradeConfig `json:"arbitrageTrade" yaml:"arbitrageTrade" toml:"arbitrageTrade"`
}


```

### 3.TradeAlgorithmContextインターフェイスを備えた構造体を作る

```
type My struct {
	config         *TradeConfig
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
	config         *ArbitrageTradeConfig
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
	configFilePathPrefix := path.Join(configDir, algorithmName)
	cf, err := configurator.NewConfigurator(configFilePathPrefix)
	if err != nil {
		return nil, errors.Errorf("can not create configurator  (config file path prefix = %v)", configFilePathPrefix)
	}
	config := new(Config)
	err = cf.Load(config)
	if err != nil {
		return nil, errors.Errorf("can not load config (config file path prefix = %v)", configFilePathPrefix)
	}
	return &My{
	    name: "my",
	    config: config.Trade,
	}, nil
}

func newArbitrageMy(config interface{}) (algorithm.ArbitrageTradeAlgorithmContext, error) {
	configFilePathPrefix := path.Join(configDir, algorithmName)
	cf, err := configurator.NewConfigurator(configFilePathPrefix)
	if err != nil {
		return nil, errors.Errorf("can not create configurator (config file path prefix = %v)", configFilePathPrefix)
	}
	config := new(Config)
	err = cf.Load(config)
	if err != nil {
		return nil, errors.Errorf("can not load config (config file path prefix = %v)", configFilePathPrefix)
	}
	return &ArbitrageMy{
	    name: "my",
	    config: config.ArbitrageTrade,
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

### 1. ACT/robot/import.goにmyパッケージを登録する

```
import (
	_ "github.com/AutomaticCoinTrader/ACT/algorithm/example"
	_ "github.com/AutomaticCoinTrader/ACT/algorithm/my"
)

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
