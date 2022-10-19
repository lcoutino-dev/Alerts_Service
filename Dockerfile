FROM golang:1.14-alpine as build

RUN apk add --no-cache git

WORKDIR /src 
RUN go get -d -v ./...
COPY . .
RUN go build consumer.go


FROM alpine as runtime
RUN addgroup -S app && adduser -S -G app app
RUN mkdir /app && chown -R app:app /app
COPY --chown=app:app --from=build /src/consumer /app/consumer
USER app
CMD [ "/app/consumer" ]
