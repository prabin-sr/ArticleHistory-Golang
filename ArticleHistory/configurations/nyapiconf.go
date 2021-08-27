package configurations

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// NYTimesAPIConfig - Struct for NY Times API Configurations
type NYTimesAPIConfig struct {
	Credentials  map[string]string `json:"credentials"`
	APIEndpoints map[string]string `json:"api_endpoints"`
}

// GetNYTimesAPIConfig - Retrieves the email data
func GetNYTimesAPIConfig() NYTimesAPIConfig {
	jsonFile, err := os.Open("./configurations/static/nyapiconf.json")

	nyapiconfiguration := NYTimesAPIConfig{}

	// if os.Open returns an error then handle it
	if err != nil {
		return nyapiconfiguration
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	defer jsonFile.Close()

	json.Unmarshal(byteValue, &nyapiconfiguration)

	return nyapiconfiguration
}
