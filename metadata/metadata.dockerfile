FROM alpine:latest
RUN mkdir /app
COPY metadataApp /app
CMD [ "/app/metadataApp"]
