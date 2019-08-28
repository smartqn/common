package libredis

/*************
1:20160302 带连接池连接维护功能，redis双写双读操作
2:注意，修改了 开源框架redigo源代码，在没有连接时，直接返回。更新golib库源代码需注意！！！！
3:函数第一个返回值都代表操作情况(0:主备都操作成功，-1:主失败备成功, -2:主成功备失败,-100:主备都失败)
4:key不存在 string类型返回空字符串
5：压力测试结论
*************/

import (
	"container/list"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aiwuTech/fileLogger"
	redigo "github.com/garyburd/redigo/redis"
	"github.com/linnv/logx"
)

const (
	RPair_RedisPass = ""

	RPair_ErrReconnCount = 10
)

var RPair_Log *fileLogger.FileLogger

var RPair_RPUSH_ErrBuf_Len = 102400
var RPair_RPUSH_KVSplit = "20160401_RPUSH_KVSplit_1052"

type RedisPair struct {
	RedisPoll0 *redigo.Pool
	RedisSvr0  string //"127.0.0.1:6379"

	RedisPoll1 *redigo.Pool
	RedisSvr1  string

	rPair_RPUSH_ErrBuf_List list.List
	rPair_RPUSH_RWMux       sync.RWMutex
	rPair_LRANGE_RANDOM     int
}

func (pair *RedisPair) Close() {
	pair.RedisPoll0.Close()
	pair.RedisPoll1.Close()
}

func (pair *RedisPair) getReturnCode(err0 error, err1 error) int {
	// check is ErrNil
	if err0 == redigo.ErrNil {
		err0 = nil
	}

	if err1 == redigo.ErrNil {
		err1 = nil
	}

	if err0 == nil && err1 == nil {
		return 0 // c0 c1 ok
	} else if err1 == nil {
		return -1 // c0 err
	} else if err0 == nil {
		return -2 // c1 err
	} else {
		RPair_Log.E("REPORTALARM=ALL_REDIS_ERROR,err0=%s,err1=%s", err0, err1)
		return -100 // c0 c1 err
	}
}

func (pair *RedisPair) Expire(k string, timeout int) int {

	c0 := pair.RedisPoll0.Get2(0)
	c1 := pair.RedisPoll1.Get2(0)

	defer func() {
		c0.Close()
		c1.Close()
	}()

	_, err0 := c0.Do("EXPIRE", k, timeout)
	_, err1 := c1.Do("EXPIRE", k, timeout)

	return pair.getReturnCode(err0, err1)
}

func (pair *RedisPair) Setex(k string, timeout int, v string) int {

	c0 := pair.RedisPoll0.Get2(0)
	c1 := pair.RedisPoll1.Get2(0)

	defer func() {
		c0.Close()
		c1.Close()
	}()

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	v = timestamp + ":" + v

	_, err0 := c0.Do("SETEX", k, timeout, v)
	_, err1 := c1.Do("SETEX", k, timeout, v)

	return pair.getReturnCode(err0, err1)
}

func (pair *RedisPair) GetSet(k string, v string) (int, string) {

	c0 := pair.RedisPoll0.Get2(0)
	c1 := pair.RedisPoll1.Get2(0)

	defer func() {
		c0.Close()
		c1.Close()
	}()

	v0, err0 := redigo.String(c0.Do("GETSET", k, v))
	v1, err1 := redigo.String(c1.Do("GETSET", k, v))

	v0_t, v0_s := spiltStr(&v0)
	v1_t, v1_s := spiltStr(&v1)

	if v0_t >= v1_t {
		return pair.getReturnCode(err0, err1), v0_s
	} else {
		return pair.getReturnCode(err0, err1), v1_s
	}

}

func (pair *RedisPair) Set(k string, v string) int {

	c0 := pair.RedisPoll0.Get2(0)
	c1 := pair.RedisPoll1.Get2(0)

	defer func() {
		c0.Close()
		c1.Close()
	}()

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	v = timestamp + ":" + v

	_, err0 := c0.Do("SET", k, v)
	_, err1 := c1.Do("SET", k, v)

	return pair.getReturnCode(err0, err1)

}

func (pair *RedisPair) Del(k string) int {

	c0 := pair.RedisPoll0.Get2(0)
	c1 := pair.RedisPoll1.Get2(0)

	defer func() {
		c0.Close()
		c1.Close()
	}()

	_, err0 := c0.Do("DEL", k)
	_, err1 := c1.Do("DEL", k)

	return pair.getReturnCode(err0, err1)

}

func (pair *RedisPair) Hset(k string, v map[string]string) int {

	c0 := pair.RedisPoll0.Get2(0)
	c1 := pair.RedisPoll1.Get2(0)

	defer func() {
		c0.Close()
		c1.Close()
	}()

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	var err0, err1 error

	for k1, v1 := range v {
		v1 = timestamp + ":" + v1
		_, err0 = c0.Do("HSET", k, k1, v1)
		_, err1 = c1.Do("HSET", k, k1, v1)
	}

	return pair.getReturnCode(err0, err1)
}

func (pair *RedisPair) Hget(k string, k1 string) (int, string) {

	c0 := pair.RedisPoll0.Get2(0)
	c1 := pair.RedisPoll1.Get2(0)

	defer func() {
		c0.Close()
		c1.Close()
	}()

	v0, err0 := redigo.String(c0.Do("HGET", k, k1))
	v1, err1 := redigo.String(c1.Do("HGET", k, k1))

	v0_t, v0_s := spiltStr(&v0)
	v1_t, v1_s := spiltStr(&v1)

	if v0_t >= v1_t {
		return pair.getReturnCode(err0, err1), v0_s
	} else {
		return pair.getReturnCode(err0, err1), v1_s
	}
}

func (pair *RedisPair) Hgetall(k string) (int, map[string]string) {

	c0 := pair.RedisPoll0.Get2(0)
	c1 := pair.RedisPoll1.Get2(0)

	defer func() {
		c0.Close()
		c1.Close()
	}()

	dict := make(map[string]string)

	v0, err0 := redigo.Strings(c0.Do("HGETALL", k))
	v1, err1 := redigo.Strings(c1.Do("HGETALL", k))

	if len(v0) >= len(v1) {
		for i := 0; i < len(v0); i += 2 {
			k0_t := v0[i]
			v0_t := v0[i+1]
			_, v0_s := spiltStr(&v0_t)
			dict[k0_t] = v0_s
		}
		return pair.getReturnCode(err0, err1), dict
	} else {
		for i := 0; i < len(v1); i += 2 {
			k1_t := v1[i]
			v1_t := v1[i+1]
			_, v1_s := spiltStr(&v1_t)
			dict[k1_t] = v1_s
		}
		return pair.getReturnCode(err0, err1), dict
	}
}

func (pair *RedisPair) Hdel(k string, k1 string) int {
	c0 := pair.RedisPoll0.Get2(0)
	c1 := pair.RedisPoll1.Get2(0)

	defer func() {
		c0.Close()
		c1.Close()
	}()

	_, err0 := c0.Do("HDEL", k, k1)
	_, err1 := c1.Do("HDEL", k, k1)

	return pair.getReturnCode(err0, err1)
}

func (pair *RedisPair) Hlen(k string) (int, int) {
	c0 := pair.RedisPoll0.Get2(0)
	c1 := pair.RedisPoll1.Get2(0)

	defer func() {
		c0.Close()
		c1.Close()
	}()

	v0, err0 := redigo.Int(c0.Do("HLEN", k))
	v1, err1 := redigo.Int(c1.Do("HLEN", k))

	if v0 >= v1 {
		return pair.getReturnCode(err0, err1), v0
	} else {
		return pair.getReturnCode(err0, err1), v1
	}

}

func (pair *RedisPair) Hkeys(k string) (int, []string) {
	c0 := pair.RedisPoll0.Get2(0)
	c1 := pair.RedisPoll1.Get2(0)

	defer func() {
		c0.Close()
		c1.Close()
	}()

	v0, err0 := redigo.Strings(c0.Do("HKEYS", k))
	v1, err1 := redigo.Strings(c1.Do("HKEYS", k))

	if len(v0) >= len(v1) {
		return pair.getReturnCode(err0, err1), v0
	} else {
		return pair.getReturnCode(err0, err1), v1
	}

}

func spiltStr(pOrg *string) (string, string) {
	len := len(*pOrg)
	if len < 11 {
		return "", *pOrg
	}

	timestamp := (*pOrg)[0:11]
	if strings.HasSuffix(timestamp, ":") {
		return timestamp[0:10], (*pOrg)[11:len]
	}

	return "", *pOrg
}

func (pair *RedisPair) Keys(prefix string) (list []string, err error) {
	c1, err1 := pair.RedisPoll0.Dial()
	if err1 != nil {
		logx.Warnf("err1: %+v\n", err1)
		return list, err1
	}

	c2, err2 := pair.RedisPoll0.Dial()
	if err2 != nil {
		logx.Warnf("err2: %+v\n", err2)

		return list, err2
	}

	var getScript = redigo.NewScript(1, `return redis.call('keys', KEYS[1])`)

	keys1, err := redigo.Strings(getScript.Do(c1, prefix))
	if err != nil {
		logx.Warnf("err: %+v\n", err)
		return list, err
	}

	keys2, err := redigo.Strings(getScript.Do(c2, prefix))
	if err != nil {
		logx.Warnf("err: %+v\n", err)
		return list, err
	}
	if len(keys1) != len(keys2) {
		err = fmt.Errorf("data of redis pairs broken")
		if len(keys1) > len(keys2) {
			return keys2, err
		}
	}
	return keys1, nil
}

func (pair *RedisPair) Get(k string) (int, string) {
	return pair.Get0(k, 1)
}

//check timestamp
func (pair *RedisPair) Get0(k string, useTimestamp int) (int, string) {

	c0 := pair.RedisPoll0.Get2(0)
	c1 := pair.RedisPoll1.Get2(0)

	defer func() {
		c0.Close()
		c1.Close()
	}()

	v0, err0 := redigo.String(c0.Do("GET", k))
	v1, err1 := redigo.String(c1.Do("GET", k))

	if useTimestamp == 1 {
		v0_t, v0_s := spiltStr(&v0)
		v1_t, v1_s := spiltStr(&v1)

		if v0_t >= v1_t {
			return pair.getReturnCode(err0, err1), v0_s
		} else {
			return pair.getReturnCode(err0, err1), v1_s
		}
	} else {
		if v0 != "" {
			return pair.getReturnCode(err0, err1), v0
		} else {
			return pair.getReturnCode(err0, err1), v1
		}
	}

}

func (pair *RedisPair) Zadd(k string, score int, member string) int {

	c0 := pair.RedisPoll0.Get2(0)
	c1 := pair.RedisPoll1.Get2(0)

	defer func() {
		c0.Close()
		c1.Close()
	}()

	_, err0 := c0.Do("ZADD", k, score, member)
	_, err1 := c1.Do("ZADD", k, score, member)

	return pair.getReturnCode(err0, err1)

}

func (pair *RedisPair) Zrem(k string, member string) int {

	c0 := pair.RedisPoll0.Get2(0)
	c1 := pair.RedisPoll1.Get2(0)

	defer func() {
		c0.Close()
		c1.Close()
	}()

	_, err0 := c0.Do("ZREM", k, member)
	_, err1 := c1.Do("ZREM", k, member)

	return pair.getReturnCode(err0, err1)

}

func (pair *RedisPair) Zscore(k string, member string) (int, int) {

	c0 := pair.RedisPoll0.Get2(0)
	c1 := pair.RedisPoll1.Get2(0)

	defer func() {
		c0.Close()
		c1.Close()
	}()

	v := 0
	v0, err0 := redigo.Int(c0.Do("ZSCORE", k, member))
	v1, err1 := redigo.Int(c1.Do("ZSCORE", k, member))

	if err0 == nil {
		v = v0
	} else if err1 == nil {
		v = v1
	}

	return pair.getReturnCode(err0, err1), v

}

func (pair *RedisPair) Zrange(k string, start int, end int) (int, []string) {

	c0 := pair.RedisPoll0.Get2(0)
	c1 := pair.RedisPoll1.Get2(0)

	defer func() {
		c0.Close()
		c1.Close()
	}()

	v0, err0 := redigo.Strings(c0.Do("ZRANGE", k, start, end))
	if err0 == nil && len(v0) != 0 {
		return pair.getReturnCode(err0, nil), v0
	}

	v1, err1 := redigo.Strings(c1.Do("ZRANGE", k, start, end))

	return pair.getReturnCode(err0, err1), v1

}

func (pair *RedisPair) Zrangebyscore(k string, min int, max int) (int, []string) {

	c0 := pair.RedisPoll0.Get2(0)
	c1 := pair.RedisPoll1.Get2(0)

	defer func() {
		c0.Close()
		c1.Close()
	}()

	v0, err0 := redigo.Strings(c0.Do("ZRANGEBYSCORE", k, min, max))
	if err0 == nil && len(v0) != 0 {
		return pair.getReturnCode(err0, nil), v0
	}

	v1, err1 := redigo.Strings(c1.Do("ZRANGEBYSCORE", k, min, max))

	return pair.getReturnCode(err0, err1), v1

}

func (pair *RedisPair) Zremrangebyscore(k string, min int, max int) (int, int) {

	c0 := pair.RedisPoll0.Get2(0)
	c1 := pair.RedisPoll1.Get2(0)

	defer func() {
		c0.Close()
		c1.Close()
	}()

	v0, err0 := redigo.Int(c0.Do("ZREMRANGEBYSCORE", k, min, max))

	v1, err1 := redigo.Int(c1.Do("ZREMRANGEBYSCORE", k, min, max))

	if v0 >= v1 {
		return pair.getReturnCode(err0, err1), v0
	} else {
		return pair.getReturnCode(err0, err1), v1
	}

}

func (pair *RedisPair) Zcount(k string, min int, max int) (int, int) {

	c0 := pair.RedisPoll0.Get2(0)
	c1 := pair.RedisPoll1.Get2(0)

	defer func() {
		c0.Close()
		c1.Close()
	}()

	v := 0
	v0, err0 := redigo.Int(c0.Do("ZCOUNT", k, min, max))
	v1, err1 := redigo.Int(c1.Do("ZCOUNT", k, min, max))

	if err0 == nil && err1 == nil {
		if v1 > v0 {
			v = v1
		} else {
			v = v0
		}
	} else if err0 == nil {
		v = v0
	} else if err1 == nil {
		v = v1
	}

	return pair.getReturnCode(err0, err1), v

}

func (pair *RedisPair) Lpush(k string, v string) int {

	c0 := pair.RedisPoll0.Get2(0)
	c1 := pair.RedisPoll1.Get2(0)

	defer func() {
		c0.Close()
		c1.Close()
	}()

	// logx.Debugf("k: %+v v :%v\n", k, v)
	err, err0 := c0.Do("LPUSH", k, v)
	if err != nil {
		logx.Warnf("err: %+v err0 %v\n", err, err0)
	}

	err, err1 := c1.Do("LPUSH", k, v)
	if err != nil {
		logx.Warnf("err: %+v err1 %v\n", err, err1)
	}

	return pair.getReturnCode(err0, err1)

}

func (pair *RedisPair) Rpush(k string, v string) int {

	c0 := pair.RedisPoll0.Get2(0)
	c1 := pair.RedisPoll1.Get2(0)

	defer func() {
		c0.Close()
		c1.Close()
	}()

	_, err0 := c0.Do("RPUSH", k, v)
	_, err1 := c1.Do("RPUSH", k, v)

	return pair.getReturnCode(err0, err1)

}

//zhichen 对双机都失败的，提供一个缓冲区，后台程序自动重试
func (pair *RedisPair) Rpush_WithRetry(k string, v string) int {

	c0 := pair.RedisPoll0.Get2(0)
	c1 := pair.RedisPoll1.Get2(0)

	defer func() {
		c0.Close()
		c1.Close()
	}()

	_, err0 := c0.Do("RPUSH", k, v)
	_, err1 := c1.Do("RPUSH", k, v)

	res := pair.getReturnCode(err0, err1)

	if res == -100 {
		pair.rPair_RPUSH_RWMux.Lock()
		str := k + RPair_RPUSH_KVSplit + v
		pair.rPair_RPUSH_ErrBuf_List.PushBack(str)

		RPair_Log.I("Rpush_WithRetry str=%s", str)
		pair.rPair_RPUSH_RWMux.Unlock()
	}

	return res

}

func (pair *RedisPair) Lpop(k string) (int, string) {

	c0 := pair.RedisPoll0.Get2(0)
	c1 := pair.RedisPoll1.Get2(0)

	defer func() {
		c0.Close()
		c1.Close()
	}()

	v0, err0 := redigo.String(c0.Do("LPOP", k))
	v1, err1 := redigo.String(c1.Do("LPOP", k))

	if err0 == nil {
		return pair.getReturnCode(err0, err1), v0
	} else {
		return pair.getReturnCode(err0, err1), v1
	}

}

func (pair *RedisPair) Rpop(k string) (int, string) {

	c0 := pair.RedisPoll0.Get2(0)
	c1 := pair.RedisPoll1.Get2(0)

	defer func() {
		c0.Close()
		c1.Close()
	}()

	v0, err0 := redigo.String(c0.Do("RPOP", k))
	v1, err1 := redigo.String(c1.Do("RPOP", k))

	if err0 == nil {
		return pair.getReturnCode(err0, err1), v0
	} else {
		return pair.getReturnCode(err0, err1), v1
	}

}

func (pair *RedisPair) Llen(k string) (int, int) {

	c0 := pair.RedisPoll0.Get2(0)
	c1 := pair.RedisPoll1.Get2(0)

	defer func() {
		c0.Close()
		c1.Close()
	}()

	v0, err0 := redigo.Int(c0.Do("LLEN", k))
	v1, err1 := redigo.Int(c1.Do("LLEN", k))

	if v0 >= v1 {
		return pair.getReturnCode(err0, err1), v0
	} else {
		return pair.getReturnCode(err0, err1), v1
	}

}

func (pair *RedisPair) Lrange(k string, start int, end int) (int, []string) {

	c0 := pair.RedisPoll0.Get2(0)
	c1 := pair.RedisPoll1.Get2(0)

	defer func() {
		c0.Close()
		c1.Close()
	}()

	v0, err0 := redigo.Strings(c0.Do("LRANGE", k, start, end))
	v1, err1 := redigo.Strings(c1.Do("LRANGE", k, start, end))

	if len(v0) >= len(v1) {
		return pair.getReturnCode(err0, err1), v0
	} else {
		return pair.getReturnCode(err0, err1), v1
	}

}

func (pair *RedisPair) lrange_0(k string, start int, end int) (error, []string) {

	c0 := pair.RedisPoll0.Get2(0)
	defer func() {
		c0.Close()
	}()

	v0, err0 := redigo.Strings(c0.Do("LRANGE", k, start, end))

	return err0, v0
}

func (pair *RedisPair) lrange_1(k string, start int, end int) (error, []string) {

	c1 := pair.RedisPoll1.Get2(0)
	defer func() {
		c1.Close()
	}()

	v1, err1 := redigo.Strings(c1.Do("LRANGE", k, start, end))

	return err1, v1
}

func (pair *RedisPair) Lrange_WithReadRandom(k string, start int, end int) (int, []string) {

	pair.rPair_LRANGE_RANDOM = pair.rPair_LRANGE_RANDOM + 1
	if pair.rPair_LRANGE_RANDOM > 10240 {
		pair.rPair_LRANGE_RANDOM = 1
	}

	if pair.rPair_LRANGE_RANDOM%2 == 0 {
		err0, v0 := pair.lrange_0(k, start, end)
		if err0 == nil {
			return pair.getReturnCode(err0, nil), v0
		} else {
			err1, v1 := pair.lrange_1(k, start, end)
			return pair.getReturnCode(err0, err1), v1
		}

	} else {
		err1, v1 := pair.lrange_1(k, start, end)
		if err1 == nil {
			return pair.getReturnCode(nil, err1), v1
		} else {
			err0, v0 := pair.lrange_0(k, start, end)
			return pair.getReturnCode(err0, err1), v0
		}

	}
}

func getMsgSeq(pOrg *string, getlen int) string {
	len := len(*pOrg)
	if len < getlen {
		return ""
	}

	substr := (*pOrg)[0:getlen]

	seqindex := strings.Index(substr, "msgSeq")
	if seqindex <= 0 {
		return ""
	}

	seqindex = seqindex + 6

	substr = substr[seqindex:getlen]
	tagindex := strings.Index(substr, ",")

	if tagindex <= 0 {
		return ""
	}

	mseq := substr[0:tagindex]

	mseq = strings.Replace(mseq, "\"", "", -1)
	mseq = strings.Replace(mseq, ":", "", -1)

	return mseq
}

func (pair *RedisPair) Lrange_WithReadAll(k string, start int, end int) (int, []string) {

	c0 := pair.RedisPoll0.Get2(0)
	c1 := pair.RedisPoll1.Get2(0)

	v0, err0 := redigo.Strings(c0.Do("LRANGE", k, start, end))
	v1, err1 := redigo.Strings(c1.Do("LRANGE", k, start, end))

	c0.Close()
	c1.Close()

	allm := make(map[string]string)
	if err0 == nil {
		for _, m := range v0 {
			seq := getMsgSeq(&m, 50)
			if seq != "" {
				allm[seq] = m
			}
		}
	}

	if err1 == nil {
		for _, m := range v1 {
			seq := getMsgSeq(&m, 50)
			if seq != "" {
				allm[seq] = m
			}
		}
	}

	var keys []string
	for k := range allm {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var v []string
	for _, k := range keys {
		v = append(v, allm[k])
	}

	return pair.getReturnCode(err0, err1), v

}

//自动重入两条redis都Rpush失败的msg数据
func autoReIRpushMsg(pair *RedisPair, checkInter time.Duration) {

	for {
		pair.rPair_RPUSH_RWMux.Lock()
		//RPair_Log.I("autoReIRpushMsg len=%d,beg", RPair_RPUSH_ErrBuf_List.Len())
		for i := 0; i < 1000; i++ {
			elem := pair.rPair_RPUSH_ErrBuf_List.Front()
			if elem == nil {
				break
			}

			kb := strings.Split(elem.Value.(string), RPair_RPUSH_KVSplit)

			if len(kb) == 2 {
				res := pair.Rpush(kb[0], kb[1])

				//RPair_Log.I("autoReIRpushMsg Rpush res=%s,kb=%s",res,kb )

				if res == -100 {
					break
				}
			}

			pair.rPair_RPUSH_ErrBuf_List.Remove(elem)

			//RPair_Log.I("autoReIRpushMsg Remove=%s,len=%d", elem.Value.(string),RPair_RPUSH_ErrBuf_List.Len())
		}

		//强制清空
		if pair.rPair_RPUSH_ErrBuf_List.Len() > RPair_RPUSH_ErrBuf_Len {
			RPair_Log.E("REPORTALARM=AUTO_RERPUSH_ERROR,len=%s", pair.rPair_RPUSH_ErrBuf_List.Len())
			pair.rPair_RPUSH_ErrBuf_List.Init()
		}

		pair.rPair_RPUSH_RWMux.Unlock()
		time.Sleep(checkInter * time.Second)
	}

}

func autoCheckPoll(pollname string, poll *redigo.Pool, checkInter time.Duration, checkCount int) {
	if poll == nil {
		return
	}
	//check conn stat
	isStart := true
	monitorCount := checkCount / 3
	if monitorCount <= 0 {
		monitorCount = 1
	}
	for {
		// RPair_Log.I("autoCheckPoll pollname=%s,beg", pollname)

		var connArray []redigo.Conn
		if isStart {
			connArray = make([]redigo.Conn, checkCount)
			isStart = false
		} else {
			connArray = make([]redigo.Conn, monitorCount)
		}

		//get
		for index, _ := range connArray {
			connArray[index] = poll.Get()
		}

		//check
		errConnNum := 0
		allConnNum := len(connArray)
		for index, c := range connArray {
			if c.Err() == nil {
				_, err := c.Do("PING")
				if err == nil {
					//Log.I("autoCheckPoll pollname=%s,index=%d end", pollname,index)
				} else {
					RPair_Log.E("autoCheckPoll pollname=%s,index=%d,err=%s end", pollname, err, index)
					c.Close()
					connArray[index] = poll.Get()
					errConnNum++
				}
			} else {
				RPair_Log.E("autoCheckPoll pollname=%s,index=%d,err=%s end", pollname, c.Err(), index)
				errConnNum++
				//c.Close()
				//connArray[index] = poll.Get() // no need
			}
		}

		if errConnNum == allConnNum {
			isStart = true //reconnet all
		}

		//close
		for index := 0; index < len(connArray); index++ {
			connArray[index].Close()
		}
		// RPair_Log.I("autoCheckPoll pollname=%s,end", pollname)
		time.Sleep(checkInter * time.Second)
	}

}
func InitWithPasswd(mainSvr string, backSvr string, connCount int, logPath, password string) *RedisPair {
	RPair_Log = fileLogger.NewSizeLogger(logPath, "redispair.Log", "", 100, 100, fileLogger.MB, 300, 10000)
	RPair_Log.SetLogLevel(fileLogger.INFO) //trace Log will not be print

	pair := &RedisPair{
		RedisPoll0:          nil,
		RedisSvr0:           mainSvr,
		RedisPoll1:          nil,
		RedisSvr1:           backSvr,
		rPair_LRANGE_RANDOM: 1,
	}

	pair.RedisPoll0 = initRedis(pair.RedisSvr0, password, connCount)
	pair.RedisPoll1 = initRedis(pair.RedisSvr1, password, connCount)

	go autoCheckPoll("poll0", pair.RedisPoll0, 60, connCount)
	go autoCheckPoll("poll1", pair.RedisPoll1, 60, connCount)

	go autoReIRpushMsg(pair, 1)

	time.Sleep(2 * time.Second)

	return pair
}

func Init(mainSvr string, backSvr string, connCount int, logPath string) *RedisPair {
	RPair_Log = fileLogger.NewSizeLogger(logPath, "redispair.Log", "", 100, 100, fileLogger.MB, 300, 10000)
	RPair_Log.SetLogLevel(fileLogger.INFO) //trace Log will not be print

	pair := &RedisPair{
		RedisPoll0:          nil,
		RedisSvr0:           mainSvr,
		RedisPoll1:          nil,
		RedisSvr1:           backSvr,
		rPair_LRANGE_RANDOM: 1,
	}

	pair.RedisPoll0 = initRedis(pair.RedisSvr0, RPair_RedisPass, connCount)
	pair.RedisPoll1 = initRedis(pair.RedisSvr1, RPair_RedisPass, connCount)

	go autoCheckPoll("poll0", pair.RedisPoll0, 60, connCount)
	go autoCheckPoll("poll1", pair.RedisPoll1, 60, connCount)

	go autoReIRpushMsg(pair, 1)

	time.Sleep(2 * time.Second)

	return pair
}

func initRedis(host string, password string, connCount int) *redigo.Pool {

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
		TestOnBorrow: nil,
		Dial: func() (redigo.Conn, error) {
			RPair_Log.I("connect redis begin host=%s,password=%s,connCount=%d", host, password, connCount)
			c, err := redigo.DialTimeout("tcp", host, 2*time.Second, 2*time.Second, 2*time.Second)

			if err != nil {
				logx.Warnf("err: %+v\n", err)
				RPair_Log.E("REPORTALARM=CONNECT_REDIS_ERROR,connect redis err1 host=%s,password=%s,connCount=%d,err=%s", host, password, connCount, err)
				return nil, err
			}
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
		},
	}
}
