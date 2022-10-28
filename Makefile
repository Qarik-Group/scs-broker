push:
	bin/push-broker

run:
	go run ./main.go

manifest:
	export SCS_BROKER_CONFIG=$(spruce merge broker_config.yml secrets.yml | spruce json )

