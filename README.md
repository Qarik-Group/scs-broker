# Config Server Broker

This repo implements a service broker for config server.

## Demo

Here is a walkthrough of the complete provision / bind / unbind / deprovision story executed from the terminal.

Some configuration for the broker is pre-populated under [cf/broker_config.yml](./cf/broker_config.yml).

Before running the broker you must populate a file with sensitive information.
First we should create a uaa client that can create other clients for the `bind` oporation.
```
$ uaac client add config-server-broker --name config-server-broker --authorized_grant_types client_credentials --authorities clients.write,clients.read,clients.admin
```
When prompted enter the password.

Then fill out the file:
```
$ cat <<<EOF > cf/secrets.yml
broker_auth:
  user: admin
  password: <broker-password>
cloud_foundry_config:
  api_url: <cf-api-url>
  skip_ssl_validation: true
  cf_username: <cf-username>
  cf_password: <cf-password>
  uaa_client_id: <client-id-for-uaa>
  uaa_client_secret: "client-secret-for-uaa>"
instance_space_guid: "<space-guid-where-to-push-instances>"
instance_domain: "<cf-domain-to-use-for-routing-to-instance>"
EOF
```

Then load the config into the environment:
```
$ export CONFIG_SERVER_BROKER_CONFIG=$(spruce merge cf/broker_config.yml cf/secrets.yml | spruce json )
```
And finally run the broker:
```
$ make run
go run ./main.go
{"timestamp":"1592851636.363440990","source":"config-server-broker","message":"config-server-broker.Starting Config Server broker","log_level":1,"data":{}}
```

Once its running you can run through the service lifecycle manually using [eden](https://github.com/starkandwayne/eden).
First provision a service:
```
$ eden --client admin --client-secret <service-broker-password> --url http://localhost:8080 provision -s config-server -p basic -P '{"gitRepoUrl": "https://github.com/spring-cloud-samples/config-repo"}'
provision:   config-server/basic - name: config-server-basic-602effa1-13f5-408b-b25b-15a7d6aa2500
provision:   done
```

Then create a binding:
```
$ eden --client admin --client-secret <service-broker-password> --url http://localhost:8080 provision -s config-server -p basic -P '{"gitRepoUrl": "https://github.com/spring-cloud-samples/config-repo"}'
provision:   config-server/basic - name: config-server-basic-602effa1-13f5-408b-b25b-15a7d6aa2500
provision:   done
eden --client admin --client-secret admin --url http://localhost:8080 bind -i 602effa1-13f5-408b-b25b-15a7d6aa2500
Success

Run 'eden credentials -i config-server-basic-602effa1-13f5-408b-b25b-15a7d6aa2500 -b config-server-df0a2f93-9700-47a5-83b6-0ef14b3d76c7' to see credentials
$ eden --client admin --client-secret admin --url http://localhost:8080 credentials -i config-server-basic-602effa1-13f5-408b-b25b-15a7d6aa2500 -b config-server-df0a2f93-9700-47a5-83b6-0ef14b3d76c7
{
  "client_id": "config-server-binding-a553345a-4897-439e-a89e-f3aa7418ba3a",
  "client_secret": "Q5lv4CBapnzo3WOQANIGAB4u9Jq6Qa"
}
```

Now we can verify that things are working by testing wether we can access config via the provided credentials. For that we first need to get a token via uaac (input the client_secret when it asks for a password). Note that this currently only works with ssl validation turned on.

```
$ uaac token client get config-server-binding-a553345a-4897-439e-a89e-f3aa7418ba3a

Successfully fetched token via client credentials grant.
Target: https://uaa.hol.starkandwayne.com
Context: config-server-binding-a553345a-4897-439e-a89e-f3aa7418ba3a, from client config-server-binding-a553345a-4897-439e-a89e-f3aa7418ba3a
```

Then look up the token:
```
$ % uaac context config-server-binding-a553345a-4897-439e-a89e-f3aa7418ba3a

[0]*[https://uaa.hol.starkandwayne.com]
  skip_ssl_validation: true

  [1]*[config-server-binding-a553345a-4897-439e-a89e-f3aa7418ba3a]
      client_id: config-server-binding-a553345a-4897-439e-a89e-f3aa7418ba3a
      access_token: eyJhbGciOiJSUzI1NiIsImprdSI6Imh0dHBzOi8vdWFhLmhvbC5zdGFya2FuZHdheW5lLmNvbS90b2tlbl9rZXlzIiwia2lkIjoia2V5LTEiLCJ0eXAiOiJKV1QifQ.eyJqdGkiOiJlMWI1MzNhYjQzODY0NjYwYTc5NmNlMTM0MzI2Yjc4NyIsInN1YiI6ImNvbmZpZy1zZXJ2ZXItYmluZGluZy1hNTUzMzQ1YS00ODk3LTQzOWUtYTg5ZS1mM2FhNzQxOGJhM2EiLCJhdXRob3JpdGllcyI6WyJ1YWEubm9uZSJdLCJzY29wZSI6WyJ1YWEubm9uZSJdLCJjbGllbnRfaWQiOiJjb25maWctc2VydmVyLWJpbmRpbmctYTU1MzM0NWEtNDg5Ny00MzllLWE4OWUtZjNhYTc0MThiYTNhIiwiY2lkIjoiY29uZmlnLXNlcnZlci1iaW5kaW5nLWE1NTMzNDVhLTQ4OTctNDM5ZS1hODllLWYzYWE3NDE4YmEzYSIsImF6cCI6ImNvbmZpZy1zZXJ2ZXItYmluZGluZy1hNTUzMzQ1YS00ODk3LTQzOWUtYTg5ZS1mM2FhNzQxOGJhM2EiLCJncmFudF90eXBlIjoiY2xpZW50X2NyZWRlbnRpYWxzIiwicmV2X3NpZyI6IjM3YWE2OTAyIiwiaWF0IjoxNTkyODUyNjAzLCJleHAiOjE1OTI4OTU4MDMsImlzcyI6Imh0dHBzOi8vdWFhLmhvbC5zdGFya2FuZHdheW5lLmNvbS9vYXV0aC90b2tlbiIsInppZCI6InVhYSIsImF1ZCI6WyJjb25maWctc2VydmVyLWJpbmRpbmctYTU1MzM0NWEtNDg5Ny00MzllLWE4OWUtZjNhYTc0MThiYTNhIl19.ocG4LreYNLCsDgP-Xy2KYXgJ7qMC9zy6ok7YSOw-meSSx_X6jsbYqADK598BsHkxVBklhvUpn8FK0DdLKIvFJcYGOS3uZRw56_biaO1BTRR_akRopSRxNh4nR1qyNPtkGC2pfxqNeyPQDeTSXPVROfFv6TRJQri2AVpDt0r6uz0o45buZDaSr81oHKT2ExnKx5qsYScFHgY6bVRk-wmztpvEezwlzgE42g3y726610WjKYA19VhP1rVcBwY-kD9tPkfviaDFhbiaGQDR94a45kvoKAcK2QhBeGooIoBoOaJzrx8SSlseSJHMEOjvrfXhpRC0P5DfsOGqh3dFJ20U-Q
      token_type: bearer
      expires_in: 43199
      scope: uaa.none
      jti: e1b533ab43864660a796ce134326b787
$ TOKEN=eyJhbGciOiJSUzI1NiIsImprdSI6Imh0dHBzOi8vdWFhLmhvbC5zdGFya2FuZHdheW5lLmNvbS90b2tlbl9rZXlzIiwia2lkIjoia2V5LTEiLCJ0eXAiOiJKV1QifQ.eyJqdGkiOiJlMWI1MzNhYjQzODY0NjYwYTc5NmNlMTM0MzI2Yjc4NyIsInN1YiI6ImNvbmZpZy1zZXJ2ZXItYmluZGluZy1hNTUzMzQ1YS00ODk3LTQzOWUtYTg5ZS1mM2FhNzQxOGJhM2EiLCJhdXRob3JpdGllcyI6WyJ1YWEubm9uZSJdLCJzY29wZSI6WyJ1YWEubm9uZSJdLCJjbGllbnRfaWQiOiJjb25maWctc2VydmVyLWJpbmRpbmctYTU1MzM0NWEtNDg5Ny00MzllLWE4OWUtZjNhYTc0MThiYTNhIiwiY2lkIjoiY29uZmlnLXNlcnZlci1iaW5kaW5nLWE1NTMzNDVhLTQ4OTctNDM5ZS1hODllLWYzYWE3NDE4YmEzYSIsImF6cCI6ImNvbmZpZy1zZXJ2ZXItYmluZGluZy1hNTUzMzQ1YS00ODk3LTQzOWUtYTg5ZS1mM2FhNzQxOGJhM2EiLCJncmFudF90eXBlIjoiY2xpZW50X2NyZWRlbnRpYWxzIiwicmV2X3NpZyI6IjM3YWE2OTAyIiwiaWF0IjoxNTkyODUyNjAzLCJleHAiOjE1OTI4OTU4MDMsImlzcyI6Imh0dHBzOi8vdWFhLmhvbC5zdGFya2FuZHdheW5lLmNvbS9vYXV0aC90b2tlbiIsInppZCI6InVhYSIsImF1ZCI6WyJjb25maWctc2VydmVyLWJpbmRpbmctYTU1MzM0NWEtNDg5Ny00MzllLWE4OWUtZjNhYTc0MThiYTNhIl19.ocG4LreYNLCsDgP-Xy2KYXgJ7qMC9zy6ok7YSOw-meSSx_X6jsbYqADK598BsHkxVBklhvUpn8FK0DdLKIvFJcYGOS3uZRw56_biaO1BTRR_akRopSRxNh4nR1qyNPtkGC2pfxqNeyPQDeTSXPVROfFv6TRJQri2AVpDt0r6uz0o45buZDaSr81oHKT2ExnKx5qsYScFHgY6bVRk-wmztpvEezwlzgE42g3y726610WjKYA19VhP1rVcBwY-kD9tPkfviaDFhbiaGQDR94a45kvoKAcK2QhBeGooIoBoOaJzrx8SSlseSJHMEOjvrfXhpRC0P5DfsOGqh3dFJ20U-Q
```

And curl the config endpoint:
```
$ curl -v -H "Authorization: bearer ${TOKEN}" config-server-602effa1-13f5-408b-b25b-15a7d6aa2500.hol.starkandwayne.com/foo/foo
{"name":"foo","profiles":["foo"],"label":null,"version":"bb51f4173258ae3481c61b95b503c13862ccfba7","state":null,"propertySources":[{"name":"https://github.com/spring-cloud-samples/c
onfig-repo/foo.properties","source":{"foo":"from foo props","democonfigclient.message":"hello spring io"}},{"name":"https://github.com/spring-cloud-samples/config-repo/application.y
ml (document #0)","source":{"info.description":"Spring Cloud Samples","info.url":"https://github.com/spring-cloud-samples","eureka.client.serviceUrl.defaultZone":"http://localhost:8
761/eureka/","foo":"baz"}}]}* Closing connection 0
```
Once things are verified we can destroy everything again:
```
$ eden --client admin --client-secret admin --url http://localhost:8080 unbind -i config-server-basic-602effa1-13f5-408b-b25b-15a7d6aa2500 -b config-server-a553345a-4897-439e-a89e-f3aa7418ba3a
Success
$ eden --client admin --client-secret admin --url http://localhost:8080 deprovision -i 602effa1-13f5-408b-b25b-15a7d6aa2500
deprovision: config-server/basic - guid: 602effa1-13f5-408b-b25b-15a7d6aa2500
deprovision: done
```

## Image

The image that gets run (starkandwayne/spring-cloud-config-server:1.1.0) is defined on the [oauth](https://github.com/starkandwayne/spring-cloud-config-server/tree/oauth) branch of the S&W for of [hyness/spring-cloud-config-server](https://github.com/hyness/spring-cloud-config-server). It currently doesn't support any kind of `skip-ssl-verification` option so the `jwk_set_key` endpoint of the uaa must use ssl with a verifieable cert.
