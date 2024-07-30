# DDNS Client

## How to install DDNS-CLIENT as dependency

Simply run the command `go get github.com/datahearth/ddnsclient`

## Run the client

You have 2 options to run the DDNS client.  
You can run it as: 
- docker container:  
```
docker run -v /path/to/config/ddnsclient.yaml:/ddnsclient.yaml --name ddnsclient ghcr.io/datahearth/ddnsclient:latest
```
or with a custom config path:  
```
docker run -e CONFIG_PATH=/path/inside/container/custom.yaml -v /path/to/config/ddnsclient.yaml:/path/inside/container/custom.yaml --name ddnsclient ghcr.io/datahearth/ddnsclient:latest
```

- binary executable:
```
git clone https://github.com/datahearth/ddnsclient.git
cd ddnsclient
go build -o ddnsclient cmd/main.go
./ddnsclient
```
make sure the config is in the same directory with the name `ddnsclient.yaml` or set the `CONFIG_PATH` variable

## Supported providers

Any provider using the standard for DDNS should be supported by default thanks to the generic configuration.  
You just need to get your credentials (obviously) and the update URL.  
If you face any kind of issue, feel free to open an issue and ping me in it. If necessary, a branch will be open to fix the problem.  

| Provider   	| Configuration key 	| Implemented 	| Tested 	|
|------------	|-------------------	|-------------	|--------	|
| OVH        	| ovh               	| YES         	| YES    	|
| GOOGLE     	| google            	| YES         	| YES    	|
| DuckDNS    	| duckdns           	| YES         	| YES     |
| No-IP      	| noip              	| YES          	| NO     	|
| DynDNS     	| dyndns            	| YES          	| NO     	|
| CloudFlare 	| cloudflare        	| NO          	| NO     	|

Note: 
For DDNS providers using basic authentication inside URL (e.g: `https://{username}:{password}@ddns.something.com/...`), remove the `username`and `password` part to get only the "classical" URL (e.g: `https://ddns.something.com/...`). Then fill the `username` and `password` fields in the provider configuration.

## Contributing

You can contribute to the project by submitting an issue and resolve issues by creating PRs. I'll look at them and validate your changes if they're correct as soon as possible. 

## TO-DO

- Add HRM to configuration file
- Add more DDNS provider (see the table above)

## Useful links
- Google DDNS doc: https://support.google.com/domains/answer/6147083?hl=en#zippy=%2Cusing-the-api-to-update-your-dynamic-dns-record
- OVH DDNS doc: https://docs.ovh.com/us/en/domains/hosting_dynhost/
- DuckDNS DDNS doc: https://www.duckdns.org/spec.jsp
- No-IP DDNS doc: https://www.noip.com/integrate/request
- DynDNS DDNS doc: https://help.dyn.com/remote-access-api/perform-update/
