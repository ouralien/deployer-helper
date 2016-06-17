package deployerHelper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var config *Config

type deployerPayload struct {
	GitURL           string `json:"gitUrl"`
	DeployKey        string `json:"deployKey"`
	Registry         string `json:"registry"`
	RegistryLogin    string `json:"registryLogin"`
	RegistryPassword string `json:"registryPassword"`
	RegistryEmail    string `json:"registryEmail"`
	ImageName        string `json:"imageName"`
	ImageVersion     string `json:"imageVersion"`
	WebhookToken     string `json:"webhookToken"`
	ExtraVars        string `json:"extraVars"`
}

type deployerResponse struct {
	Output string `json:"output"`
}

// Config struct
type Config struct {
	deployKey string
	service   string
	image     string
	host      string
	token     string
	repo      string
	registry  string
	login     string
	password  string
	email     string
	extraVars string
}

// Init is a func
func Init() {
	config = &Config{}
	config.image = "master"
}

// Deploy is a func
func Deploy() {
	fmt.Println("Host:", config.service, config.host)

	payload := deployerPayload{
		GitURL:           config.repo,
		DeployKey:        config.deployKey,
		Registry:         config.registry,
		RegistryLogin:    config.login,
		RegistryPassword: config.password,
		RegistryEmail:    config.email,
		ImageName:        config.service,
		ImageVersion:     config.image,
		WebhookToken:     config.token,
		ExtraVars:        config.extraVars,
	}

	r := callService(payload)
	responseHandler(r)
}

func callService(payload deployerPayload) *http.Response {
	serviceURL := fmt.Sprintf("https://%s/deploy", config.host)
	fmt.Printf("Calling Service : %v \n", serviceURL)
	fmt.Printf("With payload : %v \n", payload)

	var URL *url.URL
	URL, err := url.Parse(serviceURL)
	if err != nil {
		fmt.Printf("Error Parsing URL : %v \n", serviceURL)
		os.Exit(-1)
	}

	bArray, err := json.Marshal(payload)
	req, err := http.NewRequest("POST", URL.String(), bytes.NewBuffer(bArray))
	if err != nil {
		fmt.Printf("Error Creating Request : %v \n", err.Error())
		os.Exit(-1)
	}

	req.Header.Set("content-type", "application/json")
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error Calling Client : %v \n", err.Error())
		os.Exit(-1)
	}

	return res
}

func responseHandler(r *http.Response) {
	res := &deployerResponse{}

	if r.StatusCode < 200 || r.StatusCode > 400 {
		fmt.Printf("Status %v: NOT OK\n", r.StatusCode)
		os.Exit(-1)
	}

	if err := json.NewDecoder(r.Body).Decode(res); err != nil {
		fmt.Printf("UNABLE to parse JSON response from service\n")
		os.Exit(-1)
	}

	fmt.Printf("\nResponse From Service : \n\n%v", res.Output)

	if strings.Contains(res.Output, "Error:") {
		os.Exit(-1)
	}
}

// SetDeployKey is a func
func SetDeployKey(value string) {
	config.deployKey = value
}

// SetService is a func
func SetService(value string) {
	config.service = value
}

// SetImage is a func
func SetImage(value string) {
	config.image = value
}

// SetHost is a func
func SetHost(value string) {
	config.host = value
}

// SetToken is a func
func SetToken(value string) {
	config.token = value
}

// SetRepo is a func
func SetRepo(value string) {
	config.repo = value
}

// SetExtraVars is a func
func SetExtraVars(value string) {
	config.extraVars = value
}

// SetRegistry is a func
func SetRegistry(host string, login string, password string, email string) {
	config.registry = host
	config.login = login
	config.password = password
	config.email = email
}
