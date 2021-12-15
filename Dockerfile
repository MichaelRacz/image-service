FROM golang:1.17.5-alpine3.15

RUN apk add --update docker 
#openrc
# RUN apk add openrc
# RUN openrc 
# RUN addgroup root docker
# RUN /etc/init.d/docker start
# RUN service docker start
# RUN rc-update add docker boot

WORKDIR /go/src
COPY . /go/src/
RUN go build ./image-service.go

EXPOSE 8080
CMD ["./image-service"]

# NOTE: For prod usage, another build stage can
# be definded that slims down the image and does
# not run as root. 
