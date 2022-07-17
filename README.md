# gcsfuse Driver for Docker

This repo implements a Docker driver for gcsfuse. There are many like it, but
this one is mine.

This plugin focuses on:

* Facilitating the creation of many volumes that use (different subdirs of) the
  same GCS bucket. [Some existing Docker
  drivers](https://github.com/lorenzleutgeb/docker-volume-gcs/blob/45155d931d8fe4b0492f440257c87e829a1cb537/main.go#L68)
  seem to map `gcsfuse` processes to buckets; if trying to mount a subdir of the
  same bucket, multiple mounts would use the same `gcsfuse` process, which might
  cause its caching to thrash with enough load.
* Enough logging/monitoring to aid debugging in a production setting

## Usage

TODO: This section is a work-in-progress.

1. Build and run `gcsfuse_docker`:

   ```
   go build ./cmd/gcsfuse_docker
   sudo ./gcsfuse_docker --v=1 --alsologtostderr
   ```

1. Test a mount with a Docker container:

   ```
   docker run \
     -it \
     --mount "type=volume,src=foobar,dst=/gcsfuse_test,volume-driver=gcsfuse,volume-opt=opt1=foo,volume-opt=opt2=bar" \
     ubuntu \
     /bin/bash
   ```

   TODO: Clarify the parameters in the mount string