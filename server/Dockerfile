# Builder
FROM golang:1.20 as builder

WORKDIR /app

COPY . .

RUN make engine

# UPX
FROM hairyhenderson/upx:latest as upx

WORKDIR /app

COPY --from=builder /app/engine /app

RUN upx /app/engine

# Distribution
FROM frolvlad/alpine-glibc:latest

WORKDIR /app

COPY --from=upx /app/engine /app

# Expose port 8989
EXPOSE 9898

ENTRYPOINT ["/app/engine"]