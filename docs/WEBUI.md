# WEBUIの作り方

## 注意

まだ、ちゃんと決めきれていない。逆に言えば自由。以下決めたことだけ。

## フレームワーク

  - ginを使って実装する 
    - https://github.com/gin-gonic/gin

## ルーティング

### intergrator.goのsetupRoutingにルーティングを追加

```
func (i *Integrator) setupRouting(engine *gin.Engine) {
	engine.HEAD( "/", i.index)
	engine.GET( "/", i.index)
}
```

## ハンドルの追加

### integrator.goのhandle.goにハンドル関数を追加していく

```
func (i *Integrator) index(context *gin.Context) {

}
```

## アセットの組み込み
  
  - template,CSS,javascript,imageのようなアセット関連は https://github.com/jteeuwen/go-bindata を使ってコンパイル時にバイナリに組み込む 

```
    cd integrator && go-bindata -pkg integrator asset/...
	go build
```






