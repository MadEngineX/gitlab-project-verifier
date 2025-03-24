ARG PROXY_REGISTRY

FROM $PROXY_REGISTRY/golang:1.22.4 AS dependency

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy

FROM dependency as build
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o /out/verifier cmd/verifier/main.go

FROM $PROXY_REGISTRY/alpine:3.19.1 AS runtime
WORKDIR /app
COPY --from=build /out/verifier /app/
ENTRYPOINT ["/app/verifier"]
