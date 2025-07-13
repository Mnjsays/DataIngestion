
FROM  golang:latest AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN echo $(ls)
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o main ./cmd/server/main.go

FROM alpine:latest AS release-stage
WORKDIR /app/di
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/main .
COPY .env .
EXPOSE 8094
CMD ["./main"]