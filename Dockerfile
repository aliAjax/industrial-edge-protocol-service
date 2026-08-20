FROM golang:1.23 AS build
WORKDIR /src
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/edge-api ./cmd/edge-api
FROM gcr.io/distroless/static-debian12
COPY --from=build /out/edge-api /edge-api
EXPOSE 8099
ENTRYPOINT ["/edge-api"]
