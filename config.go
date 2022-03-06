package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// todo: alternatively read from environment variables, if they exist and are populated
// os.Getenv("HABITICA_API_USER")
// os.Getenv("HABITICA_API_KEY")
// os.Getenv("HABITICA_API_CLIENT")

type Config struct {
	UserID    string `json:"user_id"`
	UserToken string `json:"user_token"`
	APIClient string `json:"-"`
}

func NewConfig(apiClient string) *Config {
	config := newConfig(apiClient)
	config.setup()
	return config
}

func newConfig(apiClient string) *Config {
	return &Config{
		APIClient: apiClient,
	}
}

func (this *Config) setup() {
	this.ensureAPIClient()
	if configFileExists() {
		err := this.readConfigFile()
		if err != nil {
			log.Println("[WARN] there was a problem reading from the configuration file:", err)
		} else {
			return
		}
	}

	this.readConfigFromUser()

	if !writeConfigFile(this) {
		fmt.Println("Unable to write the configuration file.")
	}
}

func configFileExists() bool {
	_, err := os.Stat(configFilename)
	if err != nil {
		return false
	}
	return true
}

func (this *Config) readConfigFile() error {
	fileBytes, err := ioutil.ReadFile(configFilename)
	if err != nil {
		log.Println("[ERROR] reading config file:", err)
		return err
	}
	var config Config
	err = json.Unmarshal(fileBytes, &config)
	if err != nil {
		log.Println("[ERROR] json.Unmarshal config file:", err)
		return err
	}
	log.Printf("[DEV] config: %#v\n", config)
	this.UserToken = config.UserToken
	this.UserID = config.UserID
	return nil
}

func (this *Config) readConfigFromUser() {
	var err error
	fmt.Println(configText)
	fmt.Print("Enter your Habitica User ID:")
	reader := bufio.NewReader(os.Stdin)
	this.UserID, err = reader.ReadString('\n')
	if err != nil {
		log.Printf("[ERROR] reading user id: %s; err: %s\n", this.UserID, err)
	}
	this.UserID = strings.TrimSpace(this.UserID)
	fmt.Print("Enter your Habitica API Token:")
	this.UserToken, err = reader.ReadString('\n')
	if err != nil {
		log.Printf("[ERROR] reading api token: %s; err: %s\n", this.UserToken, err)
	}
	this.UserToken = strings.TrimSpace(this.UserToken)
}
func writeConfigFile(config *Config) bool {
	configBytes, err := json.Marshal(config)
	if err != nil {
		log.Println("[ERROR] writeConfigFile json.Marshal:", err)
		return false
	}
	err = ioutil.WriteFile(configFilename, configBytes, 0644)
	if err != nil {
		log.Println("[ERROR] writeConfigFile ioutil.WriteFile:", err)
		return false
	}
	return true
}

func (this *Config) ensureAPIClient() {
	apiClient := os.Getenv("HABITICA_API_CLIENT")
	if this.APIClient == "" && apiClient == "" {
		log.Fatal("[ERROR] API Client field is not set.")
	}
	this.APIClient = apiClient
}

const (
	configFilename string = "config.json"
	configText     string = `
This app requires you to enter your Habitica User ID and API Token.
You can find both of those here:
https://habitica.com/user/settings/api

These credentials will be stored on your computer so that
you do not need to enter them each time you run this application.`
)
