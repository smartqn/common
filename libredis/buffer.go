package libredis

import (
	"encoding/json"
)

type ConfigRobot struct {
	PrivateKey     string `bson:"privateKey"`
	Enabled        int    `bson:"enabled"`
	AutoCoordinate int    `json:"autoCoordinate"`
	//@TODO  add
	SwitchAfterErrorTimes int64 `json:"switchAfterErrorTimes"`
}

const PrefixConfigRobot = "ConfigRobot_"

const ErrorRedis = -100

func GetConfigRobot(client *RedisOpt, enterpriseID string) (err error, config *ConfigRobot) {
	if client == nil {
		return
	}

	config = new(ConfigRobot)
	key := PrefixConfigRobot + enterpriseID
	_, bs := client.Get(key)
	err = json.Unmarshal([]byte(bs), config)
	return
}
