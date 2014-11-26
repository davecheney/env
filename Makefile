PACKAGE_ROOT="github.com/rockpoollabs/env"
PACKAGES=$(PACKAGE_ROOT)/...

doc:
	pkill godoc; godoc -http=":7080" &

deps:
	goapp get -u github.com/smartystreets/goconvey/convey

test:
	goapp test $(PACKAGES)

live-test:
	goconvey

fmt:
	echo "Running gofmt on everything ..."
	goapp fmt -x $(PACKAGES)