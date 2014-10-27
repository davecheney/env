
an endpoint to load yaml into datastore and memcache based upon appID.

The question then becomes about how to load these values from Datastore and map them.


Must pattern ..


```

import (
    "l"
)

//option 1 - one value returned per call
client.Post(l(ctx, "QUANTUM_VIEW_URL")... )

//option 2 - return a configuration struct which has every value mapped for the application

```



```
    

```





A configuration layer for variables which vary according to environment, except/ using GAE's `[appengine.AppID(ctx)](https://cloud.google.com/appengine/docs/go/reference#AppID)` to support the generic idea of server environments (production, staging, etc.).


import(
    env "gae_env"
)


func handler(c appengine.Context, w http.ResponseWriter, r *http.Request) {
    
}

func init() {
    env.load()
}

in init use env.load() ...

then, when you



synonyms - so generic environments can be used



GAE App


appid-env

throw away the language of "production", "staging". This is application aware env vars.

map, err := env.scan(ctx)



```yaml
rockfish-project:
rockfish-staging:

```

```golang

import(
    "env"
    "appengine"
)

environment, err := env.Is("production") 


```




import "env"



need a way to declare environments dynamically - e.g. staging2, integration, etc..


url, err := env(c, "QUANTUM_URL")

internally env will 



env.IsProduction
env.IsStaging
env.IsTest
env.IsDevelopment

the above call appengine.


Write the yaml with environment vars:

production:
    - QUANTUM_VIEW_URL: "http://prod"
staging:
    - QUANTUM_VIEW_URL: "http://staging"

  