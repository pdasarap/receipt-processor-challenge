FROM golang:1.22.5-alpine

RUN mkdir /app
WORKDIR /app
COPY . .

RUN go build -o receipt_processor

EXPOSE 8080
CMD ["./receipt_processor"]