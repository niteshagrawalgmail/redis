//package main
//
//import (
//	"fmt"
//	"context"
//	"github.com/go-redis/redis/v8"
//)
//var ctx = context.Background()
//func main() {
//	fmt.Println("Go Redis Tutorial")
//
//	rdb := redis.NewClient(&redis.Options{
//		Addr: "localhost:6379",
//		Password: "",
//		DB: 0,
//	})
//
//	redis.NewUniversalClient()
//
//	err := rdb.Set(ctx, "key", "value", 0).Err()
//	if err != nil {
//		panic(err)
//	}
//
//	val, err := rdb.Get(ctx, "key").Result()
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println("key", val)
//
//}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var client redis.UniversalClient

func main() {
	storeDataInRedis := http.HandlerFunc(storeData)
	retriveDataFromRedis := http.HandlerFunc(retriveData)
	http.Handle("/store", storeDataInRedis)
	http.Handle("/get", retriveDataFromRedis)
	http.ListenAndServe(":8080", nil)
}

func retriveData(w http.ResponseWriter, r *http.Request) {
	k := getKeyToRetrive(r)
	v := getDataFromredis(k)

	w.WriteHeader(200)
	fmt.Fprintf(w, v)
}

func storeData(w http.ResponseWriter, r *http.Request) {

	k, v := getKeyValToStore(r)
	log.Printf("Got the key:%v and the value:%v to store", k, v)
	persistDataInRedis(k, v)

	w.WriteHeader(200)
	fmt.Fprintf(w, "Data stored in redis")

}

func getKeyToRetrive(r *http.Request) string {
	query := r.URL.Query()
	key := query["key"][0]
	return key
}

func getKeyValToStore(r *http.Request) (string, string) {
	query := r.URL.Query()
	key := query["key"][0]
	value := query["value"][0]

	return key, value
}

func persistDataInRedis(k string, v string) {
	redisClient := getRedisClient()
	pingErr := redisClient.Ping(ctx).Err()

	if pingErr != nil {
		log.Println("Error occurred while pinging Redis")
		panic(pingErr)
	}

	log.Println("Successfully pinged Redis")

	fmt.Printf("Setting key %v and value %v in redis", k, v)
	err := redisClient.Set(ctx, k, v, 0).Err()
	if err != nil {
		log.Printf("Error occurred while setting key %v and value %v in redis", k, v)
		panic(err)
	}

	log.Printf("Successfully store key %v and value %v in redis", k, v)

}

func getDataFromredis(k string) string {

	redisClient := getRedisClient()
	pingErr := redisClient.Ping(ctx).Err()

	if pingErr != nil {
		log.Println("Error occurred while pinging Redis")
		panic(pingErr)
	}
	value, err := redisClient.Get(ctx, k).Result()
	if err != nil {
		log.Printf("Error occurred while getting key %v from redis", k)
		panic(err)
	}
	log.Printf("Successfully got key %v from redis", k)
	return value
}

func getRedisClient() redis.UniversalClient {
	if client != nil {
		log.Println("Redis client exists. Using the same")
		return client
	} else {
		isRedisHASetup := isRedisHA()
		if isRedisHASetup {
			log.Println("Redis is HA setup")
			redisUrl := getRedisURL()
			password := getPassword()
			if password != "" {
				client = redis.NewUniversalClient(&redis.UniversalOptions{
					Addrs:      []string{redisUrl},
					Password:   password,
					MasterName: "mymaster",
				})
			} else {
				client = redis.NewUniversalClient(&redis.UniversalOptions{
					Addrs:      []string{redisUrl},
					MasterName: "mymaster",
				})
			}

		} else {
			log.Println("Redis is non-HA setup")
			redisUrl := getRedisURL()
			password := getPassword()
			if password != "" {
				log.Println("Redis client with password")
				client = redis.NewUniversalClient(&redis.UniversalOptions{
					Addrs:    []string{redisUrl},
					Password: password,
				})
			} else {
				log.Println("Redis client without password")
				client = redis.NewUniversalClient(&redis.UniversalOptions{
					Addrs: []string{redisUrl},
				})
			}

		}

		log.Println("Creating a new Redis Client and using it")
		return client
	}
}

func getRedisURL() string {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	url := host + ":" + port
	return url
}

func isRedisHA() bool {
	isHA := os.Getenv("IS_HA")
	val, _ := strconv.ParseBool(isHA)
	return val
}

func getPassword() string {
	password := os.Getenv("REDIS_PASSWORD")
	return password
}
