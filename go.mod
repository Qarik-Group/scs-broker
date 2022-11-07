module github.com/starkandwayne/scs-broker

go 1.14

replace github.com/docker/docker => github.com/docker/engine v17.12.0-ce-rc1.0.20200531234253-77e06fda0c94+incompatible

replace github.com/SermoDigital/jose => github.com/SermoDigital/jose v0.9.2-0.20161205224733-f6df55f235c2

replace github.com/mailru/easyjson => github.com/mailru/easyjson v0.0.0-20180323154445-8b799c424f57

replace github.com/cloudfoundry/sonde-go => github.com/cloudfoundry/sonde-go v0.0.0-20171206171820-b33733203bb4

replace code.cloudfoundry.org/go-log-cache => code.cloudfoundry.org/go-log-cache v1.0.1-0.20200316170138-f466e0302c34

require (
	code.cloudfoundry.org/cli v6.51.0+incompatible
	code.cloudfoundry.org/lager v2.0.0+incompatible
	github.com/cloudfoundry-community/go-cf-clients-helper v1.0.1
	github.com/cloudfoundry-community/go-cfclient v0.0.0-20220930021109-9c4e6c59ccf1
	github.com/cloudfoundry-community/go-uaa v0.3.1
	github.com/drewolson/testflight v1.0.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/gorilla/mux v1.7.4 // indirect
	github.com/mattn/go-colorable v0.1.11 // indirect
	github.com/pborman/uuid v1.2.0 // indirect
	github.com/pivotal-cf/brokerapi v6.4.2+incompatible
	github.com/starkandwayne/spring-cloud-services-cli-config-parser v1.0.2
	github.com/stretchr/testify v1.7.0 // indirect
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	golang.org/x/net v0.0.0-20220425223048-2871e0cb64e4 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)
