# Start from golang base image
FROM golang:alpine

# Add maintainer info
LABEL maintainer="Naufal Fawwaz Andriawan"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

CMD ["./main"]