package env

import (
	"appengine"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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

func Load(relativePathToFile string) (err error) {
	return Env.Load(relativePathToFile)
}

func Get(context appengine.Context, field string) (interface{}, error) {
	return Env.Get(context, field)
}

func Name(context appengine.Context) string {
	return Env.Name(context)
}

/*
	Loads the config file from relativePathToFile and stores it in Env.
	Returns an error if this wasn't possible.
*/
func (e *env) Load(relativePathToFile string) (err error) {
	e.emap = map[string]interface{}{}
	e.raw = []byte{}

	data, err := ioutil.ReadFile(relativePathToFile)
	if err != nil {
		return err
	}
	e.raw = data
	err = json.Unmarshal(e.raw, &e.emap)
	if err != nil {
		return err
	}

	//Check for mappings, error if they're not present
	_, ok := e.emap[MAPPING_KEY_NAME]
	if !ok {
		errorMsg := fmt.Sprintf("JSON Malformed. Missing top level property named '%v' to determine projectId", MAPPING_KEY_NAME)
		return errors.New(errorMsg)
	}

	return nil
}

/*
	Retrieves an environment variable from the loaded json
	Returns an error if this wasn't possible.
*/
func (e *env) Get(context appengine.Context, field string) (interface{}, error) {
	currentEnvName, err := e.getCurrentEnvName(context)
	if err != nil {
		return "", err
	}
	retrievedKey, err := e.getKeyFromEnv(context, currentEnvName, field)
	if (err != nil && currentEnvName != DEFAULT_MAP_NAME) {
		//Try the default map if we can't find it in the current map
		retrievedKey, err = e.getKeyFromEnv(context, DEFAULT_MAP_NAME, field)
		if err != nil {
			return "", err
		}
		return retrievedKey, nil
	} else if (err != nil) {
		return "", err
	}

	return retrievedKey, nil
}

func (e *env) getKeyFromEnv(context appengine.Context, envName string, field string) (interface{}, error) {
	currentEnvData, err := e.getEnvData(context, envName)
	if err != nil {
		return "", err
	}
	currentEnvDataMap := currentEnvData.(map[string]interface{})

	retrievedKey, ok := currentEnvDataMap[field]

	if !ok {
		currentAppId := appengine.AppID(context)
		currentEnvName, err := e.getCurrentEnvName(context)
		if err != nil {
			return "", err
		}
		errorMsg := fmt.Sprintf("Cannot retrieve environment variable. Missing field %v for environment %v with App Id %v", field, currentEnvName, currentAppId)
		return "", errors.New(errorMsg)
	}
	return retrievedKey, nil
}

/*
	Determines the currently running environment name, ie "Production"
	If an unknown app id is used, it attempts to get the default environment.
	NOTE: If the default environment does not exist, the error is suppressed.
*/
func (e *env) Name(context appengine.Context) string {
	currentEnvName, _ := e.getCurrentEnvName(context)
	return currentEnvName
}

func (e *env) getEnvData(context appengine.Context, currentEnvName string) (interface{}, error) {
	currentEnvInterface, ok := e.emap[currentEnvName]
	if !ok {
		currentAppId := appengine.AppID(context)
		errorMsg := fmt.Sprintf("JSON Malformed. Missing top level property named %v associated with current App Id, %v", currentEnvName, currentAppId)
		return map[string]interface{}{}, errors.New(errorMsg)
	}
	return currentEnvInterface, nil
}

func (e *env) getCurrentEnvName(context appengine.Context) (string, error) {
	currentAppId := appengine.AppID(context)

	envInterface, ok := e.emap[MAPPING_KEY_NAME]
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
