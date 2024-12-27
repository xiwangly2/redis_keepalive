package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pelletier/go-toml"
)

type Config struct {
	Redis struct {
		Addr     string `toml:"addr"`
		Username string `toml:"username"`
		Password string `toml:"password"`
		DB       int    `toml:"db"`
	} `toml:"redis"`
}

var rdb *redis.Client

func main() {
	// 读取配置文件
	config, err := readConfig("config.toml")
	if err != nil {
		log.Fatalf("读取配置文件失败: %v\n", err)
	}

	// 创建Redis客户端
	rdb = redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Username: config.Redis.Username,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})

	// 测试连接
	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("连接失败: %v\n", err)
	} else {
		fmt.Println("连接成功:", pong)
	}

	// 增：写入数据
	writeData("timestamp", time.Now().Unix())

	// 查：读取数据
	readData("timestamp")

	// 改：更新数据
	updateData("timestamp", time.Now().Unix()+10000)

	// 查：读取更新后的数据
	readData("timestamp")

	// 删：删除数据
	deleteData("timestamp")

	// 查：验证删除后的数据
	readData("timestamp")
}

func readConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config Config
	err = toml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// 增：写入数据
func writeData(key string, value int64) {
	err := rdb.Set(context.Background(), key, value, 0).Err()
	if err != nil {
		log.Printf("写入数据失败: %v\n", err)
	} else {
		fmt.Printf("已写入数据: %s -> %d\n", key, value)
	}
}

// 查：读取数据
func readData(key string) {
	val, err := rdb.Get(context.Background(), key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			fmt.Printf("键 %s 不存在\n", key)
		} else {
			log.Printf("读取数据失败: %v\n", err)
		}
	} else {
		fmt.Printf("读取数据: %s -> %s\n", key, val)
	}
}

// 改：更新数据
func updateData(key string, newValue int64) {
	// 直接使用 Set 来更新数据
	err := rdb.Set(context.Background(), key, newValue, 0).Err()
	if err != nil {
		log.Printf("更新数据失败: %v\n", err)
	} else {
		fmt.Printf("已更新数据: %s -> %d\n", key, newValue)
	}
}

// 删：删除数据
func deleteData(key string) {
	err := rdb.Del(context.Background(), key).Err()
	if err != nil {
		log.Printf("删除数据失败: %v\n", err)
	} else {
		fmt.Printf("已删除数据: %s\n", key)
	}
}
