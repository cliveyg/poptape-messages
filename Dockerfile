FROM golang:1.19-alpine AS build

RUN apk --no-cache add git

ENV GO111MODULE=on

RUN mkdir /app
ADD . /app
WORKDIR /app

# remove any go module files and get deps
RUN rm -f go.mod go.sum
RUN go mod init github.com/cliveyg/poptape-messages
RUN go mod tidy
RUN go mod download

# need these flags or alpine image won't run due to dynamically linked libs in binary
RUN CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build -a -ldflags '-w' -o messages


FROM alpine:latest

RUN mkdir -p /messages
COPY --from=build /app/messages /messages
COPY --from=build /app/.env /messages
WORKDIR /messages

# Make port 8090 available to the world outside this container
EXPOSE 8090

# Run messages binary when the container launches
CMD ["./messages"]
