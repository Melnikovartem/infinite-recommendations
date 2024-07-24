FROM golang:1.22-alpine

WORKDIR /app

# Install build dependencies if needed
RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

ENV GORSE_SERVER_HOST=localhost
ENV GORSE_SERVER_PORT=8087
ENV GORSE_API_KEY=""

COPY . .
RUN go build -o main main.go

EXPOSE 8080

CMD ["./main"]