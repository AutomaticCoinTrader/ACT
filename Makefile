all: act zaif-proxy

act:
	go build
zaif-proxy:
	cd tools/proxy/zaif && go build -o zaif-proxy
