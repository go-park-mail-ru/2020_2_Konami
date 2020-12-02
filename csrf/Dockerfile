FROM golang:latest AS build_step
ENV GO111MODULE=on
WORKDIR  /go/src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/build/main /go/src/cmd/csrf_dispatcher/main.go

FROM alpine
WORKDIR /app
COPY --from=build_step /go/build/main /app/main
RUN chmod +x /app/main
EXPOSE 8003/tcp
EXPOSE 6379/tcp
ENTRYPOINT /app/main
