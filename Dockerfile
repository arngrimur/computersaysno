# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.17.6-alpine3.15 AS build
ENV USER=nonroot
ENV UID=10001

RUN apk update && apk add --no-cache git
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /go/src

COPY . ./

RUN go mod download
RUN go mod tidy
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags "-s -w" -o /csn

##
## Deploy
##
FROM scratch
LABEL arch="amd64"

WORKDIR /

COPY --from=build /csn /csn
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group

#Expose port higher that 1000 so we don't need any priveleges
EXPOSE 10443

USER nonroot:nonroot

ENTRYPOINT ["/csn"]