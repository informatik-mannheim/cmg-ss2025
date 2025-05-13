FROM golang:1.24.1

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o worker-daemon ./cmd

CMD ["./worker-daemon"]