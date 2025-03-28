FROM golang:alpine AS builder
ENV GO111MODULE=on
RUN apk update && apk add --no-cache git
WORKDIR /server
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /server/cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server/main .
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=builder server/. .
WORKDIR /server/cmd/merch
COPY --from=builder server/cmd/server .
EXPOSE 8080
CMD ["server/main"]