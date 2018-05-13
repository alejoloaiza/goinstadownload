package db

import (
	"database/sql"
	"fmt"

	"github.com/go-redis/redis"
)

var dbpostgre *sql.DB
var err error
var dbredis *redis.Client

func DBConnectRedis() {
	dbredis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	//fmt.Println(">>>>>>>>>>>>>>>>> Successfully connected to Database <<<<<<<<<<<<<<<<<")

}
func DBInsertRedis(id string, info string) {
	err := dbredis.Set(id, info, 0).Err()
	if err != nil {
		panic(err)
	}
}
func DBDeleteRedis(id string) {
	err := dbredis.Del(id).Err()
	if err != nil {
		panic(err)
	}
}
func DBGetAllKeysRedis() []string {
	var ReturnData []string
	allkeys, _ := dbredis.Keys("*").Result()
	for _, currentkey := range allkeys {
		//fmt.Println("KEY>> " + currentkey)
		currentvalue, _ := dbredis.Get(currentkey).Result()
		fmt.Println("KEY>> " + currentvalue)
		ReturnData = append(ReturnData, currentvalue)
	}
	return ReturnData
}
