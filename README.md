#GAE Environment Variable Library#

Available under the MIT (Expat) License - see bottom of README.

## Installation

You will need to be setup with google app engine and `goapp`. See the [google app engine documentation for details](https://cloud.google.com/appengine/docs/go/gettingstarted/introduction)

Dependencies are installed via:
```bash
make deps
```

## Example and Usage

```golang
package main

import (
	"github.com/rockpoollabs/env"
	"os"
	"log"
)

func main() {
	//Set the env file
    env := Env
    err := env.Load("./environment.json")

    //Get current environment
    name := env.Name(ctx) 
    log.Printf("Environment: %v",msg) //prints "Environment: testing" in testing

	//Get environment variable
	msg, err := env.Get(ctx, "Message")
	log.Printf("Message: %v",msg) // prints "Message: TestingMsg" in testing
	
	///...
}
```

The corresponding environment.json file at the top level of the  would be:
```json
{
    "mappings" : {
        "production" : "App-Prod",
        "staging" : "App-Stage",
        "testing" : "App-Test"
    },
	"default" : {
		"Message" : "DefaultMsg"
	},
	"production" : {
		"Message" : "ProductionMsg"
	},
	"staging" : {
		"Message" : "StagingMsg"
	},
	"testing" : {
		"Message" : "TestingMsg"
	}
}
```

The tests for this project use the sample-environment.json file. You can look at that for another example.

## Godoc
[http://godoc.org/github.com/rockpoollabs/environment](http://godoc.org/github.com/rockpoollabs/environment)
TODO: Get correct address when it's live

## Testing

```bash
make test
make live-test //then browse to http://localhost:8080/
```

## Formatting

```bash
make fmt
```

## Development

Ensure that you are running GO 1.2.1

## Contributing

You're welcome to make a Pull Request; please include tests for anything you want to contribute.

## MIT License

The MIT License (MIT)

Copyright (c) 2014 Rockpool Labs

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.