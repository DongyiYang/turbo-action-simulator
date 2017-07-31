FROM alpine:3.3

ADD ./_output/turbo-simulator.linux /bin/turbo-simulator
EXPOSE 8087

ENTRYPOINT ["/bin/turbo-simulator"]
