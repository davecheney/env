package env

import (
	"testing"

	c "github.com/smartystreets/goconvey/convey"

	"appengine_internal"
)

const TESTING_ENV_NAME = "App-Test"

func TestLoad(t *testing.T) {

	c.Convey("With a json file, it will load it into memory", t, func() {
		fileLocation := "sample-environment.json"
		err := Load(fileLocation)

		c.So(err, c.ShouldBeNil)
		c.So(Env, c.ShouldNotBeNil)
	})

	c.Convey("With a json file missing mappings it errors", t, func() {
		fileLocation := "mappings-missing.json"
		err := Load(fileLocation)
		c.So(err, c.ShouldNotBeNil)
	})
}

func TestGet(t *testing.T) {
	ctx := ContextMock{}
	ctx.SetAppId(TESTING_ENV_NAME)

	c.Convey("Before loading the json file", t, func() {
		_, ok := GetOk(ctx, "message")
		c.So(ok, c.ShouldBeFalse)
	})

	c.Convey("With a complete loaded json file", t, func() {
		fileLocation := "./sample-environment.json"
		err := Load(fileLocation)
		c.So(err, c.ShouldBeNil)

		c.Convey("It retrieves the correct vars", func() {
			msg, ok := GetOk(ctx, "message")
			c.So(ok, c.ShouldBeTrue)
			c.So(msg, c.ShouldEqual, "I am a testing Msg")
			msg, ok = GetOk(ctx, "tolerance")
			c.So(ok, c.ShouldBeTrue)
			c.So(msg, c.ShouldEqual, 0.14)
			msg, ok = GetOk(ctx, "acceptableRank")
			c.So(ok, c.ShouldBeTrue)
			c.So(msg, c.ShouldResemble, []interface{}{"8", "9", "10"})
		})

		c.Convey("ok is false if no var is in json", func() {
			_, ok := GetOk(ctx, "UnknownVar")
			c.So(ok, c.ShouldBeFalse)
		})

		c.Convey("It retrieves a default var if the app is unknown", func() {
			ctx.SetAppId("Unknown-AppId")
			msg, ok := GetOk(ctx, "message")
			c.So(ok, c.ShouldBeTrue)
			c.So(msg, c.ShouldEqual, "I am a default Msg")
		})
		c.Convey("It uses default vars if there is no match in the current env", func() {
			ctx.SetAppId(TESTING_ENV_NAME)
			msg, ok := GetOk(ctx, "greeting")
			c.So(ok, c.ShouldBeTrue)
			c.So(msg, c.ShouldEqual, "default greeting")
		})
	})

	c.Convey("With a json file missing default properties", t, func() {
		ctx.SetAppId(TESTING_ENV_NAME)
		fileLocation := "./missing-default.json"
		err := Load(fileLocation)
		c.So(err, c.ShouldBeNil)

		c.Convey("It retrieves an available var", func() {
			msg, ok := GetOk(ctx, "message")
			c.So(ok, c.ShouldBeTrue)
			c.So(msg, c.ShouldEqual, "I am a testing Msg")
		})

		c.Convey("It errors if both the environment is unknown and no default vars are present", func() {
			ctx.SetAppId("Unknown-AppId")
			_, err := GetOk(ctx, "Message")
			c.So(err, c.ShouldNotBeNil)
		})
	})
}

func TestName(t *testing.T) {
	ctx := ContextMock{}
	ctx.SetAppId(TESTING_ENV_NAME)

	fileLocation := "./sample-environment.json"

	c.Convey("With a complete json file and mapping it retrieves the correct env name", t, func() {
		err := Load(fileLocation)
		c.So(err, c.ShouldBeNil)

		name := Name(ctx)
		c.So(name, c.ShouldEqual, "test")
	})

	c.Convey("With an unknown project id it returns the default environment", t, func() {
		err := Load(fileLocation)
		c.So(err, c.ShouldBeNil)

		ctx.SetAppId("Unknown-AppId")
		name := Name(ctx)
		c.So(name, c.ShouldEqual, "default")
	})
}

func TestIs(t *testing.T) {
	ctx := ContextMock{}
	ctx.SetAppId(TESTING_ENV_NAME)
	Load("./sample-environment.json")

	c.Convey("it returns accurate environment name mappings", t, func() {
		c.So(Is(ctx, "test"), c.ShouldBeTrue)
	})
}

func TestMustLoad(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fail()
		}
	}()

	MustLoad("./nonexistent.json")
}

type ContextMock struct {
	OverrideAppId string
}

func (env ContextMock) Debugf(format string, args ...interface{})    {}
func (env ContextMock) Infof(format string, args ...interface{})     {}
func (env ContextMock) Warningf(format string, args ...interface{})  {}
func (env ContextMock) Errorf(format string, args ...interface{})    {}
func (env ContextMock) Criticalf(format string, args ...interface{}) {}
func (env ContextMock) Call(service, method string, in, out appengine_internal.ProtoMessage, opts *appengine_internal.CallOptions) error {
	return nil
}
func (env ContextMock) Request() interface{} { return nil }

func (env ContextMock) FullyQualifiedAppID() string {
	return env.OverrideAppId
}

func (env *ContextMock) SetAppId(appId string) {
	env.OverrideAppId = appId
}
