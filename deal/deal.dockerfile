FROM alpine:latest
RUN mkdir /app
COPY dealApp /app
CMD [ "/app/dealApp"]
