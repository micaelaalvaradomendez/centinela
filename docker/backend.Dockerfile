FROM golang:1.27.1-alpine AS build

WORKDIR /src

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/centinela-api ./cmd/api

FROM alpine:3.22

RUN apk add --no-cache ca-certificates busybox-extras

COPY --from=build /out/centinela-api /usr/local/bin/centinela-api

EXPOSE 8080

USER 65532:65532

ENTRYPOINT ["/usr/local/bin/centinela-api"]