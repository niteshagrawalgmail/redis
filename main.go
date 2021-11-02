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
	"github.com/go-redis/redis/v8"
	"net/http"
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

func retriveData(w http.ResponseWriter, r *http.Request){
	k:= getKeyToRetrive(r)
	v := getDataFromredis(k)

	w.WriteHeader(200)
	fmt.Fprintf(w, v)
}

func storeData(w http.ResponseWriter, r *http.Request) {

	k, v := getKeyValToStore(r)
	persistDataInRedis(k,v)

	w.WriteHeader(200)
	fmt.Fprintf(w, "Data stored in redis")

}

func getKeyToRetrive(r *http.Request) (string){
	query := r.URL.Query()
	key := query["key"][0]
	return key
}

func getKeyValToStore(r *http.Request) (string, string){
	query := r.URL.Query()
	key := query["key"][0]
	value := query["value"][0]

	return key, value
}



func persistDataInRedis(k string, v string){
		err := getRedisClient().Set(ctx, k, v, 0).Err()
		if err != nil {
			panic(err)
		}
}

func getDataFromredis(k string) string {
	value, err := getRedisClient().Get(ctx,k).Result()
	if err != nil{
		panic(err)
	}
	return value
}

func getRedisClient() redis.UniversalClient{
	if client != nil{
		return client
	}else{
		client = redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs: []string{"rfs-redisfailover:26379"},
			MasterName: "mymaster",
		})

		return client
	}
}