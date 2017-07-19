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

## アルゴリズム

 - example 
   - 動作確認用サンプル実装。実際のトレードに使うことはできません。

## 取引所対応状況

  - [x] Zaif
  - [ ] bitFlyer
  - [ ] bitbank
  - [ ] Btcbox
  - [ ] coincheck
  - [ ] kraken
  - [ ] QUOINEX
  - [ ] Lemuria
  - [ ] BITPoint
  - [ ] Money365
  - [ ] みんなのビットコイン
  - [ ] Fisco
  - [ ] FIREX
  - [ ] Z.comコインbyGMO

## ビルド方法

```
go build
```

## 設定例
 - config/act.yaml

```
robot:
  algorithmPluginDir: "plugin"
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
