package env

import(
	"fmt"
)

type env struct {
	emap map[string]interface{}
}

/*
	Loads the config file from relativePathToFile and stores it in Env.
	Returns an error if this wasn't possible.
 */
func (e env) Load(relativePathToFile string) (err error) {
	return
}

func (e env) Get(context interface{}, field string) interface{} {
	return interface{}
}

func (e env) Name(context interface{}) string {
	return "" //would return an environment name like "staging" by looking at appengine.AppID(context)
}


var Env env

func init() {
	Env = env{ emap: map[string]interface{}{} }
}


