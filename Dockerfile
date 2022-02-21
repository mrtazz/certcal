FROM golang:1.17-alpine3.15 AS builder

RUN mkdir /app

WORKDIR /app

ADD . /app/

RUN go build certcal.go -o certcal

FROM alpine:3.15
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/certcal /app/certcal
RUN chown -R nobody /app && chmod +x /app/certcal

USER nobody
ENV PORT=3000
ENV CERTCAL_HOSTS="unwiredcouch.com,github.com"
ENTRYPOINT ["/app/certcal"]
