# 開発環境

## MAC OS
1. home brewもしくはmac portsをインストールする

2. jdkをインストール
  - http://www.oracle.com/technetwork/java/javase/downloads/jdk8-downloads-2133151.html

3. goをインストールする
  - go 1.8以降推奨

```
brew install go 
```

または

```
sudo port install go
```

4. goglandをインストールする

  - https://www.jetbrains.com/go/

5. GOPATHディレクトリを作っておく

```
mkdir -p act/src/github.com/AutomaticCoinTrader
```

6. githubからcloneする

```
git clone https://github.com/AutomaticCoinTrader/ACT.git act/src/github.com/AutomaticCoinTrader/ACT
```

7. GOPATH環境変数を指定しつつ、コマンドラインからgoglandを起動

```
GOPATH=$(pwd)/act open /Applications/Gogland\ 1.0\ EAP.app
```
8. 起動したら、open projectで　act/src/github.com/AutomaticCoinTrader/ACT のパスを開く

## windows 10

1. jdkをインストール

  - http://www.oracle.com/technetwork/java/javase/downloads/jdk8-downloads-2133151.html

2. go をインストール

  - go 1.8以降推奨
  - https://storage.googleapis.com/golang/go1.8.3.windows-amd64.msi

3. goglandをインストール

  - https://www.jetbrains.com/go/

4. gitをインストール

  - https://git-for-windows.github.io/

5. windows コマンドラインを開く

6.  GOPATHディレクトリを作っておく

```
mkdir -p act\src\github.com\AutomaticCoinTrader
```

7. githubからcloneする

```
git clone https://github.com/AutomaticCoinTrader/ACT.git act\src\github.com\AutomaticCoinTrader\ACT
```

8. gogland起動

```
(set GOPATH=%CD%\act) && "C:\Program Files\JetBrains\Gogland 171.4694.61\bin\gogland64.exe"
```
9. 起動したら、open projectで　act/src/github.com/AutomaticCoinTrader/ACT のパスを開く

## Ubuntu 16.04

1. jdkをインストール

```
sudo add-apt-repository -y ppa:webupd8team/java
sudo apt update
sudo apt install -y oracle-java8-installer oracle-java8-set-default
```

2. goをインストール

```
sudo add-apt-repository -y  ppa:longsleep/golang-backports
sudo apt-get update
sudo apt-get install -y golang-go
```

3. Goglandをインストール

```
wget https://download.jetbrains.com/go/gogland-171.4694.61.tar.gz
tar -zxvf gogland-171.4694.61.tar.gz
sudo mkdir /opt/Gogland
sudo mv Gogland-171.4694.61 /opt/Gogland/
sudo ln -s /opt/Gogland/Gogland-171.4694.61 /opt/Gogland/latest
```

4. GOPATHディレクトリを作っておく

```
mkdir -p act/src/github.com/AutomaticCoinTrader
```

5. githubからcloneする

```
git clone https://github.com/AutomaticCoinTrader/ACT.git act/src/github.com/AutomaticCoinTrader/ACT
```

6. GOPATH環境変数を指定しつつ、コマンドラインからgoglandを起動

```
GOPATH=$(pwd)/act /opt/Gogland/latest/bin/gogland.sh
```

7. 起動したら、open projectで　act/src/github.com/AutomaticCoinTrader/ACT のパスを開く