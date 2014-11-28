/*
	Google App Engine Environment Variable Library.

	A simpler way of thinking about programming environments when working with `appengine.Context' AppIDs (or Project IDs) on Google's Cloud Platform, GAE

	See: https://github.com/rockpoollabs/env for README and examples.
*/
package env

import (
	"encoding/json"
	"fmt"
	"os"

	"appengine"
)

const MAPPING_KEY_NAME = "mappings"
const DEFAULT_MAP_NAME = "default"

var Env map[string]interface{}

// Loads the config file from relativePathToFile and stores it in Env.
// Returns an error if this wasn't possible.
func Load(relativePathToFile string) error {
	r, err := os.Open(relativePathToFile)
	if err != nil {
		return err
	}

	d := json.NewDecoder(r)
	Env = make(map[string]interface{})
	err = d.Decode(&Env)
	if err != nil {
		return err
	}

	//Check for mappings, error if they're not present
	_, ok := Env[MAPPING_KEY_NAME]
	if !ok {
		return fmt.Errorf("JSON Malformed. Missing top level property named %q to determine projectId", MAPPING_KEY_NAME)
	}

	return nil
}

// MustLoad loads the config file from relativePathToFile and panics if an error occurs.
func MustLoad(relativePathToFile string) {
	err := Load(relativePathToFile)
	if err != nil {
		panic("Unable to load " + relativePathToFile + ": " + err.Error())
	}
}

// GetOk retrieves an environment variable from the loaded json.
// Returns a boolean if this wasn't possible.
func GetOk(c appengine.Context, field string) (interface{}, bool) {
	currentEnvName, err := getCurrentEnvName(c)
	if err != nil {
		return "", false
	}
	val, err := getKeyFromEnv(c, currentEnvName, field)
	if err != nil && currentEnvName != DEFAULT_MAP_NAME {
		//Try the default map if we can't find it in the current map
		val, err = getKeyFromEnv(c, DEFAULT_MAP_NAME, field)
		if err != nil {
			return "", false
		}
		return val, true
	}

	return val, err == nil
}

// Get proxies a call to GetOk, suppressing any error (the client will have to deal with this).
func Get(c appengine.Context, key string) interface{} {
	val, _ := GetOk(c, key)
	return val
}

// Name fetermines the currently running environment name, ie "Production"
// If an unknown app id is used, it attempts to get the default environment.
// NOTE: If the default environment does not exist, the error is suppressed.
func Name(c appengine.Context) string {
	currentEnvName, _ := getCurrentEnvName(c)
	return currentEnvName
}

// Is checks which environment this is in. For example: Is(c, "production")
func Is(c appengine.Context, envName string) bool {
	return Name(c) == envName
}

func getKeyFromEnv(c appengine.Context, envName string, field string) (interface{}, error) {
	currentEnvData, err := getEnvData(c, envName)
	if err != nil {
		return "", err
	}
	currentEnvDataMap := currentEnvData.(map[string]interface{})

	retrievedKey, ok := currentEnvDataMap[field]

	if !ok {
		currentAppId := appengine.AppID(c)
		currentEnvName, err := getCurrentEnvName(c)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("Cannot retrieve environment variable. Missing field %q for environment %q with App Id %q", field, currentEnvName, currentAppId)
	}
	return retrievedKey, nil
}

func getEnvData(c appengine.Context, currentEnvName string) (interface{}, error) {
	currentEnvInterface, ok := Env[currentEnvName]
	if !ok {
		currentAppId := appengine.AppID(c)
		return nil, fmt.Errorf("JSON Malformed. Missing top level property named %q associated with current App Id %q", currentEnvName, currentAppId)
	}
	return currentEnvInterface, nil
}

func getCurrentEnvName(c appengine.Context) (string, error) {
	currentAppId := appengine.AppID(c)

	envInterface, ok := Env[MAPPING_KEY_NAME]

	if !ok {
		return "", fmt.Errorf("JSON Malformed. Missing top level property named %q to determine projectId", MAPPING_KEY_NAME)
	}

	envs := envInterface.(map[string]interface{})
	currentEnvName := DEFAULT_MAP_NAME
	for envName, envProjectIdInterface := range envs {
		envProjectId := envProjectIdInterface.(string)
		if envProjectId == currentAppId {
			currentEnvName = envName
		}
	}
	return currentEnvName, nil
}
