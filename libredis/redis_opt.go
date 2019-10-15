package libredis

/*************
1:20191015 修改，单机redis操作
2:封装框架redigo源代码
3:函数第一个返回值都代表操作情况(0:主备都操作成功，-100:主备都失败)
4:key不存在 string类型返回空字符串
5：压力测试结论
*************/

import (
	"strconv"
	"strings"
	"time"

	"github.com/aiwuTech/fileLogger"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/linnv/logx"
	//redigo "github.com/smartqn/redisgofix"
)

type RedisOpt struct {
	RedisPoll *redigo.Pool
	RedisSvr  string // default "127.0.0.1:6379"
}

var RPair_Log *fileLogger.FileLogger

func InitWithPasswd(mainSvr string, connCount int, logPath, password string, userDB int) *RedisOpt {
	RPair_Log = fileLogger.NewSizeLogger(logPath, "redisOpt.Log", "", 100, 100, fileLogger.MB, 300, 10000)
	RPair_Log.SetLogLevel(fileLogger.INFO) //trace Log will not be print

	pair := &RedisOpt{
		RedisPoll: nil,
		RedisSvr:  mainSvr,
		//RedisPoll1:          nil,
		//RedisSvr1:           backSvr,
		//rPair_LRANGE_RANDOM: 1,
	}

	pair.RedisPoll = initRedis(pair.RedisSvr, password, connCount, userDB)
	//pair.RedisPoll1 = initRedis(pair.RedisSvr1, password, connCount)

	//go autoCheckPoll("poll0", pair.RedisPoll0, 60, connCount)
	//go autoCheckPoll("poll1", pair.RedisPoll1, 60, connCount)

	//go autoReIRpushMsg(pair, 1)

	//time.Sleep(2 * time.Second)

	return pair
}

func initRedis(host string, password string, connCount int, userDB int) *redigo.Pool {
	if connCount < 3 {
		connCount = 3
	}

	RPair_Log.I("initRedis host=%s,password=%s,connCount=%d", host, password, connCount)
	return &redigo.Pool{
		MaxIdle:     connCount,
		MaxActive:   1000,
		IdleTimeout: 30 * time.Minute,
		Wait:        false,
		/*
		   TestOnBorrow: func(c redigo.Conn, t time.Time) error {
		       _, err := c.Do("PING")
		       return err
		   },
		*/
		//TestOnBorrow: nil,
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			if nil != err {
				RPair_Log.E("redis ping error:"+err.Error(), "error")
			}
			return err
		},
		Dial: func() (redigo.Conn, error) {
			RPair_Log.I("connect redis begin host=%s,password=%s,connCount=%d", host, password, connCount)
			c, err := redigo.Dial("tcp", host, redigo.DialPassword(password), redigo.DialDatabase(userDB))
			//if nil != err {
			//	return nil, err
			//}
			RPair_Log.I("REPORTALARM=CONNECT_REDIS_OK,connect redis ok host=%s,password=%s,connCount=%d", host, password, connCount)
			return c, err
			//c, err := redigo.DialTimeout("tcp", host, 2*time.Second, 2*time.Second, 2*time.Second)
			//
			//if err != nil {
			//	logx.Warnf("err: %+v\n", err)
			//	RPair_Log.E("REPORTALARM=CONNECT_REDIS_ERROR,connect redis err1 host=%s,password=%s,connCount=%d,err=%s", host, password, connCount, err)
			//	return nil, err
			//}
			/*
				if password != "" {
					if _, err := c.Do("AUTH", password); err != nil {
						logx.Warnf("err: %+v\n", err)
						c.Close()
						RPair_Log.E("REPORTALARM=CONNECT_REDIS_ERROR,connect redis err2 host=%s,password=%s,connCount=%d,err=%s", host, password, connCount, err)
						// panic(err)
						return nil, err
					}
				}
				//_, err = c.Do("SELECT", config.RedisDb)
				RPair_Log.I("REPORTALARM=CONNECT_REDIS_OK,connect redis ok host=%s,password=%s,connCount=%d", host, password, connCount)
				return c, err
			*/
		},
	}
}

const NO_ERR = 0
const OCCUR_ERR = -100

func (rd *RedisOpt) getReturnCode(err error) int {
	if err != nil {
		RPair_Log.E("err: %s\n", err.Error())
		return OCCUR_ERR
	}

	return NO_ERR
}

func (rd *RedisOpt) Expire(k string, timeout int) int {
	c := rd.RedisPoll.Get()
	defer c.Close()

	_, err := c.Do("EXPIRE", k, timeout)

	return rd.getReturnCode(err)
}

func (rd *RedisOpt) Setex(k string, timeout int, v string) int {
	c := rd.RedisPoll.Get()
	defer c.Close()

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	v = timestamp + ":" + v

	_, err := c.Do("SETEX", k, timeout, v)
	return rd.getReturnCode(err)
}

func (rd *RedisOpt) GetSet(k string, v string) (int, string) {
	c := rd.RedisPoll.Get()
	defer c.Close()

	v, err := redigo.String(c.Do("GETSET", k, v))

	return rd.getReturnCode(err), v
}

func (rd *RedisOpt) Get(k string) (int, string) {
	c := rd.RedisPoll.Get()
	defer c.Close()

	v, err := redigo.String(c.Do("GET", k))
	return rd.getReturnCode(err), v
}

func (rd *RedisOpt) Get0(k string, useTimeStamp int) (int, string) {
	c := rd.RedisPoll.Get()
	defer c.Close()

	v, err := redigo.String(c.Do("GET", k))

	if useTimeStamp == 1 {
		_, vS := spiltStr(&v)
		return rd.getReturnCode(err), vS
	} else {
		return rd.getReturnCode(err), v
	}
}

func (rd *RedisOpt) Set(k string, v string) int {
	c := rd.RedisPoll.Get()
	defer c.Close()

	_, err := c.Do("SET", k, v)

	return rd.getReturnCode(err)
}

func (rd *RedisOpt) Del(k string) int {
	c := rd.RedisPoll.Get()
	defer c.Close()

	_, err := c.Do("DEL", k)

	return rd.getReturnCode(err)
}

func (rd *RedisOpt) Hset(k string, v map[string]string) int {
	c := rd.RedisPoll.Get()
	defer c.Close()

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	var err error

	for k1, v1 := range v {
		v1 = timestamp + ":" + v1
		_, err = c.Do("HSET", k, k1, v1)
	}

	return rd.getReturnCode(err)
}

func spiltStr(pOrg *string) (string, string) {
	lenth := len(*pOrg)
	if lenth < 11 {
		return "", *pOrg
	}

	timestamp := (*pOrg)[0:11]
	if strings.HasSuffix(timestamp, ":") {
		return timestamp[0:10], (*pOrg)[11:lenth]
	}

	return "", *pOrg
}

func (rd *RedisOpt) Hget(k string, k1 string) (int, string) {
	c := rd.RedisPoll.Get()
	defer c.Close()

	v, err := redigo.String(c.Do("HGET", k, k1))

	_, vS := spiltStr(&v)

	return rd.getReturnCode(err), vS
}

func (rd *RedisOpt) Hgetall(k string) (int, map[string]string) {
	c := rd.RedisPoll.Get()
	defer c.Close()

	dict := make(map[string]string)

	v, err := redigo.Strings(c.Do("HGETALL", k))
	for i := 0; i < len(v); i += 2 {
		kT := v[i]
		vT := v[i+1]
		_, v0_s := spiltStr(&vT)
		dict[kT] = v0_s
	}
	return rd.getReturnCode(err), dict
}

func (rd *RedisOpt) Hdel(k string, k1 string) int {
	c := rd.RedisPoll.Get()
	defer c.Close()

	_, err := c.Do("HDEL", k, k1)

	return rd.getReturnCode(err)
}

func (rd *RedisOpt) Hlen(k string) (int, int) {
	c := rd.RedisPoll.Get()
	defer c.Close()

	v, err := redigo.Int(c.Do("HLEN", k))

	return rd.getReturnCode(err), v
}

func (rd *RedisOpt) Hkeys(k string) (int, []string) {
	c := rd.RedisPoll.Get()
	defer c.Close()

	v, err := redigo.Strings(c.Do("HKEYS", k))
	return rd.getReturnCode(err), v
}

func (rd *RedisOpt) Keys(prefix string) (list []string, err error) {
	c, err := rd.RedisPoll.Dial()
	if err != nil {
		logx.Warnf("err1: %+v\n", err)
		return list, err
	}

	var getScript = redigo.NewScript(1, `return redis.call('keys', KEYS[1])`)

	keys, err := redigo.Strings(getScript.Do(c, prefix))
	if err != nil {
		logx.Warnf("err: %+v\n", err)
		return list, err
	}

	return keys, nil
}

func (rd *RedisOpt) Zadd(k string, score int, member string) int {
	c := rd.RedisPoll.Get()
	defer c.Close()

	_, err := c.Do("ZADD", k, score, member)

	return rd.getReturnCode(err)
}

func (rd *RedisOpt) Zrem(k string, member string) int {
	c := rd.RedisPoll.Get()
	defer c.Close()

	_, err := c.Do("ZREM", k, member)

	return rd.getReturnCode(err)
}

func (rd *RedisOpt) Zscore(k string, member string) (int, int) {
	c := rd.RedisPoll.Get()
	defer c.Close()

	v, err := redigo.Int(c.Do("ZSCORE", k, member))
	return rd.getReturnCode(err), v
}

func (rd *RedisOpt) Zrange(k string, start int, end int) (int, []string) {
	c := rd.RedisPoll.Get()
	defer c.Close()

	v, err := redigo.Strings(c.Do("ZRANGE", k, start, end))

	return rd.getReturnCode(err), v
}

//list
func (rd *RedisOpt) Lpush(k string, v string) int {
	c := rd.RedisPoll.Get()
	defer c.Close()

	reply, err := c.Do("LPUSH", k, v)
	if err != nil {
		logx.Warnf("reply: %+v err %v\n", reply, err)
	}

	return rd.getReturnCode(err)
}

func (rd *RedisOpt) Rpush(k string, v string) int {
	c := rd.RedisPoll.Get()
	defer c.Close()

	_, err := c.Do("RPUSH", k, v)

	return rd.getReturnCode(err)
}

func (rd *RedisOpt) Lpop(k string) (int, string) {
	c := rd.RedisPoll.Get()
	defer c.Close()

	v, err := redigo.String(c.Do("LPOP", k))

	return rd.getReturnCode(err), v
}

func (rd *RedisOpt) Rpop(k string) (int, string) {
	c := rd.RedisPoll.Get()
	defer c.Close()

	v, err := redigo.String(c.Do("RPOP", k))
	return rd.getReturnCode(err), v
}

func (rd *RedisOpt) Llen(k string) (int, int) {
	c := rd.RedisPoll.Get()
	defer c.Close()

	v, err := redigo.Int(c.Do("LLEN", k))

	return rd.getReturnCode(err), v
}

func (rd *RedisOpt) Lrange(k string, start int, end int) (int, []string) {
	c := rd.RedisPoll.Get()
	defer c.Close()

	v, err := redigo.Strings(c.Do("LRANGE", k, start, end))

	return rd.getReturnCode(err), v
}

//publish & subscribe
func (rd *RedisOpt) Publish(k string, v string) int {
	c := rd.RedisPoll.Get()
	defer c.Close()

	_, err := c.Do("Publish", k, v)
	return rd.getReturnCode(err)
}
