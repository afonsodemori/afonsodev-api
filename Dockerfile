FROM golang:1.25.6-alpine AS builder

WORKDIR /app

COPY go.mod ./
#COPY go.mod go.sum ./

#RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o afonso-dev-api .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/afonso-dev-api .

EXPOSE 8080

ENTRYPOINT ["./afonso-dev-api"]
