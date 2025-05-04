FROM golang:1.20-alpine

# Install Tesseract OCR and Leptonica
RUN apk update && apk add --no-cache \
    tesseract-ocr \
    tesseract-ocr-dev \
    leptonica-dev \
    g++ \
    build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

# Build the Go application
RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
