dev:
	export $$(xargs < .env.local) && go run src/cmd/v1/main.go
test:
	export $$(xargs < .env.test) && go run src/cmd/v1/main.go
build:
	env GOOS=linux GOARCH=amd64 go build ./src/cmd/v1