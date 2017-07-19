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
type tradeConfig struct {
	message string `json:"message" yaml:"message" toml:"message"`
}

type arbitrageTradeConfig struct {
	message string `json:"message" yaml:"message" toml:"message"`
}

type config struct {
	Trade          *tradeConfig          `json:"trade"          yaml:"trade"          toml:"trade"`
    ArbitrageTrade *arbitrageTradeConfig `json:"arbitrageTrade" yaml:"arbitrageTrade" toml:"arbitrageTrade"`
}


```

### 3.TradeAlgorithmContextインターフェイスを備えた構造体を作る

```
type my struct {
	config         *tradeConfig
	name           string
}

func (m *my) GetName() (string) {
    // 名前を返す
	return l.name
}

func (m *my) Initialize(tradeContext exchange.TradeContext, notifier *notifier.Notifier) (error) {
    // 初期化処理
	return nil
}

func (m *my) Update(tradeContext exchange.TradeContext, notifier *notifier.Notifier) (error) {
    // トレード情報を取得するたびに呼ばれる
    // ここでトレードをやる
    
    fmt.Print(l.config.message)
    return nil
}

func (m *my) Finalize(tradeContext exchange.TradeContext, notifier *notifier.Notifier) (error) {
    // 終了処理
	return nil
}


type arbitrageMy struct {
	config         *arbitrageTradeConfig
	name           string
}

func (m *arbitrageMy) GetName() (string) {
    // 名前を返す
	return l.name
}

func (m *arbitrageMy) Initialize(exchanges map[string]exchange.Exchange, notifier *notifier.Notifier) (error) {
    // 初期化処理
	return nil
}

func (m *arbitrageMy) Update(exchanges map[string]exchange.Exchange, notifier *notifier.Notifier) (error) {
    // 定期的に呼ばれる
    // ここでトレードをやる
    
    fmt.Print(l.config.message)
    return nil
}

func (m *arbitrageMy) Finalize(exchanges map[string]exchange.Exchange, notifier *notifier.Notifier) (error) {
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
	conf := new(Config)
	err = cf.Load(conf)
	if err != nil {
		return nil, errors.Errorf("can not load config (config file path prefix = %v)", configFilePathPrefix)
	}
	return &my{
	    name: "my",
	    config: conf.Trade,
	}, nil
}

func newArbitrageMy(config interface{}) (algorithm.ArbitrageTradeAlgorithmContext, error) {
	configFilePathPrefix := path.Join(configDir, algorithmName)
	cf, err := configurator.NewConfigurator(configFilePathPrefix)
	if err != nil {
		return nil, errors.Errorf("can not create configurator (config file path prefix = %v)", configFilePathPrefix)
	}
	conf := new(Config)
	err = cf.Load(conf)
	if err != nil {
		return nil, errors.Errorf("can not load config (config file path prefix = %v)", configFilePathPrefix)
	}
	return &arbitrageMy{
	    name: "my",
	    config: conf.ArbitrageTrade,
	}, nil
}
```

### 5. (4)で作った関数を登録する

```
func init() {
	algorithm.RegisterAlgorithm("my", newMy, newArbitrageMy) // "my"という名前で登録
}
```

## インポートに追加

### 1. ACT/robot/import.goにmyパッケージを登録する

```
import (
	_ "github.com/AutomaticCoinTrader/ACT/algorithm/example"
	_ "github.com/AutomaticCoinTrader/ACT/algorithm/my"
)

```

## 設定を追加

### 1. confDirで指定したコンフィグディレクトリ以下のalgorithmディレクトリの下にmy.yamlを追加する

```
trade:
  message: "hello!!"
arbitrageTrade:
  message: "hello!!"
```

## 起動

### 1. 設定ファイルを指定して起動する

```
go run main.go -confdir ./config
```


# プラグイン方式によるアルゴリズムの追加の仕方

## 準備

### 1. 任意のディレクト以下にファイルを作成する

```
mkdir my/
touch my/my.go
```
  
## my.goの作成

### 1. パッケージ名はmainで良い

```
    package main
```


### 2. 設定情報を保持する構造体を作る
  - tomlとyamlとjsonで解釈できるように定義をしておく

```
type tradeConfig struct {
	message string `json:"message" yaml:"message" toml:"message"`
}

type arbitrageTradeConfig struct {
	message string `json:"message" yaml:"message" toml:"message"`
}

type config struct {
	Trade          *tradeConfig          `json:"trade"          yaml:"trade"          toml:"trade"`
    ArbitrageTrade *arbitrageTradeConfig `json:"arbitrageTrade" yaml:"arbitrageTrade" toml:"arbitrageTrade"`
}


```

### 3.TradeAlgorithmContextインターフェイスを備えた構造体を作る

```
type my struct {
	config         *tradeConfig
	name           string
}

func (m *my) GetName() (string) {
    // 名前を返す
	return l.name
}

func (m *my) Initialize(tradeContext exchange.TradeContext, notifier *notifier.Notifier) (error) {
    // 初期化処理
	return nil
}

func (m *my) Update(tradeContext exchange.TradeContext, notifier *notifier.Notifier) (error) {
    // トレード情報を取得するたびに呼ばれる
    // ここでトレードをやる
    
    fmt.Print(l.config.message)
    return nil
}

func (m *my) Finalize(tradeContext exchange.TradeContext, notifier *notifier.Notifier) (error) {
    // 終了処理
	return nil
}


type arbitrageMy struct {
	config         *arbitrageTradeConfig
	name           string
}

func (m *arbitrageMy) GetName() (string) {
    // 名前を返す
	return l.name
}

func (m *arbitrageMy) Initialize(exchanges map[string]exchange.Exchange, notifier *notifier.Notifier) (error) {
    // 初期化処理
	return nil
}

func (m *arbitrageMy) Update(exchanges map[string]exchange.Exchange, notifier *notifier.Notifier) (error) {
    // 定期的に呼ばれる
    // ここでトレードをやる
    
    fmt.Print(l.config.message)
    return nil
}

func (m *arbitrageMy) Finalize(exchanges map[string]exchange.Exchange, notifier *notifier.Notifier) (error) {
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
	conf := new(Config)
	err = cf.Load(conf)
	if err != nil {
		return nil, errors.Errorf("can not load config (config file path prefix = %v)", configFilePathPrefix)
	}
	return &my{
	    name: "my",
	    config: conf.Trade,
	}, nil
}

func newArbitrageMy(config interface{}) (algorithm.ArbitrageTradeAlgorithmContext, error) {
	configFilePathPrefix := path.Join(configDir, algorithmName)
	cf, err := configurator.NewConfigurator(configFilePathPrefix)
	if err != nil {
		return nil, errors.Errorf("can not create configurator (config file path prefix = %v)", configFilePathPrefix)
	}
	conf := new(Config)
	err = cf.Load(conf)
	if err != nil {
		return nil, errors.Errorf("can not load config (config file path prefix = %v)", configFilePathPrefix)
	}
	return &arbitrageMy{
	    name: "my",
	    config: conf.ArbitrageTrade,
	}, nil
}
```

### 5. pluginの初期化はinit関数で行う

```
func init() {
}
```

### 6. (4)で作った関数やアルゴリズム名を返すGetRegistrationInfo関数を作成する

```
func GetRegistrationInfo() (string, algorithm.TradeAlgorithmNewFunc, algorithm.ArbitrageTradeAlgorithmNewFunc)
    return "my", newMy, newArbitrageMy
}
```

### pluginでビルドする
  - できたものをpluginディレクトリに配置する

```
go build -buildmode=plugin
```



## 設定を追加

## 1. confDirで指定したコンフィグディレクトリ以下のact.yamlにpluginディレクトリのパスを追加する

```
robot:
  algorithmPluginDir: "plugin"
integrator:
  debug: true
  addrPort: 127.0.0.1:38080
  exchanges:
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
notifier:
  mail:
```

## 2. confDirで指定したコンフィグディレクトリ以下のalgorithmディレクトリの下にmy.yamlを追加する

```
trade:
  message: "hello!!"
arbitrageTrade:
  message: "hello!!"
```

## 起動

### 1. 設定ファイルを指定して起動する

```
go run main.go -confdir ./config
```

