FROM alpine
ADD WorkWeb-service /WorkWeb-service
ENTRYPOINT [ "/WorkWeb-service" ]
