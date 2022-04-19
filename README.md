# cobra-starter

## github actions secrets

[dockerhub access_token](https://hub.docker.com/settings/security)  
[dockerhub account settings](https://hub.docker.com/settings/general)
```env
DOCKER_HUB_USERNAME={{ your dockerhub username}}
DOCKER_HUB_ACCESS_TOKEN={{ your dockerhub access_token}}
```


## Docker

```bash
docker run ghstahl/cobra_starter 
docker run ghstahl/cobra_starter version
```

## ASYNQ Tests

### Docker-Compose

```bash
docker-compose up
```
Navigate to [asynqmon](http://asynqmon.docker.localhost/)  

Redis and the asyncq monitor UI should be up.

### Run Several Consumers

### 2 Servers 

These 2 will listen to the default queues;
```go
[]string{"critical:6", "default:3", "low:1"}
```

```bash
go run .\cmd\cli tasks handler 
```
```bash
go run .\cmd\cli tasks handler 
```

### 2 Other Servers 

These 2 will listen to a single herb queue; 
```go
[]string{herb:1}
```

```bash
go run .\cmd\cli\ tasks handler -q herb:1
```
```bash
go run .\cmd\cli\ tasks handler -q herb:1 
```

## Publish some messages

We need to send at a minimum 1000 through because they get processed so fast that a single server will get all of them if it is just 10 or so messages.  


```bash
go run .\cmd\cli\ tasks publisher -q critical -c 1000
go run .\cmd\cli\ tasks publisher -q default -c 1000
go run .\cmd\cli\ tasks publisher -q low -c 1000
```

```bash
go run .\cmd\cli\ tasks publisher -q herb -c 1000
```

