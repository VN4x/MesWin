# Build context: repo root
FROM golang:1.22-bookworm AS backend-build
WORKDIR /src/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 go build -o /meswin .

FROM node:22-bookworm AS frontend-build
WORKDIR /src/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=backend-build /meswin /app/meswin
COPY --from=frontend-build /src/frontend/build /app/pb_public
ENV MESWIN_PUBLIC_DIR=/app/pb_public
EXPOSE 8090
VOLUME ["/app/pb_data"]
ENTRYPOINT ["/app/meswin"]
CMD ["serve", "--http=0.0.0.0:8090"]
