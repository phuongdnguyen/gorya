build:
	go build -o app  ./cmd

run:
	source scripts/local-env.sh
	go run cmd/* api
