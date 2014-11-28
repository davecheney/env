/*
	Google App Engine Environment Variable Library.

	A simpler way of thinking about programming environments when working with `appengine.Context' AppIDs (or Project IDs) on Google's Cloud Platform, GAE
	
	See: https://github.com/rockpoollabs/env for README and examples.
*/
package env

import (
	"appengine"
	"encoding/json"
	"errors"
	"io/ioutil"
	"fmt"
)

type env struct {
	emap map[string]interface{}
	raw  []byte
}

const MAPPING_KEY_NAME = "mappings"
const DEFAULT_MAP_NAME = "default"

var Env env

func init() {
	Env = env{emap: map[string]interface{}{}}
}

/*
	Loads the config file from relativePathToFile and stores it in Env.
	Returns an error if this wasn't possible.
*/
func Load(relativePathToFile string) (err error) {
	Env.emap = map[string]interface{}{}
	Env.raw = []byte{}

	data, err := ioutil.ReadFile(relativePathToFile)
	if err != nil {
		return err
	}

	Env.raw = data
	err = json.Unmarshal(Env.raw, &Env.emap)
	if err != nil {
		return err
	}

	//Check for mappings, error if they're not present
	_, ok := Env.emap[MAPPING_KEY_NAME]
	if !ok {
		errorMsg := fmt.Sprintf("JSON Malformed. Missing top level property named '%v' to determine projectId", MAPPING_KEY_NAME)
		return errors.New(errorMsg)
	}

	return nil
}

/*
	Must with Load.
*/
func MustLoad(relativePathToFile string) {
	err := Load(relativePathToFile)
	if err != nil {
		panic("Unable to load " + relativePathToFile + ": " + err.Error())
	}
}

/*
	Retrieves an environment variable from the loaded json
	Returns a boolean if this wasn't possible.
*/
func GetOk(c appengine.Context, field string) (interface{}, bool) {
	currentEnvName, err := getCurrentEnvName(c)
	if err != nil {
		return "", false
	}
	val, err := getKeyFromEnv(c, currentEnvName, field)
	if (err != nil && currentEnvName != DEFAULT_MAP_NAME) {
		//Try the default map if we can't find it in the current map
		val, err = getKeyFromEnv(c, DEFAULT_MAP_NAME, field)
		if err != nil {
			return "", false
		}

		return val, true

	} else if (err != nil) {
		return "", false
	}

	return val, true
}

/*
	Proxies a call to `GetOk`, suppressing any error (the client will have to deal with this).
 */

func Get(c appengine.Context, key string) (interface{}) {
	val, _ := GetOk(c, key)
	return val
}


/*
	Determines the currently running environment name, ie "Production"
	If an unknown app id is used, it attempts to get the default environment.
	NOTE: If the default environment does not exist, the error is suppressed.
*/

func Name(c appengine.Context) string {
	currentEnvName, _ := getCurrentEnvName(c)
	return currentEnvName
}

/*
	Boolean check for which environment this is in. For example: Is(c, "production")
*/
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
		errorMsg := fmt.Sprintf("Cannot retrieve environment variable. Missing field %v for environment %v with App Id %v", field, currentEnvName, currentAppId)
		return "", errors.New(errorMsg)
	}
	return retrievedKey, nil
}

func getEnvData(c appengine.Context, currentEnvName string) (interface{}, error) {
	currentEnvInterface, ok := Env.emap[currentEnvName]
	if !ok {
		currentAppId := appengine.AppID(c)
		errorMsg := fmt.Sprintf("JSON Malformed. Missing top level property named %v associated with current App Id, %v", currentEnvName, currentAppId)
		return map[string]interface{}{}, errors.New(errorMsg)
	}
	return currentEnvInterface, nil
}

func getCurrentEnvName(c appengine.Context) (string, error) {
	currentAppId := appengine.AppID(c)

	envInterface, ok := Env.emap[MAPPING_KEY_NAME]

	if !ok {
		errorMsg := fmt.Sprintf("JSON Malformed. Missing top level property named '%v' to determine projectId", MAPPING_KEY_NAME)
		return "", errors.New(errorMsg)
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
