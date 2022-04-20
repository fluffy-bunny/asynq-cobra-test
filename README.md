# cobra-starter

## Install golang 1.18

[installers](https://go.dev/dl/)

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
go run ./cmd/cli tasks handler 
```
```bash
go run ./cmd/cli tasks handler 
```

### 2 Other Servers 

These 2 will listen to a single herb queue; 
```go
[]string{herb:1}
```

```bash
go run ./cmd/cli tasks handler -q herb:1
```
```bash
go run ./cmd/cli tasks handler -q herb:1 
```

## Publish some messages

We need to send at a minimum 1000 through because they get processed so fast that a single server will get all of them if it is just 10 or so messages.  


```bash
go run ./cmd/cli tasks publisher -q critical -c 1000
go run ./cmd/cli tasks publisher -q default -c 1000
go run ./cmd/cli tasks publisher -q low -c 1000
```

```bash
go run ./cmd/cli tasks publisher -q herb -c 1000
```

## Fail to process a message

Lets spin up a server that fails to process a task.  It returns an error on everything.  

```bash
go run ./cmd/cli tasks handler -f 
```

What I have notices.
1. I spun up a failing server  
2. I sent 11 messages
3. The failing server failed them all and then went into a retry cycle
4. I spun up another server that would not fail.
5. I noticed that the other server was NOT pulling the failed messages.  It was like they were locked to the failing server
6. When I shut down the failing server and left the good server up, the good server started processing (slowly).   It looks like there is some pretty good backoff logic at work here.  

Same test, but 1000 upfront messages.
This time bringing up the second server helped and both were running.  The good server was now handling the messages.  

However once the numbers got down to around 40 unprocesses failed messages the failing server seems to camp on that small set and every now and then the good server will peel one off.  From the looks of it I may be sitting here for hours hoping the good server finally gets to process them all.  







