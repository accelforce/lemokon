FROM golang:1.16.2-alpine3.13 AS build

WORKDIR /opt/lemokon/

COPY . /opt/lemokon/

ARG tags

RUN go build -tags "$tags" github.com/accelforce/lemokon

FROM alpine:3.13.2

COPY --from=build /opt/lemokon/lemokon /opt/lemokon/lemokon

CMD ["/opt/lemokon/lemokon"]
