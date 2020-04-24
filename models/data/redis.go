package data

import (
	"errors"
	"female/lib/tools"

	"github.com/BurntSushi/toml"
	"github.com/gomodule/redigo/redis"
)

type redisConf struct {
	IP   string `toml:"ip"`
	Port string `toml:"port"`
}

var redisConfInfo = new(redisConf)

func init() {
	toml.DecodeFile(tools.GetCurrentDirectory()+"/conf/redis.toml", redisConfInfo)
}

func SetKey(key, value string) error {
	con, err := redis.Dial("tcp", redisConfInfo.IP+":"+redisConfInfo.Port)
	if err != nil {
		return err
	}
	defer con.Close()
	_, err = con.Do("set", key, value)
	if err != nil {
		return err
	}
	return nil
}

func GetKey(key string) (string, error) {
	con, err := redis.Dial("tcp", redisConfInfo.IP+":"+redisConfInfo.Port)
	if err != nil {
		return "", err
	}
	defer con.Close()
	// 判断key是否存在
	is, err := redis.Bool(con.Do("EXISTS", key))
	if err != nil {
		return "", err
	}
	if is {
		val, err := redis.String(con.Do("GET", key))
		if err != nil {
			return "", err
		}
		return val, nil
	}
	return "", errors.New("the key " + key + " is not exist")
}

func SetKeyExpire(key string, dur int) error {
	con, err := redis.Dial("tcp", redisConfInfo.IP+":"+redisConfInfo.Port)
	if err != nil {
		return err
	}
	defer con.Close()
	_, err = con.Do("EXPIRE", key, dur)
	if err != nil {
		return err
	}
	return nil
}

func DeleteKey(key string) error {
	con, err := redis.Dial("tcp", redisConfInfo.IP+":"+redisConfInfo.Port)
	if err != nil {
		return err
	}
	defer con.Close()
	con.Do("DEL", key)
	return nil
}
