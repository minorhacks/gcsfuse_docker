package gcsfuse

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/docker/go-plugins-helpers/volume"
)

type Driver struct {
	// root is the path under which mounts are made
	root string

	// mu protects the vars in this group
	mu sync.Mutex
	// vols maps unique volume names to the Volume data
	vols map[string]*Volume
}

func NewDriver(volumeRoot string, keyPath string) *Driver {
	return &Driver{
		root: volumeRoot,

		vols: map[string]*Volume{},
	}
}

func (d *Driver) Create(r *volume.CreateRequest) (retErr error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Don't allow duplicates to be created; Docker calls Get() first on any
	// given name before resorting to Create(). If two different containers are
	// racing to create a volume with the same name, fail one of them - names
	// should be unique.
	if _, ok := d.vols[r.Name]; ok {
		return fmt.Errorf("volume %q already exists", r.Name)
	}

	// Ensure that all required options are set
	if err := checkRequiredCreateOptions(r.Options, []string{
		"bucket",
		"dir",
	}); err != nil {
		return fmt.Errorf("checking options for %q: %w", r.Name, err)
	}

	// Record volume parameters
	vol := &Volume{
		bucket:   r.Options["bucket"],
		hostPath: filepath.Join(d.root, r.Name),
	}
	d.vols[r.Name] = vol

	// TODO: Make directory for volume

	return nil
}

func (d *Driver) List() (*volume.ListResponse, error) {
	return nil, fmt.Errorf("List() not implemented")
}

func (d *Driver) Get(r *volume.GetRequest) (*volume.GetResponse, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	vol, ok := d.vols[r.Name]
	if !ok {
		return nil, fmt.Errorf("volume %q not found", r.Name)
	}

	return &volume.GetResponse{
		Volume: &volume.Volume{
			Name:       r.Name,
			Mountpoint: vol.hostPath,
		},
	}, nil
}

func (d *Driver) Remove(r *volume.RemoveRequest) error {
	return fmt.Errorf("Remove() not implemented")
}

func (d *Driver) Path(r *volume.PathRequest) (*volume.PathResponse, error) {
	return nil, fmt.Errorf("Path() not implemented")
}

func (d *Driver) Mount(r *volume.MountRequest) (*volume.MountResponse, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Get volume
	vol, ok := d.vols[r.Name]
	if !ok {
		return nil, fmt.Errorf("volume %q not found", r.Name)
	}

	// Mount volume
	if err := vol.Mount(); err != nil {
		return nil, fmt.Errorf("mount for %q failed: %w", r.Name, err)
	}

	// Return response
	return &volume.MountResponse{
		Mountpoint: vol.hostPath,
	}, nil
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
