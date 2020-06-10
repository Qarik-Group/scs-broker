push:
	cf push -f cf/manifest.yml --vars-file cf/vars.yml

run:
	go run ./main.go
