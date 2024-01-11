FROM golang:1.21-alpine3.19 AS builder

ARG APP_NAME
ARG VARIANT
ENV APP_NAME $APP_NAME
ENV VARIANT $VARIANT

WORKDIR /build

COPY . .

RUN go mod download

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/$VARIANT/$APP_NAME


FROM scratch

COPY --from=builder /build/app /

ENTRYPOINT ["/app"]