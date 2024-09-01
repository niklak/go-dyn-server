FROM golang:1.22.1-alpine AS build

RUN apk add --no-cache build-base bash
ENV APP_ROOT=/dyn-server
ENV APP_NAME=dyn-server

COPY . ${APP_ROOT}

WORKDIR ${APP_ROOT}
RUN GOOS=linux GOARCH=amd64 bash ./build-plugins.sh


WORKDIR ${APP_ROOT}/cmd/${APP_NAME}


RUN  go build -o ${APP_NAME}


#FROM debian:bookworm
FROM alpine:3.20

RUN apk add --no-cache \
	ca-certificates 

ENV USER=dyn-server
ENV APP_NAME=dyn-server

RUN adduser -D ${USER}
#RUN useradd ${USER}

COPY --chown=${USER}:${USER} --from=build /${APP_NAME}/cmd/${APP_NAME}/${APP_NAME} /usr/local/bin/${APP_NAME}
COPY --chown=${USER}:${USER} --from=build /${APP_NAME}/include /include

USER ${USER}

CMD ["dyn-server"]
