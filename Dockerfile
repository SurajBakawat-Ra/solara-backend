
# Simple multi-stage Dockerfile for the Go API
FROM golang:1.22 AS build
WORKDIR /app
COPY . .
RUN go mod init solara-backend || true
RUN go mod tidy || true
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server

FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=build /app/server /server
ENV PORT=8080
ENV ALLOW_ORIGIN=*
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/server"]
