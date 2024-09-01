FROM golang:1.22.1-bookworm AS build


ENV APP_ROOT=/dyn-server
ENV APP_NAME=dyn-server

COPY . ${APP_ROOT}
WORKDIR ${APP_ROOT}/cmd/${APP_NAME}

RUN go build -o ${APP_NAME}


FROM debian:bookworm 

ENV USER=dyn-server
ENV APP_NAME=dyn-server

#RUN adduser -D ${USER}
RUN useradd ${USER}

COPY --chown=${USER}:${USER} --from=build /${APP_NAME}/cmd/${APP_NAME}/${APP_NAME} /usr/local/bin/${APP_NAME}

USER ${USER}

CMD ["dyn-server"]
