package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	jsonFile, err := os.Open("./api_key.json")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer jsonFile.Close()

	byteValue, readAllErr := ioutil.ReadAll(jsonFile)
	if readAllErr != nil {
		log.Fatalln("Error while parsing ./api_key.json")
		return
	}

	var apiKey ApiKeyData

	unmarshalErr := json.Unmarshal(byteValue, &apiKey)
	if unmarshalErr != nil {
		log.Fatalln("Error while unmarshalling ./api_key.json")
		return
	}

	challenge := NewChallenge(apiKey)
	challenge.Run()
}
