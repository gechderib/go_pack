FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /salon_db 
EXPOSE 8080
CMD [ "/salon_db" ]

# # Build stage
# FROM golang:1.22 AS builder

# WORKDIR /app
# COPY . .

# RUN go mod tidy
# RUN go build -o app ./cmd/api

# # Run stage
# FROM alpine:latest

# WORKDIR /root/
# COPY --from=builder /app/app .

# EXPOSE 8080

# CMD ["./app"]