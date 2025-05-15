#step 1: build the service in a builder container
FROM golang:1.24-alpine3.21 AS builder
WORKDIR /app
COPY . ./
RUN go mod download
RUN go get ./...
RUN go build -o /job-scheduler .

#step 2: copy the binary (and nothing else) into the run container
#FROM alpine:3.21
FROM scratch
EXPOSE 8080
COPY --from=builder /job-scheduler /app/job-scheduler
CMD ["/app/job-scheduler"]
