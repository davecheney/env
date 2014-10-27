# Makefile for this GAE/Go project

APPENGINE_VERSION=1.9.13
APPENGINE_FILE=go_appengine_sdk_linux_amd64-$(APPENGINE_VERSION).zip
SRC=src
ALLOCATION=$(SRC)/allocation
APPENGINE=go_appengine
APPENGINE_SRC="./.appengine-src-${APPENGINE_VERSION}"
VENDOR_GOPATH=.vendor
LOCAL_REMOTE=http://localhost:8080/_ah/remote_api
PROD_REMOTE=https://rockfish-project.appspot.com/_ah/remote_api
PACKAGES=env
GPM_URL="https://raw.githubusercontent.com/markmandel/gpm/feature/new-flags/bin/gpm"

#discover the OS to override Linux defaults for Darwin.
UNAME := $(shell uname)

ifeq ($(UNAME), Darwin)
        APPENGINE_FILE=go_appengine_sdk_darwin_amd64-$(APPENGINE_VERSION).zip
endif

configure:
	if [ ! -d ${APPENGINE_SRC} ];then mkdir ${APPENGINE_SRC}; fi
	if [ ! -f "${APPENGINE_SRC}/${APPENGINE_FILE}" ]; then \
		wget -P ${APPENGINE_SRC} "https://storage.googleapis.com/appengine-sdks/featured/${APPENGINE_FILE}"; \
	fi;
	rm -rf go_appengine
	unzip -q "${APPENGINE_SRC}/$(APPENGINE_FILE)"	
	if [ ! -f ${APPENGINE_SRC}/gpm ]; then \
		wget -P ${APPENGINE_SRC} ${GPM_URL}; \
	fi
	cp ${APPENGINE_SRC}/gpm ${APPENGINE}	
	chmod +x $(APPENGINE)/gpm
	$(MAKE) clean-deps

clean-configure:
	rm -rf $(APPENGINE)

clean-all:
	$(MAKE) clean-configure clean-deps

clean-deps:
	rm -rf $(VENDOR_GOPATH)
	mkdir $(VENDOR_GOPATH)
	rm -rf pkg

fmt:
	gofmt -w -s $(SRC)

deps:
	gpm install -g "$(CURDIR)/$(APPENGINE)/goapp" -u 0
	# build binaries as necessary.
	goapp install github.com/smartystreets/goconvey

test:
	goapp test $(PACKAGES)

convey:
	cd src && goconvey -gobin="$(CURDIR)/$(APPENGINE)/goapp" -port=7000 -depth=3

serve:
	dev_appserver.py --storage_path=./data ./$(SRC)/notifications_module/app.yaml

doc:
	killall godoc; godoc -http=":7080" &

deploy:
	appcfg.py --oauth2 --application=$(PROJECT) update $(SRC)/notifications_module/app.yaml
	appcfg.py --oauth2 --application=$(PROJECT) update_cron $(SRC)/notifications_module/

