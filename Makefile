.PHONY: run
run:
	cd cmd/ && go run main.go -conf ../configs/config.yaml

.PHONY: start-env
start-env:
	cd deploy/bluebell/ && docker-compose up -d

.PHONY: stop-env
stop-env:
	cd deploy/bluebell/ && docker-compose down

.DEFAULT_GOAL := help