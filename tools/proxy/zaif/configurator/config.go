package configurator

type serverConfig struct {
	AddrPort string `json:"addrPort"  yaml:"addrPort"  toml:"addrPort"`
}

type loggerConfig struct {
	Output string `json:"output"    yaml:"output"    toml:"output"`
}

type ZaifProxyConfig struct {
	Retry              int                         `json:"retry"              yaml:"retry"              toml:"retry"`
	RetryWait          int                         `json:"retryWait"          yaml:"retryWait"          toml:"retryWait"`
	Timeout            int                         `json:"timeout"            yaml:"timeout"            toml:"timeout"`
	ReadBufSize        int                         `json:"readBufSize"        yaml:"readBufSize"        toml:"readBufSize"`
	WriteBufSize       int                         `json:"writeBufSize"       yaml:"writeBufSize"       toml:"writeBufSize"`
	PollingConcurrency int                         `json:"pollingConcurrency" yaml:"pollingConcurrency" toml:"pollingConcurrency"`
	CurrencyPairs      []string                    `json:"currencyPairs"      yaml:"currencyPairs"      toml:"currencyPairs"`
	Server             *serverConfig               `json:"server"             yaml:"server"             toml:"server"`
	Logger             *loggerConfig               `json:"logger"             yaml:"logger"             toml:"logger"`
}
