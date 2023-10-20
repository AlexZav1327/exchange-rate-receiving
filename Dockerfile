FROM golang:latest as builder
ADD . /src/app
WORKDIR /src/app
RUN CGO_ENABLED=0 GOOS=linux go build -o currencies-service ./cmd/currencies-service/main.go
EXPOSE 8080

FROM alpine:edge
COPY --from=builder /src/app/currencies-service /currencies-service
RUN chmod +x ./currencies-service
ENTRYPOINT ["/currencies-service"]