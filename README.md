#GAE Environment Variable Library#

Available under the MIT (Expat) License - see bottom of README.

## Rationale

### Mapping GAE Projects to Programming Environments

A simpler way of thinking about programming environments when working with `appengine.Context' AppIDs (or Project IDs) on [Google's Cloud Platform, GAE](https://cloud.google.com/appengine/docs)

So, instead of

```
func SomeHandler(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	if appengine.AppID(c) == "rockpool-production" {
		//do this
	} else if appengine.AppID(c) == "rockpool-staging" {
		//do that
	}
```
you can do this

```
func SomeHandler(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	if env.Is("production") {
		//do this
	} else env.Is("staging") {
		//do that
	}

```
As well we being cleaner, it means you can meaningfully refer to environments in your code without having to hardcode the name of the Google Cloud Project.

### Environment Specific Variables

You may then access environment specific variables with `env.Get(c, appengine.Context, key string)` or `env.GetOk(c, appengine.Context, key string)`:

```
func SomeHandler(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	db_password := env.Get(c, "db_password").(string)
	//or, use the GetOk pattern
	db_password2, ok := env.GetOk(c, "db_password")
}
```

### Example JSON Config

```
{
    "mappings" : {
        "production" : "rockpool-production",
        "staging" : "rockpool-staging",
        "test" : "rockpool-test"
    },
	"default" : {
		"question" : "How about a nice game of chess?",
	},
	"production" : {
		"message" : "I am a production Msg",
		"tolerance" : 0.12,
		"acceptableRank" : ["2","3","4"]
	},
	"staging" : {
		"message" : "I am a staging Msg",
		"tolerance" : 0.13,
		"acceptableRank" : ["5","6","7"]
		"question" : "How about a nice game of chess (in staging)?",
	}
```


### The default stanza

Default variables mean that every environment will gain this value, unless they choose you choose to override it within the environment's stanza.

In the example configuration above, production will have a question with `"How about a nice game of chess?"` as its value but staging will have `"How about a nice game of chess (in staging)?"`.

### How to load the JSON config

Typically setup is down within a init function called early in the lifecycle of one of your packages.

Use `MustLoad` to create a runtime panic if the JSON config cannot be loaded.
Use `Load` to return an error to handle yourself if the JSON config cannot be loaded.

```
func init() {
	env.MustLoad("./environment.json")
	// or ...
	err := env.Load("./environment.json")

	if err != nil {
		...
	}
}
```

## Installation

`goapp get github.com/rockpoollabs/env`

Also see Configuration below.

## Godoc
[http://godoc.org/github.com/rockpoollabs/environment](http://godoc.org/github.com/rockpoollabs/environment)
TODO: Get correct address when it's live

## Configuration

You will need to be setup with google app engine and `goapp`. See the [google app engine documentation for details](https://cloud.google.com/appengine/docs/go/gettingstarted/introduction)

Dependencies are installed via:
```bash
make deps
```

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

Ensure that you are running Goapp. This was developed on goapp 1.9.15 but should work with earlier versions supporting `appengine.AppId`.

```
	goapp version
	go version go1.2.1 (appengine-1.9.15) darwin/amd64
```

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