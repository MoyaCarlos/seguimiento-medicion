FROM node:24-alpine AS frontend
WORKDIR /web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

FROM golang:1.26.5-alpine AS backend
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ cmd/
COPY internal/ internal/
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api
RUN mkdir /data

FROM gcr.io/distroless/static-debian13:nonroot
WORKDIR /app
COPY --from=backend /out/api /app/api
COPY --from=frontend /web/dist /app/web/dist
COPY --from=backend --chown=nonroot:nonroot /data /data
ENV PORT=8080 DB_PATH=/data/seguimiento.db WEB_DIR=/app/web/dist
VOLUME /data
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/api"]
