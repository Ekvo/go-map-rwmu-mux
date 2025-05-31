FROM golang:1.24.1 AS builder

LABEL stage=builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /usr/src/build

ADD go.mod ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg

RUN go build -o quotebook ./cmd/app/main.go

FROM scratch

LABEL autors="ekvo"

ENV SERVER_HOST=0.0.0.0
ENV SERVER_PORT=8080

WORKDIR /usr/src/app

COPY --from=builder /usr/src/build/quotebook /usr/src/app/quotebook

EXPOSE ${SERVER_PORT}

ENTRYPOINT ["/usr/src/app/quotebook"]