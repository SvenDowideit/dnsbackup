module github.com/SvenDowideit/dnsbackup

go 1.14

replace github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.3.0

replace github.com/docker/docker => github.com/docker/engine v0.0.0-20190822205725-ed20165a37b4

replace github.com/moby/moby => github.com/docker/engine v0.0.0-20190822205725-ed20165a37b4

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.1.0

require (
	github.com/alecthomas/kong v0.2.11
	github.com/libdns/digitalocean v0.0.0-20200817185712-f11d70f2506c
	github.com/libdns/gandi v1.0.1
	github.com/libdns/libdns v0.1.0
	github.com/onaci/docker-ona v0.0.0-20200520071557-1181f04dd130
)
