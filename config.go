package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

// todo: alternatively read from environment variables, if they exist and are populated
// os.Getenv("HABITICA_API_USER")
// os.Getenv("HABITICA_API_KEY")
// os.Getenv("HABITICA_API_CLIENT")

type Config struct {
	UserID         string `json:"user_id"`
	UserToken      string `json:"user_token"`
	apiClient      string // `json:"-"`
	configPath     string
	resetRequested bool
}

func NewConfig(apiClient string, reset bool) *Config {
	config := newConfig(apiClient)
	config.resetRequested = reset
	config.setup()
	return config
}

func newConfig(apiClient string) *Config {
	return &Config{
		apiClient: apiClient,
	}
}

func (this *Config) setup() {
	this.ensureAPIClient()
	this.ensureConfigPath()
	this.getConfigValues()

	if !writeConfigFile(this) { // todo: do not write the config file every time, only when we have new config data
		fmt.Println("Unable to write the configuration file.")
	}
}

func (this *Config) getConfigValues() {
	// if file does not exist -> get configs from user
	// if file does exist & user wants new configs -> get configs from user
	// if file does exist & user wants old configs -> get configs from file
	if this.configFileExists() && !this.resetRequested {
		fmt.Printf("Welcome back!\n\nShall we use the same credentials that\n  we used last time?\n\n")
		fmt.Print("Yes, No, Cancel (y,n,c): ")
		keepOrNew := readFromStdin()
		if strings.ContainsRune(keepOrNew, 'c') || strings.Contains(keepOrNew, "cancel") {
			os.Exit(1)
		}
		if strings.ContainsRune(keepOrNew, 'n') || strings.Contains(keepOrNew, "no") {
			fmt.Println("[DEV] getting new config values:", keepOrNew) // DEV
			this.readConfigFromUser()
			return
		}
		fmt.Println("[DEV] getting configuration from file:", keepOrNew) // DEV

		err := this.readConfigFile()
		if err == nil {
			return
		}
		log.Println("[WARN] there was a problem reading from the configuration file:", err)
	}
	log.Println("no config file or reset requested")
	this.readConfigFromUser()
}

func (this *Config) configFileExists() bool {
	_, err := os.Stat(this.configPath)
	if err != nil {
		//log.Println("configFileExists, false:", err)
		return false
	}
	//log.Println("configFileExists, true:", err)
	return true
}

func (this *Config) readConfigFile() error {
	fileBytes, err := ioutil.ReadFile(this.configPath)
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
	//log.Printf("[DEV] config: %#v\n", config) // DEV
	this.UserToken = config.UserToken
	this.UserID = config.UserID
	return nil
}

func (this *Config) readConfigFromUser() {
	clearScreen()
	fmt.Println(configText)
	fmt.Print("Enter your Habitica User ID: ")
	this.UserID = readFromStdin()
	fmt.Print("Enter your Habitica API Token: ")
	this.UserToken = readFromStdin()

	fmt.Println()
}

func readFromStdin() string {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("[ERROR] reading user input: %s\n", err)
	}
	//fmt.Println("this was the input:", input)
	return strings.ToLower(strings.TrimSpace(input))
}

func writeConfigFile(config *Config) bool {
	configBytes, err := json.Marshal(config)
	if err != nil {
		log.Println("[ERROR] writeConfigFile json.Marshal:", err)
		return false
	}
	err = ioutil.WriteFile(config.configPath, configBytes, 0644)
	if err != nil {
		log.Println("[ERROR] writeConfigFile ioutil.WriteFile:", err)
		return false
	}
	fmt.Printf("")
	return true
}

func (this *Config) ensureAPIClient() {
	apiClient := os.Getenv("HABITICA_API_CLIENT")
	if this.apiClient == "" && apiClient == "" {
		log.Fatal("[ERROR] API Client field is not set.")
	}
	this.apiClient = apiClient
}

func (this *Config) ensureConfigPath() {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal("[ERROR] could not retrieve config directory:", err)
	}
	configDir := path.Join(userConfigDir, configDirectory)
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		log.Fatal("[ERROR] could not make config directory:", err)
	}
	this.configPath = path.Join(configDir, configFilename)
}

const (
	configDirectory string = "HabiticaArmoireChecker"
	configFilename  string = "config.json"
	configText      string = `
This application requires you to enter
your Habitica User ID and API Token.

You can find both of those here:
  https://habitica.com/user/settings/api

These credentials will be stored on
  your computer so that you do not need
  to enter them each time you run this
  application.
`
)
