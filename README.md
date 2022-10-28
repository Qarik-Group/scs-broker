# Config Server Broker

This repo implements a service broker for config server.

## Setup

1. Clone repo to directory or home.
2. Create a Config-Server space (optional) under system.
```bash
$cf target -o system
$cf create-space config-server
$cf target -o system -s config-server
```
3. Modify the contents of cf/secrets.yml add the credentials needed, the instance guid for the space being used, etc.
#### Broker Auth Section
- **user** - *The user needed for creating the service-broker*
- **password** - *The password you create for creating the service-broker*
#### Cloud Foundry Config Section
- **api_url** - *The Cloud Foundry API URL*
- **cf_username** - *The user in CF that has permissions to create services with the API*
- **cf_password** - *the password for above user*
- **uaa_client_id** - *The Client ID, which can possibly be found in CredHub. If not the UAA client needs to be created*
- **uaa_client_secret** - *The secret for the client id above*
- **instance_space_guid** - *The GUID for the current space, This can be found by doing `cf space config-server --guid`*
- **instance_domain** - *The instance domain which is likely the same as the API base domain*
- **config_server_download_uri** - *This can be either a file:// or an https:// protocol. The use of file will look within the service-broker itself and should reference the location of the config-server.jar that has been packaged along with the broker*

4. While within the config-server-broker base directory, run `make push` if the binary doesn't need to be rebuilt. The only reason the binary may need to be rebuilt is if there is a code change. To build, run `make build` first, then run `make push`. All jars should be located under the app/artifacts directory, the binary should be located under the app directory. The app directory is what gets pushed to the CF container.

5. Create the service broker.
```
$ cf create-service-broker [user from the Broker Auth section] [password from the Broker Auth section] [App URL created for the running app]
```
6. Once the push is complete, the service broker is now running. You should now be able to create a service like the following example.
```
$cf create-service config-server default test-service -c "whatever json configuration you wish to use for config-server - see config-server docs from Spring.io"
```


## History ##

* v0.0.2 - Now with package fetching
* v0.0.1 - Initial Release
