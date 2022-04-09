#use small build environment
FROM golang:alpine as build 

WORKDIR /app

#get dependencies
COPY go.mod ./
COPY go.sum ./

#copy and build the project
COPY . ./
RUN go build -o /main cmd/main.go

#Use deploy env
FROM alpine:latest

WORKDIR /

COPY --from=build /main /
ENTRYPOINT [ "/main" ]