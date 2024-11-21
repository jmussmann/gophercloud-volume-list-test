FROM fedora:41 AS build
RUN dnf install -y golang libnbd-devel
WORKDIR /src
COPY src/gophercloud-volume-list-test/ .
RUN go mod download
RUN go build -o /usr/local/bin/gophercloud-volume-list-test .

ENTRYPOINT ["/usr/local/bin/gophercloud-volume-list-test"]
