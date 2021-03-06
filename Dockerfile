FROM golang:1.17-alpine AS build
# FROM golang:1.16.4
RUN apk add build-base

ADD ./plugins /src

RUN cd /src && go build  -buildmode=plugin -o krakend-debugger.so krakend-debugger.go

FROM devopsfaith/krakend:latest

USER root

COPY ./krakend.plugin.json krakend.json


COPY --from=build /src plugins

RUN krakend run -d -p 8080 
# RUN ls