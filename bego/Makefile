run-dev:
	go run . run-monitor

run-app:
	./EmMonitor run-monitor

run-web:
	./EmMonitor run-web

build-for-ubuntu:
	rm -rf EmMonitor
	GOOS=linux GOARCH=amd64 go build -o EmMonitor
	chmod +x EmMonitor

build-for-alpine:
	rm -rf EmMonitor
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o EmMonitor
	chmod +x EmMonitor