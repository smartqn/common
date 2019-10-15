package libredis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"testing"
	"time"
	"unsafe"
)

const (
	srv       = "127.0.0.1:6379"
	connCount = 10
	logPath   = "./"
)

func TestRedisOpt_Expire(t *testing.T) {
	rdOpt := InitWithPasswd(srv, connCount, logPath, "", 0)
	ret := rdOpt.Expire("yangjc2", 10)
	if ret == OCCUR_ERR {
		t.Fatalf("rdOpt.Expire err\n")
		return
	}

	return
}

func TestRedisOpt_Publish(t *testing.T) {
	rdOpt := InitWithPasswd(srv, connCount, logPath, "", 0)
	ret := rdOpt.Publish("yangjcPub", "nice")
	if ret == OCCUR_ERR {
		t.Fatalf("rdOpt.Expire err\n")
		return
	}

	return
}

func TestRedisOpt_Subscribe(t *testing.T) {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		t.Fatalf("redis dial failed.")
		return
	}

	var client redis.PubSubConn
	client = redis.PubSubConn{conn}

	err = client.Subscribe("yangPublish")
	if err != nil {
		t.Fatalf("redis Subscribe error.\n")
		return
	}

	c := redis.PubSubConn{conn}

	go func() {
		for {
			log.Println("wait...")
			switch res := c.Receive().(type) {
			case redis.Message:
				channel := (*string)(unsafe.Pointer(&res.Channel))
				message := (*string)(unsafe.Pointer(&res.Data))
				log.Printf("channel: %s, message: %s\n", *channel, *message)
			case redis.Subscription:
				fmt.Printf("%s: %s %d\n", res.Channel, res.Kind, res.Count)
			case error:
				log.Println("error handle...")
				panic("err: " + res.Error())
			}
		}
	}()

	time.Sleep(10 * time.Second)
}
