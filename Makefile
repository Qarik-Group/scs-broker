push:
	mkdir tmp || true
	echo "config_server_broker_config: '$(spruce merge cf/broker-config.yml cf/secrets.yml | spruce json)'" > tmp/vars.yml
	cf push -f cf/manifest.yml --vars-file tmp/vars.yml

run:
	go run ./main.go
