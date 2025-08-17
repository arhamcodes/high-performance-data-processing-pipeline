FROM golang:1.24-alpine
WORKDIR /high-performance-data-processing-pipeline/worker-go
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o process
CMD ["./process"]