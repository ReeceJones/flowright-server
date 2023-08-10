FROM golang:alpine

RUN apk add build-base
# isntall podman
RUN apk add podman
RUN apk add fuse-overlayfs

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -v -o /usr/local/bin/ ./...

EXPOSE 9000

ENV PATH $PATH:/usr/local/bin/flowright

VOLUME [ "/var/lib/containers" ]
ADD ./containers.conf /etc/containers/containers.conf


CMD ["proxy", "-host", "0.0.0.0", "-control-plane-host", "control_plane", "-public-host", "proxy"]