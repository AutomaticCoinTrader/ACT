# ACT

## 概要

仮想通貨のトレード用に開発しているトレードボットです。
アルゴリズムと取引所を簡単に追加できるように作っています。

## サポート環境

 - Linux
 - MacOS (制限あり)
 - Windows (制限あり)

## goのバージョン

 - 1.8以降

## 機能

 - アルゴリズムによる自動トレード (未実装)
   - 各銘柄毎の取引、裁定取引に対応予定
 - webuiからのマニュアルによるトレード (未実装)
 - アラート通知機能 (未実装)

## アルゴリズム

 - example 
   - 動作確認用サンプル実装。実際のトレードに使うことはできません。
 - lazydog
   - Stochastic RSI と DMI を併用したアルゴリズムのトレードボッド。(予定)

## 取引所対応状況

  - [ ] bitbank (*)
  - [ ] bitFlyer (*)
  - [ ] BITPoint
  - [ ] coincheck
  
  - [ ] DMMBitcoin
  
  - [ ] GMOCoin
  
  - [ ] QUOINEX (*)
  - [x] Zaif    (*)


## ビルド方法

```
go build
```

## 設定例
 - 設定はyaml,toml,jsonいずれかで記述する
 - config/act.yaml

```
robot:
  algorithmPluginDir: "plugin"
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
server:
  debug: true
  addrPort: 127.0.0.1:38080
notifier:
  mail:
    hostPort: "smtp.gmail.com:465"
    username: "username"
    password: "password"
    authType: "plain"
    useTls: true
    useStartTls: "false"
    from: "exmaple@gmail.com"
    to: "exmaple@gmail.com"
```

## 起動

```
act -confdir ./config
```

## 停止

```
pkill act
```

## その他

  - [開発情報](/docs/DEVELOP.md)
  - [寄付](/docs/DONATION.md)
