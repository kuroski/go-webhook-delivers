FROM golang:1.23.5 AS builder
WORKDIR /app
COPY . .
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go build -o app ./cmd/web/

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 80
CMD ["./app", "-addr", ":80"]