package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

type Requester struct {
	userID    string
	userToken string
	apiClient string
	User      User
}

func NewRequester(config *Config) *Requester {
	return &Requester{
		userID:    config.UserID,
		userToken: config.UserToken,
		apiClient: config.apiClient, // aka User-Agent?
	}
}

func (this *Requester) getGoldAmount() error {
	response, err := this.doTheRequest(http.MethodGet, buildAddress("user?userFields=stats.gp"))
	if err != nil {
		return err
	}
	return checkResponse(response)
}
func (this *Requester) checkArmoire() error {
	//https://habitica.com/api/v3/user/buy-armoire
	response, err := this.doTheRequest(http.MethodPost, buildAddress("user/buy-armoire"))
	if err != nil {
		return err
	}
	return checkResponse(response)
}

func (this *Requester) doTheRequest(method, action string) (*http.Response, error) {
	//log.Println("[DEV] action:", action)
	request, err := http.NewRequest(method, action, nil)
	if err != nil {
		log.Println("[ERROR] http.NewRequest:", err)
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-api-user", this.userID)
	request.Header.Set("x-api-key", this.userToken)
	request.Header.Set("x-client", this.apiClient)
	log.Println("[DEV] header values:", this.userID, this.userToken, this.apiClient)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Printf("[ERROR] client.Do response: %#v\n", response)
		log.Println("[ERROR] client.Do err:", err)
		return response, err
	}
	return response, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func checkResponse(response *http.Response) error {
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("[WARN] could not read response from Habitica:", err)
	}
	if responseBytes != nil {
		var errorResponse ErrorResponse
		json.Unmarshal(responseBytes, &errorResponse)
		log.Printf("[INFO] response: %v\n", errorResponse)
	}
	if !responseIsOk(response.StatusCode) {
		log.Printf("[ERROR] response is not ok: %#v\n", err)
		return errors.New(response.Status)
	}
	return nil
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

func responseIsOk(code int) bool {
	if code == http.StatusOK {
		return true
	}
	if code == http.StatusUnauthorized {
		fmt.Println(unauthorizedText)
	} else {
		fmt.Println("[ERROR] Habitica response code:", code)
	}
	return false
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func buildAddress(action string) string {
	return LiveHost + path.Join(ApiPath, ApiVersion, action)
}

const ApiPath string = "api"
const ApiVersion string = "v3"
const LiveHost string = "https://habitica.com/"

const (
	unauthorizedText = `
* * * * * * * * * * * * * * * * * * * *
Habitica reports that the credentials
  you provided are not authorized to
  access your account.
* * * * * * * * * * * * * * * * * * * *`
)
