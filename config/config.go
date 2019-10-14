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
	NotifyRobotAddr     string `yaml:"notifyRobotAddr"`
	RobotFinishCallReq  string `yaml:"robotFinishCallReq"`
	EnterpriseStr       string `yaml:"enterpriseStr"`
	EnterpriseStrEx     string `yaml:"enterpriseStrEx"`
	FlowIdStr           string `yaml:"flowIdStr"`
	FlowIdStrEx         string `yaml:"flowIdStrEx"`
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
		panic("SmsSrvAddr can't be empty")
	}
	if config.MsgSendApi == "" {
		config.MsgSendApi = "/esms/sendsms"
	}
	if config.SMSIndustryAccount == "" {
		panic("SMSIndustryAccount can't be empty")
	}
	if config.SMSIndustryPassword == "" {
		panic("SMSIndustryPassword can't be empty")
	}
	if config.SMSMarketAccount == "" {
		panic("SMSMarketAccount can't be empty")
	}
	if config.SMSMarketPassword == "" {
		panic("SMSMarketPassword can't be empty")
	}

	rootConfig = config
	return
}

//var HttpClient *http.Client
var rootConfig *Configuration

func Config() *Configuration {
	return rootConfig
}
