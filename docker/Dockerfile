# ---------- Build stage ----------
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/sungrow-go-docker

# ---------- Final stage ----------

FROM alpine
LABEL MAINTAINER="stefan.roeck@gmail.com"

RUN mkdir /sungrow
COPY --from=builder /app/bin/sungrow-go-docker /sungrow/sungrow-go-docker
WORKDIR /sungrow

COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]