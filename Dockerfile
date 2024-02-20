FROM  golang:1.21-alpine AS build

RUN apk add --no-cache git

WORKDIR /tmp/go-opentelemetry

COPY src .

RUN go mod download

RUN go build -o ./out/go-opentelemetry .

FROM alpine:3.18

COPY --from=build /tmp/go-opentelemetry/out/go-opentelemetry /app/go-opentelemetry

EXPOSE 8080


CMD ["/app/go-opentelemetry"]