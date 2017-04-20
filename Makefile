build:
	go build -ldflags -s

run: build
	./slackpgp
