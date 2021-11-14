ARG APPNAME=similar-server

FROM golang:alpine AS build
RUN apk add --no-cache --update git
ADD . /go/src/app
WORKDIR /go/src/app/main
RUN go generate
RUN go build --mod=vendor -o app

FROM alpine:latest

EXPOSE 8080

RUN apk --no-cache add ca-certificates tzdata
COPY --from=build /go/src/app/main/app /usr/bin/app
USER nobody
CMD ["/usr/bin/app"]