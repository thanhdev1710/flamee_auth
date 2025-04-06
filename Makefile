APP_NAME=flamee_auth
ENTRY=cmd/server/main.go

.PHONY: run build tidy fmt clean

run:
	go run $(ENTRY)

build:
	go build -o $(APP_NAME) $(ENTRY)

tidy:
	go mod tidy

fmt:
	go fmt ./...

clean:
	rm -f $(APP_NAME)