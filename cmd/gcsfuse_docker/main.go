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

var (
	volumeRoot    = flag.String("volume_root", "", "Root path under which volumes are mounted")
	gcloudKeyPath = flag.String("gcloud_json_key_path", "", "Path to service key JSON to use as credentials")
)

func main() {
	flag.Parse()

	driver := gcsfuse.NewLogging(gcsfuse.NewDriver(*volumeRoot, *gcloudKeyPath))
	handler := volume.NewHandler(driver)
	if err := handler.ServeUnix(DriverName, 0); err != nil {
		glog.Exitf("Error serving %s.sock: %s", DriverName, err)
	}
}
