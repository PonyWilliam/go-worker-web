FROM alpine:latest

RUN mkdir /app
WORKDIR /app
ADD WorkWeb /app/WorkWeb

CMD ["./WorkWeb"]