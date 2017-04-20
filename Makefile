build:
	time go build -x -ldflags -s

run: build
	./slackpgp
