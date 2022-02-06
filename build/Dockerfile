FROM golang:1.17.6-alpine

WORKDIR /app

COPY ./ ./

RUN ls -l ./

RUN go mod tidy
RUN go build ./cmd/generator


CMD [ "./generator" ]