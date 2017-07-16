package notifier

import "github.com/AutomaticCoinTrader/ACT/utility"

type Notifier struct {
	mailClient *utility.MailClient
}

func (n *Notifier) SendMail(subject string, body string) (error) {
	if n.mailClient == nil {
		return nil
	}
	return n.mailClient.SendMail(subject, body)
}

type MailConfig struct {
	HostPort    string `json:"hostPort"    yaml:"hostPort"    toml:"hostPort"`
	Username    string `json:"username"    yaml:"username"    toml:"username"`
	Password    string `json:"password"    yaml:"password"    toml:"password"`
	AuthType    string `json:"authType"    yaml:"authType"    toml:"authType"`
	UseTLS      bool   `json:"useTls"      yaml:"useTls"      toml:"useTls"`
	UseStartTLS bool   `json:"useStartTls" yaml:"useStartTls" toml:"useStartTls"`
	From        string `json:"from"        yaml:"from"        toml:"from"`
	To          string `json:"to"          yaml:"to"          toml:"to"`
}

type Config struct {
	Mail *MailConfig `json:"mail" yaml:"mail" toml:"mail"`
}

func NewNotifier(config *Config) (*Notifier, error) {
	var mailClient *utility.MailClient
	if config != nil {
		mailClient = utility.NewMailClient(config.Mail.HostPort, config.Mail.Username,
			config.Mail.Password, utility.GetSMTPAuthType(config.Mail.AuthType),
			config.Mail.UseTLS, config.Mail.UseStartTLS, config.Mail.From, config.Mail.To)
	}
	return &Notifier{
		mailClient: mailClient,
	}, nil
}
