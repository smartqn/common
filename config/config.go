package config

import (
	"sync"

	"github.com/linnv/logx"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	SMSIndustryAccount  string `yaml:"smsIndustryAccount"`
	SMSIndustryPassword string `yaml:"smsIndustryPassword"`
	SMSMarketAccount    string `yaml:"smsMarketAccount"`
	SMSMarketPassword   string `yaml:"smsMarketPassword"`
	SmsSrvAddr          string `yaml:"smsSrvAddr"`
	MsgSendApi          string `yaml:"msgSendApi"`
}

var once sync.Once

func Init(bs []byte) {
	once.Do(func() {
		_, err := initConfig(bs)
		if err != nil {
			panic(err.Error())
		}
	})
}

func initConfig(bs []byte) (config *Configuration, err error) {
	config = new(Configuration)
	err = yaml.Unmarshal(bs, &config)
	if err != nil {
		logx.Warnf("err: %+v\n", err)
		return nil, err
	}

	if config.SmsSrvAddr == "" {
		config.SmsSrvAddr = "http://39.106.59.58:9999"
	}
	if config.MsgSendApi == "" {
		config.MsgSendApi = "/esms/sendsms"
	}
	if config.SMSIndustryAccount == "" {
		config.SMSIndustryAccount = "a00038"
	}
	if config.SMSIndustryPassword == "" {
		config.SMSIndustryPassword = "qn123456"
	}
	if config.SMSMarketAccount == "" {
		config.SMSMarketAccount = "a00038"
	}
	if config.SMSMarketPassword == "" {
		config.SMSMarketPassword = "qn123456"
	}

	rootConfig = config
	return
}

//var HttpClient *http.Client
var rootConfig *Configuration

func Config() *Configuration {
	return rootConfig
}
