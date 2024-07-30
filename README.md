# DDNS Client

## How to install DDNS-CLIENT

Simply run the command `go get github.com/datahearth/ddnsclient`

## Run the client

You have 2 options to run the DDNS client.  
You can run it as: 
- docker container:
`docker run -v /path/to/config/ddnsclient.yaml:/ddnsclient.yaml --name ddnsclient datahearth/ddnsclient:latest`  

- binary executable:
`./ddnsclient` (make sure the config is in the same directory with the name `ddnsclient.yaml`)

## Supported providers

- OVH
- Google (only one subdomain accepted for now)

Note: For now, ddnsclient supports only one credential for each provider.

## Use the library

You can also plug the library to your own system. Just get the module and you'll find everything needed to start it.
If something is missing or is not working properly, please create an issue so I can fix it.

## Contributing

You can contribute to the project by submitting an issue and resolve issues by creating PRs. I'll look at them and validate your changes if they're correct. 