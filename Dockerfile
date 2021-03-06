FROM golang:alpine AS build

RUN mkdir -p /go/src/github.com/zdnscloud/csi-lvm-plugin
COPY . /go/src/github.com/zdnscloud/csi-lvm-plugin

WORKDIR /go/src/github.com/zdnscloud/csi-lvm-plugin
RUN CGO_ENABLED=0 GOOS=linux go build -o lvmcsi cmd/lvmcsi.go


FROM alpine

LABEL maintainers="Kubernetes Authors"
LABEL description="LVM CSI Plugin"

RUN apk update && apk --no-cache add blkid file util-linux e2fsprogs
COPY --from=build /go/src/github.com/zdnscloud/csi-lvm-plugin/lvmcsi /lvmcsi

ENTRYPOINT ["/lvmcsi"]
