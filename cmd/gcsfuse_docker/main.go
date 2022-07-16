package main

import (
	"flag"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/golang/glog"
	"github.com/minorhacks/gcsfuse_docker/gcsfuse"
)

const (
	DriverName = "gcsfuse"
)

func main() {
	flag.Parse()

	driver := gcsfuse.NewDriver()
	handler := volume.NewHandler(driver)
	if err := handler.ServeUnix(DriverName, 0); err != nil {
		glog.Exitf("Error serving %s.sock: %s", DriverName, err)
	}
}
