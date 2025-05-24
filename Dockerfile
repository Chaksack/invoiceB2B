FROM node:18-alpine AS nuxt-builder

WORKDIR /client

COPY client/package*.json ./
RUN npm install

COPY client/ .

RUN npm run generate

RUN echo "--- Contents of /client/.output/public/ after Nuxt build ---" && \
    ls -lA /client/.output/public/ && \
    echo "--- End of /client/.output/public/ contents ---"

FROM golang:1.23-alpine AS go-builder

WORKDIR /app

RUN apk update && apk add --no-cache sqlite build-base

RUN go install github.com/air-verse/air@latest

ENV CGO_ENABLED=1

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY --from=nuxt-builder /client/.output/public ./client/dist

RUN echo "--- Contents of /app/client/dist/ after COPY ---" && \
    ls -lA /app/client/dist/ && \
    echo "--- End of /app/client/dist/ contents ---"



EXPOSE 3000

CMD ["air", "-c", ".air.toml"]