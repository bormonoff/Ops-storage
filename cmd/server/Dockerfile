FROM golang:1.22 AS build

ENV APP_HOME=/app

COPY cmd/ $APP_HOME/cmd/
COPY internal/ $APP_HOME/internal/
COPY go* Makefile $APP_HOME

WORKDIR $APP_HOME

RUN make build

FROM ubuntu:24.04 AS deploy

RUN groupadd -r metrics && useradd -r -g metrics muser
USER muser

COPY --from=build /app/build/server /app/

ENTRYPOINT ["/app/server"]
