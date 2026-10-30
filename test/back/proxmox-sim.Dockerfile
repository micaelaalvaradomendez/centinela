# Simulador de Proxmox VE (backend/cmd/proxmox-simulador) para la suite de aceptación.
# Se usa para verificar BAC-30 (creación, config y rrddata) sin depender del Proxmox real.
FROM golang:1.27.1-alpine AS build

WORKDIR /src

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/proxmox-simulador ./cmd/proxmox-simulador

FROM alpine:3.22

COPY --from=build /out/proxmox-simulador /usr/local/bin/proxmox-simulador

EXPOSE 8081

ENTRYPOINT ["/usr/local/bin/proxmox-simulador"]
