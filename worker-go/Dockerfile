FROM golang:1.22-alpine AS builder
WORKDIR /high-performance-data-processing-pipeline/worker-go
COPY . .
RUN go mod init worker-go
RUN CGO_ENABLED=0 go build -o process
FROM alpine:latest
WORKDIR /high-performance-data-processing-pipeline/worker-go
COPY --from=builder /high-performance-data-processing-pipeline/worker-go/process .
CMD ["./process"]