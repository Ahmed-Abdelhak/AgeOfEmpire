package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-redis/redis"
)

var client *redis.Client
var result string

func main() {

	// disclaimer: I don't like writing comments in code, because it supposed to be readable
	// but It's my first time to write code in Go, that's why I keep commenting as a logger 													for what I do

	// reading input from user
	var name string = getUnitName()

	// Configurations of Redis
	client = redisConfig()

	// get from Cache
	result, err := client.Get(name).Result()
	check(err)

	// HTTP Call to endpoint if cache is empty
	if len(result) == 0 {
		result = getUnitByNameOrId(name)
	}

	// saving to Cache
	err = client.Set(name, result, 0).Err() // cache timeout is zero
	check(err)

	// logging the output to the user
	log.Println("Unit Details are below: ")
	log.Println(result)
}

func getUnitName() string {
	var name string

	fmt.Println("Please enter the Unit name")
	_, err := fmt.Scan(&name)

	check(err)

	return name
}

func redisConfig() (client *redis.Client) {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func getUnitByNameOrId(name string) string {

	resp, err := http.Get(fmt.Sprintf("https://age-of-empires-2-api.herokuapp.com/api/v1/unit/%v", name))

	check(err)

	body, err := ioutil.ReadAll(resp.Body)

	check(err)

	defer resp.Body.Close()

	return string(body)
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}
