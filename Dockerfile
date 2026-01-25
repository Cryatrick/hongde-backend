# ---- Build stage ----
FROM golang:1.24-bullseye AS builder

# Set the working directory inside the container
WORKDIR /

# Install build dependencies for CGO
RUN apt-get update && apt-get install -y \
    libwebp-dev \
    pkg-config \
    gcc \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# Ensure CGO is enabled (usually default, but be explicit)
ENV CGO_ENABLED=1

COPY /WebServer/CoreAPI/go.mod /WebServer/CoreAPI/go.sum ./
COPY /WebServer/CoreAPI/web/ ./web/
COPY /WebServer/CoreAPI/temp/ ./temp/

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY /WebServer/CoreAPI/ ./

# Build the Go application
RUN go build -o main ./cmd 

# ---- Runtime stage ----
FROM debian:bullseye-slim
RUN apt-get update && apt-get install -y \
    wkhtmltopdf \
    xfonts-75dpi xfonts-base \
    libxrender1 libjpeg62-turbo libpng16-16 \
    libssl1.1 fontconfig && \
    apt-get clean && rm -rf /var/lib/apt/lists/*
WORKDIR /

COPY --from=builder /main .
COPY --from=builder /.env .
COPY --from=builder /web ./web
COPY --from=builder /temp ./temp

# Expose the port the app runs on
EXPOSE 8090

# Command to run the application
CMD ["./main"]