package gcsfuse

import (
	"fmt"

	"github.com/docker/go-plugins-helpers/volume"
)

type Driver struct {
}

func NewDriver() *Driver {
	return &Driver{}
}

func (d *Driver) Create(r *volume.CreateRequest) (retErr error) {
	return fmt.Errorf("Create() not implemented")
}

func (d *Driver) List() (*volume.ListResponse, error) {
	return nil, fmt.Errorf("List() not implemented")
}

func (d *Driver) Get(r *volume.GetRequest) (*volume.GetResponse, error) {
	return nil, fmt.Errorf("Get() not implemented")
}

func (d *Driver) Remove(r *volume.RemoveRequest) error {
	return fmt.Errorf("Remove() not implemented")
}

func (d *Driver) Path(r *volume.PathRequest) (*volume.PathResponse, error) {
	return nil, fmt.Errorf("Path() not implemented")
}

func (d *Driver) Mount(r *volume.MountRequest) (*volume.MountResponse, error) {
	return nil, fmt.Errorf("Mount() not implemented")
}

func (d *Driver) Unmount(r *volume.UnmountRequest) error {
	return fmt.Errorf("Unmount() not implemented")
}

func (d *Driver) Capabilities() *volume.CapabilitiesResponse {
	return &volume.CapabilitiesResponse{
		Capabilities: volume.Capability{
			Scope: "global", // TODO: what is the best scope here?
		},
	}
}
