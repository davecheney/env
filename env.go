package env

import(
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)


func Load(e string) (conf map[string]interface{}, err error) {
	b, err = ioutil.ReadFile("./env.json")
	var f interface{}
	err := json.Unmarshal(b, &f)
	conf := f.(map[string]interface{})
	return
}

//func Handler1(w http.ResponseWriter, r *http.Request) {
//	c := appengine.NewContext(r)
//
//	c.Infof("Starting 1 ... ")
//	err := os.Setenv("arnie", "I'll be back!")
//	if err != nil {
//		c.Errorf("Could not set env var")
//	}
//
//	b, err := ioutil.ReadFile("./env.json")
//
//	c.Infof(">>>> ", string(b))
//
//	w.Write([]byte("hello 1"))
//}
//
//func Handler2(w http.ResponseWriter, r *http.Request) {
//	c := appengine.NewContext(r)
//	c.Infof("Starting 2... ")
//	c.Infof(">> arnie? :", os.Getenv("arnie"))
//
//	w.Write([]byte("hello 2"))
//}


//func init() {
//	http.HandleFunc("/hello1", Handler1)
//	http.HandleFunc("/hello2", Handler2)
//}
