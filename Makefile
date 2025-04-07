APP_NAME=flamee_auth
ENTRY=cmd/server/main.go
PLATFORM_LINUX=linux
PLATFORM_MAC=darwin
PLATFORM_WIN=windows
ARCH_X64=amd64
ARCH_X86=386

.PHONY: run build tidy fmt clean build-mac build-win

# Command to run the app
run:
	go run $(ENTRY)

# Command to build the app for the host OS
build:
	go build -o $(APP_NAME) $(ENTRY)

# Command to build the app for macOS
build-mac:
	GOOS=$(PLATFORM_MAC) GOARCH=$(ARCH_X64) go build -o $(APP_NAME)-mac $(ENTRY)

# Command to build the app for Windows
build-win:
	GOOS=$(PLATFORM_WIN) GOARCH=$(ARCH_X64) go build -o $(APP_NAME)-win.exe $(ENTRY)

# Command to tidy the Go modules
tidy:
	go mod tidy

# Command to format the Go code
fmt:
	go fmt ./...

# Command to clean up the generated binary
clean:
	rm -f $(APP_NAME) $(APP_NAME)-mac $(APP_NAME)-win.exe
