FROM golang:1.16.3
WORKDIR /go/src/github.com/drgarcia1986/street-fair
COPY . .
RUN CGO_ENABLED=0 make build-api

FROM alpine:latest  
RUN apk --no-cache add ca-certificates \
 && apk --no-cache upgrade musl
WORKDIR /app/
COPY --from=0 /go/src/github.com/drgarcia1986/street-fair/street-fair .
CMD ["./street-fair"]
