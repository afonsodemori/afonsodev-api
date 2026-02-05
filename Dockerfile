FROM golang:1.25.6-alpine AS builder
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY go.mod ./
#COPY go.mod go.sum ./
#RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o afonso-dev-api .

FROM scratch
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/afonso-dev-api .
EXPOSE 8080
ENTRYPOINT ["./afonso-dev-api"]
