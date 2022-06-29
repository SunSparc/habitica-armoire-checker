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
		apiClient: config.apiClient,
	}
}

func (this *Requester) getGoldAmount() error {
	response, err := this.doTheRequest(http.MethodGet, buildAddress("user?userFields=stats.gp"))
	if err != nil {
		return err
	}
	return this.processResponse(response)
}
func (this *Requester) checkArmoire() error {
	//https://habitica.com/api/v3/user/buy-armoire
	response, err := this.doTheRequest(http.MethodPost, buildAddress("user/buy-armoire"))
	if err != nil {
		return err
	}
	return this.processResponse(response)
}

func (this *Requester) doTheRequest(method, action string) (*http.Response, error) {
	request, err := http.NewRequest(method, action, nil)
	if err != nil {
		log.Println("[ERROR] http.NewRequest:", err)
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-api-user", this.userID)
	request.Header.Set("x-api-key", this.userToken)
	request.Header.Set("x-client", this.apiClient)

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

func (this *Requester) processResponse(response *http.Response) error {
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("[WARN] could not read response from Habitica:", err)
	}

	if responseBytes != nil {
		var habiticaResponse HabiticaResponse
		err := json.Unmarshal(responseBytes, &habiticaResponse)
		if err != nil {
			log.Printf("[ERROR] unmarshal response: %v\n", habiticaResponse)
		}

		if response.StatusCode == http.StatusOK {
			this.User.Data = habiticaResponse.Data
			return nil
		}
		fmt.Printf("Habitica says:\n%q\n", habiticaResponse.Message)
		fmt.Printf("\n\n(Press Enter/Return to continue.)\n")
		readFromStdin()
		return errors.New(response.Status)
	}
	return nil
}

type HabiticaResponse struct {
	Success bool     `json:"success"`
	Data    UserData `json:"data"`
	Error   string   `json:"error"`
	Message string   `json:"message"`
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func buildAddress(action string) string {
	return LiveHost + path.Join(ApiPath, ApiVersion, action)
}

const ApiPath string = "api"
const ApiVersion string = "v3"
const LiveHost string = "https://habitica.com/"
