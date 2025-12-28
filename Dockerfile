FROM golang:1.25-alpine AS builder

WORKDIR /build


COPY go.mod go.sum* ./
RUN go mod download

COPY . .


RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /todo-app ./cmd/main.go


FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

RUN adduser -D appuser
USER appuser

WORKDIR /app

COPY --from=builder /todo-app .

EXPOSE 8080

CMD ["./todo-app"]