package gcsfuse

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/golang/glog"
)

type Driver struct {
	// root is the path under which mounts are made
	root string
	// credsPath is the path to gcloud JSON credentials to pass on to gcsfuse.
	credsPath string

	// mu protects the vars in this group
	mu sync.Mutex
	// vols maps unique volume names to the Volume data
	vols map[string]*Volume
}

func NewDriver(volumeRoot string, credsPath string) *Driver {
	return &Driver{
		root:      volumeRoot,
		credsPath: credsPath,

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
		"subdir",
	}); err != nil {
		return fmt.Errorf("checking options for %q: %w", r.Name, err)
	}

	// Record volume parameters
	vol := &Volume{
		bucket:   r.Options["bucket"],
		subdir:   r.Options["subdir"],
		hostPath: filepath.Join(d.root, r.Name),
	}
	d.vols[r.Name] = vol

	// Make directory for volume
	if err := os.MkdirAll(vol.hostPath, 0o755); err != nil {
		return fmt.Errorf("failed to create local mountpoint %q: %w", vol.hostPath, err)
	}

	return nil
}

func (d *Driver) List() (*volume.ListResponse, error) {
	var ret []*volume.Volume
	for name, vol := range d.vols {
		ret = append(ret, &volume.Volume{
			Name:       name,
			Mountpoint: vol.hostPath,
		})
	}
	return &volume.ListResponse{
		Volumes: ret,
	}, nil
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
	d.mu.Lock()
	defer d.mu.Unlock()

	// Best-effort discovery of mount directory to delete - in case the driver
	// was restarted, and Docker is asking about a volume that no longer exists.
	var path string
	vol, ok := d.vols[r.Name]
	if !ok {
		glog.Warningf("Remove() called for unknown volume: %q", r.Name)
		path = filepath.Join(d.root, r.Name)
	} else {
		path = vol.hostPath
	}

	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("failed to remove mount path %q: %w", path, err)
	}

	delete(d.vols, r.Name)
	return nil
}

func (d *Driver) Path(r *volume.PathRequest) (*volume.PathResponse, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	vol, ok := d.vols[r.Name]
	if !ok {
		return nil, fmt.Errorf("volume %q not found", r.Name)
	}

	return &volume.PathResponse{
		Mountpoint: vol.hostPath,
	}, nil
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
	if err := vol.Mount(d.credsPath); err != nil {
		return nil, fmt.Errorf("mount for %q failed: %w", r.Name, err)
	}

	// Return response
	return &volume.MountResponse{
		Mountpoint: vol.hostPath,
	}, nil
}

func (d *Driver) Unmount(r *volume.UnmountRequest) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Get volume
	vol, ok := d.vols[r.Name]
	if !ok {
		return fmt.Errorf("volume %q not found", r.Name)
	}

	// Unmount
	if err := vol.Unmount(); err != nil {
		return fmt.Errorf("unmount for %q failed: %w", r.Name, err)
	}
	return nil
}

func (d *Driver) Capabilities() *volume.CapabilitiesResponse {
	return &volume.CapabilitiesResponse{
		Capabilities: volume.Capability{
			Scope: "global", // TODO: what is the best scope here?
		},
	}
}
