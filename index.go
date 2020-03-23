package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-redis/redis"
)

func main() {
	ageOfEmpiresGame()
}

func ageOfEmpiresGame() {
	var client *redis.Client
	var result string

	name := getUnitName()

	client = redisInitialize()

	result, err := client.Get(name).Result()
	Logger(err)

	if len(result) == 0 {
		result = getUnitByNameOrId(name)
	}

	err = client.Set(name, result, 0).Err() // cache timeout is zero
	Logger(err)

	log.Println("Unit Details are below: ")
	log.Println(result)
}

func getUnitName() string {
	var name string

	fmt.Println("Please enter the Unit name")
	_, err := fmt.Scan(&name)

	Logger(err)

	return name
}

func redisInitialize() (client *redis.Client) {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // make sure that your Redis-Server is running
		Password: "",               // no password set
		DB:       0,                // use default DB
	})
}

func getUnitByNameOrId(name string) string {

	resp, err := http.Get(fmt.Sprintf("https://age-of-empires-2-api.herokuapp.com/api/v1/unit/%v", name))

	Logger(err)

	body, err := ioutil.ReadAll(resp.Body)

	Logger(err)

	defer resp.Body.Close()

	return string(body)
}

func Logger(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}
