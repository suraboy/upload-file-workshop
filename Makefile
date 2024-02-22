run:
	go run app/core/cmd/main.go

tidy:
	export GOSUMDB=off ; go mod tidy