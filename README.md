# keycloak-grpc-service

## building for docker
```shell
 docker buildx build --platform=linux/arm64 -o type=docker --build-arg PLATFORM=arm64 -t keycloak-grpc-service:arm64 --no-cache .
```
