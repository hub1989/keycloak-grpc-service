FROM golang:1.20-bullseye AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags="-w -s" -o /bin/app .

##
## Deploy
##
FROM scratch

WORKDIR /

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /bin/app /bin/app

#USER nonroot:nonroot

ENTRYPOINT ["/bin/app"]