push:
	cf push -f cf/manifest.yml

run:
	go run ./main.go
