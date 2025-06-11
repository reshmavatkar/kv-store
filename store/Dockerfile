FROM golang:1.24-alpine

WORKDIR /app

# Copy only go.mod and go.sum first (to cache dependencies layer)
COPY go.mod go.sum ./

# Download dependencies (cached if go.mod/go.sum unchanged)
RUN go mod download

# Copy generated package and store source
COPY store/ ./store/
COPY generated/ ./generated/

WORKDIR /app/store
RUN go build -o /app/grpc-store main.go

EXPOSE 50051

CMD ["/app/grpc-store"]
