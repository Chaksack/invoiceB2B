FROM golang:1.23-alpine
LABEL authors="Chakdahah"

WORKDIR /app

RUN apk update && apk add --no-cache sqlite build-base

RUN go install github.com/air-verse/air@latest

ENV CGO_ENABLED=1

COPY go.mod go.sum ./
RUN go mod download

COPY . .

#COPY paysenta.db .

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]