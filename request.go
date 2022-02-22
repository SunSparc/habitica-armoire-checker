package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

// todo: make one request method that both of these requests use instead

func getGoldAmount(action string, config *Config) (float64, error) {
	client := http.Client{}
	request, err := http.NewRequest(http.MethodGet, BuildAddress("v3", action), nil)
	if err != nil {
		log.Println("[ERROR] http.NewRequest:", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-api-user", config.UserID)
	request.Header.Set("x-api-key", config.UserToken)
	request.Header.Set("x-client", config.APIClient)

	response, err := client.Do(request)
	if err != nil {
		log.Println("[ERROR] client.Do:", err)
		return 0, err
	}
	if response.StatusCode != http.StatusOK {
		log.Println("got a bad response:", response.StatusCode, response.Status)
		return 0, err
	}
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("[ERROR] ioutil.ReadAll:", err)
		return 0, err
	}

	var user User
	err = json.Unmarshal(responseBytes, &user)
	if err != nil {
		log.Println("[ERROR] json.Unmarshal:", err)
		return 0, err
	}

	return user.Data.Stats.Gold, nil
}

func doArmoireRequest(config *Config) User {
	//https://habitica.com/api/v3/user/buy-armoire
	client := http.Client{}
	request, err := http.NewRequest(http.MethodPost, BuildAddress("v3", "user/buy-armoire"), nil)
	if err != nil {
		log.Println("[ERROR] http.NewRequest:", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-api-user", config.UserID)
	request.Header.Set("x-api-key", config.UserToken)
	request.Header.Set("x-client", config.APIClient)

	response, err := client.Do(request)
	if err != nil {
		log.Println("[ERROR] client.Do:", err)
		//return response, err
	}
	if response.StatusCode != http.StatusOK {
		log.Println("got a bad response:", response.StatusCode, response.Status)
		//return response, err
	}
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("[ERROR] ioutil.ReadAll:", err)
		//return nil, err
	}
	//log.Printf("[doArmoireRequest][response]: %#v\n", string(responseBytes))

	var user User
	err = json.Unmarshal(responseBytes, &user)
	if err != nil {
		log.Println("json.Unmarshal error:", err)
	}
	user.StatusCode = response.StatusCode

	return user
}

func BuildAddress(version, action string) string {
	return LIVE_HOST + path.Join(API_PATH, version, action)
}

const API_PATH string = "api"
const LIVE_HOST string = "https://habitica.com/"
