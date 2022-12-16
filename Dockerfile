FROM golang:1.18-alpine AS builder
WORKDIR /source
COPY . /source
RUN go mod download && go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o lifo-rest-api ./cmd/app/.

FROM alpine:3.9
RUN mkdir /app
WORKDIR /app
RUN mkdir configs && mkdir schema
COPY --from=builder /source/configs /app/configs
COPY --from=builder /source/schema /app/schema
COPY --from=builder /source/lifo-rest-api /usr/local/bin
RUN chmod a+x /usr/local/bin/lifo-rest-api


ENTRYPOINT [ "lifo-rest-api" ]
