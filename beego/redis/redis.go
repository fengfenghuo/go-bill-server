package redis

import (
	"reflect"

	"github.com/gomodule/redigo/redis"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("redis")
)

type RedisServer struct {
	redis redis.Conn
}

func NewDataRedis(url string) (*RedisServer, error) {
	c, err := redis.Dial("tcp", url)
	if err != nil {
		return nil, err
	}

	return &RedisServer{redis: c}, nil
}

func (r *RedisServer) Set(key interface{}, data interface{}) error {
	_, err := r.redis.Do("SET", key, data)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisServer) Get(key interface{}) (interface{}, error) {
	data, err := r.redis.Do("GET", key)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *RedisServer) HashSet(key interface{}, field interface{}, data interface{}) error {
	_, err := r.redis.Do("HSET", key, field, data)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisServer) HashGet(key interface{}, field interface{}) (interface{}, error) {
	data, err := r.redis.Do("HGET", key)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *RedisServer) HashAppend(key interface{}, field interface{}, data interface{}) error {
	storeData, err := r.HashGet(key, field)
	if err != nil {
		return err
	}

	var temp []interface{}

	v := reflect.ValueOf(storeData)
	if storeData == nil {
		temp = append(temp, data)
	} else if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			temp = append(temp, v.Index(i).Interface())
		}
		temp = append(temp, data)
	} else {
		temp = append(temp, storeData, data)
	}

	err = r.HashSet(key, field, temp)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisServer) HashIsExist(key interface{}, field interface{}) bool {
	result, err := r.redis.Do("HEXISTS", key, field)
	if err != nil {
		log.Error(err.Error())
		return false
	}

	if data, ok := result.(int); ok {
		if data == 1 {
			return true
		} else {
			return false
		}
	}

	return false
}

func (r *RedisServer) SetAppend(key interface{}, data interface{}) error {
	_, err := r.redis.Do("SADD", key, data)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisServer) SetDelete(key interface{}, data interface{}) error {
	_, err := r.redis.Do("SREM", key, data)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisServer) SetGet(key interface{}) (interface{}, error) {
	data, err := r.redis.Do("SMEMBERS", key)
	if err != nil {
		return nil, err
	}
	return data, nil
}
