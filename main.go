package main

import (
	"fmt"

	redispair "github.com/smartqn/common/libredis"

	"github.com/linnv/logx"
)

func main() {
	fmt.Println()
	logx.Debugf("xx: %+v\n", 3)

	redisMain := "192.168.1.125:6379"
	redisBack := "192.168.1.125:6380"
	RedisPair := redispair.InitWithPasswd(redisMain, redisBack, 2, "./", "")
	rets, err := RedisPair.Keys("test*")
	logx.Debugf("rets: %+v\n", rets)
	logx.Debugf("err: %+v\n", err)
}
