package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

type Requester struct {
	config *Config
	User   User
}

func NewRequester() *Requester {
	return &Requester{
		config: NewConfig(APIClient),
	}
}

func (this *Requester) getGoldAmount() error {
	err := this.doTheRequest(http.MethodGet, "user?userFields=stats.gp")
	if err != nil {
		return err
	}
	return nil
}
func (this *Requester) checkArmoire() error {
	//https://habitica.com/api/v3/user/buy-armoire
	err := this.doTheRequest(http.MethodPost, "user/buy-armoire")
	if err != nil {
		return err
	}
	return nil
}

func (this *Requester) doTheRequest(method, action string) error {
	client := http.Client{}
	request, err := http.NewRequest(method, BuildAddress(ApiVersion, action), nil)
	if err != nil {
		log.Println("[ERROR] http.NewRequest:", err)
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-api-user", this.config.UserID)
	request.Header.Set("x-api-key", this.config.UserToken)
	request.Header.Set("x-client", this.config.APIClient)

	response, err := client.Do(request)
	if err != nil {
		log.Println("[ERROR] client.Do:", err)
		return err
	}
	if response.StatusCode != http.StatusOK {
		log.Println("got a bad response:", response.StatusCode, response.Status)
		return err
	}
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("[ERROR] ioutil.ReadAll:", err)
		return err
	}

	err = json.Unmarshal(responseBytes, &this.User)
	if err != nil {
		log.Println("[ERROR] json.Unmarshal:", err)
		return err
	}
	return nil
}

func BuildAddress(version, action string) string {
	return LiveHost + path.Join(ApiPath, version, action)
}

const ApiPath string = "api"
const ApiVersion string = "v3"
const LiveHost string = "https://habitica.com/"
