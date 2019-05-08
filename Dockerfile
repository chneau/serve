# build stage
ARG BASE=/go/src/github.com/chneau/serve
FROM golang:alpine AS build-env
ARG BASE
ADD . $BASE
RUN cd $BASE && CGO_ENABLED=0 go build -o /serve -ldflags '-s -w -extldflags "-static"'

FROM alpine AS prod-ready
RUN apk add --no-cache ca-certificates
COPY --from=build-env /serve /serve
ENTRYPOINT [ "/serve" ]
