push:
    bin/push-broker

run:
	go run ./main.go

manifest:
    export CONFIG_SERVER_BROKER_CONFIG=$(spruce merge broker_config.yml secrets.yml | spruce json )