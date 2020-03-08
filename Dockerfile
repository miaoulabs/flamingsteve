FROM golang:1.14-alpine as modules

RUN apk add ca-certificates git gcc g++ alsa-lib-dev libc-dev

WORKDIR /src

COPY go.mod .
COPY go.sum .

# go list -m all &&
RUN go mod download

FROM modules as build

#ENV GOOS=linux
#ENV GOARCH=amd64

WORKDIR /src

COPY . .

RUN go build -v -o /opt/dispmatrix ./cmd/dispmatrix && \
    go build -v -o /opt/dispsim ./cmd/dispsim && \
    go build -v -o /opt/glue ./cmd/glue && \
    go build -v -o /opt/muthur ./cmd/muthur && \
    go build -v -o /opt/sensim ./cmd/sensim && \
    go build -v -o /opt/sensor ./cmd/sensor && \
    go build -v -o /opt/senspad ./cmd/senspad && \
    go build -v -o /opt/sensui ./cmd/sensui && \
    go build -v -o /opt/seq ./cmd/seq

FROM alpine:3

RUN apk add --no-cache ca-certificates bash

WORKDIR /bin

EXPOSE 8080
COPY --from=build /opt/ /bin

ENTRYPOINT ["bash"]
