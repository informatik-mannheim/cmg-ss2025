FROM golang:1.24.3-alpine3.21 AS builder

WORKDIR /app
COPY . ./

RUN go mod download
RUN go get ./...
RUN cd main && go build -o /consumer-cli .

  #step 2: copy the binary (and nothing else) into the run container

FROM alpine:3.21
  # Add CA certs (for HTTPS)
RUN apk add --no-cache ca-certificates

COPY --from=builder /consumer-cli /app/consumer-cli

ENV GATEWAY_URL=https://consumer-gateway.proudground-db418b69.germanywestcentral.azurecontainerapps.io

CMD ["/app/consumer-cli"]

