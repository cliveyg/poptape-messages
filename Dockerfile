FROM golang:1.12-alpine as build

RUN apk --no-cache add git

ENV GO111MODULE=on

RUN mkdir /app
ADD . /app
WORKDIR /app

# get deps
RUN go mod init
RUN go mod tidy
RUN go mod download

# need these flags or alpine image won't run due to dynamically linked libs in binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-w' -o messages


FROM alpine:latest

RUN mkdir -p /messages
COPY --from=build /app/messages /messages
COPY --from=build /app/.env /messages
WORKDIR /messages

# Make port 8090 available to the world outside this container
EXPOSE 8090

# Run messages binary when the container launches
CMD ["./messages"]
