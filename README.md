# OpsView-GO

OpsView-GO allows you to manage your OpsView Instance(s) and integrate your favorite tools (which are built in GO).

## Build requirements
We use docker to build and test, run this project. If you don't use docker you will need to install/configure GO locally.

You can use `USE_CONTAINER=false` environment setting for make to avoid using docker.
Otherwise make sure to have these tools:
- Docker (Client & Daemon)
- gnu make tools

## Testing your changes

### From a container
```
make test
```

### Without docker
* Install golang 1.8 or better
* Install go packages listed in .travis.yml
```
USE_CONTAINER=false make test

## Reference
https://knowledge.opsview.com/reference#api-intro
https://knowledge.opsview.com/reference#api-status-filtering-service-objects