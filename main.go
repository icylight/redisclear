package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	serviceredis "github.com/MiaoSiLa/missevan-go/service/redis"
	"gopkg.in/yaml.v2"
)

var conf struct {
	Redis serviceredis.Config `yaml:"redis"`
}

func LoadConfig(confPath string) {
	configFile, err := ioutil.ReadFile(confPath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(configFile, &conf)
	if err != nil {
		panic(err)
	}
}

func main() {
	LoadConfig("./config.yml")

	redis, err := serviceredis.NewRedisClient(&conf.Redis)
	if err != nil {
		panic(err)
	}
	keys, err := redis.Keys("live-service:im:room:*").Result()

	fmt.Println("keys:")
	for i := 0; i < len(keys); i++ {
		fmt.Println(keys[i])
	}
	fmt.Print("\n")
	fmt.Println("start clear redis")

	score := int64(1556913600)
	fmt.Printf("clear time: %s\n", time.Unix(score, 0))

	for i := 0; i < len(keys); i++ {
		_, err := redis.ZRemRangeByScore(keys[i], "-inf", strconv.FormatInt(score, 10)).Result()
		if err != nil {
			fmt.Println(err)
			fmt.Println("clear failed")
			os.Exit(1)
		}
	}

	fmt.Println("\nfinish")
	return
}
