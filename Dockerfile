# Stage 1 - BUILD
FROM golang:1.18.3-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Stage 2 - RUN
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY application.env .

EXPOSE 8080
CMD [ "/app/main" ]
