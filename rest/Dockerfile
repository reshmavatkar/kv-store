FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
# Download dependencies (cached if go.mod/go.sum unchanged)
RUN go mod download

# Copy all internal packages needed for build
COPY rest/ ./rest/
COPY store/ ./store/
COPY generated/ ./generated/

# Build from rest/main.go
WORKDIR /app/rest
RUN go build -o /app/rest-api main.go

EXPOSE 8080

CMD ["/app/rest-api"]
