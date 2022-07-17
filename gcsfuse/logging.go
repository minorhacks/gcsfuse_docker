package gcsfuse

import (
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/golang/glog"
)

// Logging is a Docker volume driver that simply logs at the entry and exit of
// each function, delegating to an underlying Driver implementation:
// * Entry is logged if the logging level is set above the default to some
//   debugging threshold
// * Exit is logged if the logging level is set to a debugging threshold, or
//   whenever there is an error.
//
// This is not a complicated driver, nor does it obviate the need for a logger
// within the contained driver as well, but it does remove noisy code from the
// main logging driver.
type Logging struct {
	d volume.Driver
}

func NewLogging(d volume.Driver) *Logging {
	return &Logging{d: d}
}

func (d *Logging) Create(r *volume.CreateRequest) (retErr error) {
	glog.V(2).Infof("Create() called: %+v", r)
	defer func() {
		if retErr != nil {
			glog.Errorf("Create() failed: %v", retErr)
		} else {
			glog.V(2).Info("Create() success")
		}
	}()

	return d.d.Create(r)
}

func (d *Logging) List() (res *volume.ListResponse, retErr error) {
	glog.V(2).Info("List() called")
	defer func() {
		if retErr != nil {
			glog.Errorf("List() failed: %v", retErr)
		} else {
			glog.V(2).Infof("List() success: %+v", res)
		}
	}()

	return d.d.List()
}

func (d *Logging) Get(r *volume.GetRequest) (res *volume.GetResponse, retErr error) {
	glog.V(2).Infof("Get() called: %+v", r)
	defer func() {
		if retErr != nil {
			glog.Errorf("Get() failed: %v", retErr)
		} else {
			glog.V(2).Infof("Get() success: %+v", res)
		}
	}()

	return d.d.Get(r)
}

func (d *Logging) Remove(r *volume.RemoveRequest) (retErr error) {
	glog.V(2).Infof("Remove() called: %+v", r)
	defer func() {
		if retErr != nil {
			glog.Errorf("Remove() failed: %v", retErr)
		} else {
			glog.V(2).Info("Remove() success")
		}
	}()

	return d.d.Remove(r)
}

func (d *Logging) Path(r *volume.PathRequest) (res *volume.PathResponse, retErr error) {
	glog.V(2).Infof("Path() called: %+v", r)
	defer func() {
		if retErr != nil {
			glog.Errorf("Path() failed: %v", retErr)
		} else {
			glog.V(2).Infof("Path() success: %+v", res)
		}
	}()

	return d.d.Path(r)
}

func (d *Logging) Mount(r *volume.MountRequest) (res *volume.MountResponse, retErr error) {
	glog.V(2).Infof("Mount() called: %+v", r)
	defer func() {
		if retErr != nil {
			glog.Errorf("Mount() failed: %v", retErr)
		} else {
			glog.V(2).Infof("Mount() success: %+v", res)
		}
	}()

	return d.d.Mount(r)
}

func (d *Logging) Unmount(r *volume.UnmountRequest) (retErr error) {
	glog.V(2).Infof("Unmount() called: %+v", r)
	defer func() {
		if retErr != nil {
			glog.Errorf("Unmount() failed: %v", retErr)
		} else {
			glog.V(2).Info("Unmount() success")
		}
	}()

	return d.d.Unmount(r)
}

func (d *Logging) Capabilities() (res *volume.CapabilitiesResponse) {
	glog.V(2).Info("Capabilities() called")
	defer func() {
		glog.V(2).Infof("Capabilities(): %+v", res)
	}()

	return d.d.Capabilities()
}
