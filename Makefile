PROC_NAME = echogy
RELEASE_PATH = release
PACKAGE_PATH = release
SERVER_PATH = cmd
DEBUGGER_PATH = debugger

install:
	@go get

build-cgo:
	cd ${SERVER_PATH} &&  GOOS=linux GOARCH=amd64 go build -race -ldflags "-s -w" -o $(RELEASE_PATH)/${PROC_NAME}

build:
	cd ${SERVER_PATH} &&  GOOS=linux GOARCH=amd64 go build -o $(RELEASE_PATH)/${PROC_NAME}

build-web:
	cd ${DEBUGGER_PATH} && npm install && npm run build && rm -rf dist/stats.html

clean:
	rm -rf dist

	rm -rf ${DEBUGGER_PATH}/dist
	# clean package
	rm -rf ${PACKAGE_PATH}
	# clean server build
	rm -rf ${SERVER_PATH}/${RELEASE_PATH}


package: clean build-web build
	rm -rf ${PACKAGE_PATH}
	mkdir -p ${PACKAGE_PATH}
	cp ${SERVER_PATH}/${RELEASE_PATH}/${PROC_NAME}  ${PACKAGE_PATH}
	cp config.sample.json  ${PACKAGE_PATH}/config.json
	cp -a bin/ ${PACKAGE_PATH}
