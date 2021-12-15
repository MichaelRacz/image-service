FROM golang:1.17.5-alpine3.15

WORKDIR /go/src
COPY . /go/src/

RUN apk add --update docker 
RUN go build ./image-service.go

EXPOSE 8080
CMD ["./image-service"]

# NOTE: For prod usage, another build stage can
# be definded that slims down the image and does
# not run as root. 
