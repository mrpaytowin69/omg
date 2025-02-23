//go:build linux || solaris

package poolfreenas

import (
	"errors"
	"fmt"

	"opensvc.com/opensvc/core/driver"
	"opensvc.com/opensvc/core/pool"
	"opensvc.com/opensvc/core/xconfig"
	"opensvc.com/opensvc/drivers/arrayfreenas"
	"opensvc.com/opensvc/util/san"
	"opensvc.com/opensvc/util/sizeconv"
)

type (
	T struct {
		pool.T
	}
)

var (
	drvID = driver.NewID(driver.GroupPool, "freenas")
)

func init() {
	driver.Register(drvID, NewPooler)
}

func NewPooler() pool.Pooler {
	t := New()
	var i interface{} = t
	return i.(pool.Pooler)
}

func New() *T {
	t := T{}
	return &t
}

func (t T) Head() string {
	return fmt.Sprintf("array://%s/%s", t.arrayName(), t.diskgroup())
}

func (t T) diskgroup() string {
	return t.GetString("diskgroup")
}

func (t T) insecureTPC() bool {
	return t.GetBool("insecureTPC")
}

func (t T) compression() bool {
	return t.GetBool("compression")
}

func (t T) sparse() bool {
	return t.GetBool("sparse")
}

func (t T) blocksize() *int64 {
	return t.GetSize("blocksize")
}

func (t T) arrayName() string {
	return t.GetString("array")
}

func (t T) Capabilities() []string {
	return []string{"rox", "rwx", "roo", "rwo", "blk", "iscsi", "shared"}
}

func (t T) Usage() (pool.StatusUsage, error) {
	usage := pool.StatusUsage{}
	a := t.array()
	data, err := a.GetDataset(t.diskgroup())
	if err != nil {
		return usage, err
	}
	if i, err := sizeconv.FromSize(data.Used.Rawvalue); err != nil {
		return usage, err
	} else {
		usage.Used = float64(i / 1024)
	}
	if i, err := sizeconv.FromSize(data.Available.Rawvalue); err != nil {
		return usage, err
	} else {
		usage.Free = float64(i / 1024)
	}
	usage.Size = usage.Used + usage.Free
	return usage, nil
}

func (t T) array() *arrayfreenas.Array {
	a := arrayfreenas.New()
	a.SetName(t.arrayName())
	a.SetConfig(t.Config().(*xconfig.T))
	return a
}

func (t *T) Translate(name string, size float64, shared bool) ([]string, error) {
	data, err := t.BlkTranslate(name, size, shared)
	if err != nil {
		return nil, err
	}
	data = append(data, t.AddFS(name, shared, 1, 0, "disk#0")...)
	return data, nil
}

func (t *T) BlkTranslate(name string, size float64, shared bool) ([]string, error) {
	data := []string{
		"disk#0.type=disk",
		"disk#0.name=" + name,
		"disk#0.scsireserv=true",
		"shared=" + fmt.Sprint(shared),
		"size=" + sizeconv.ExactBSizeCompact(size),
	}
	return data, nil
}

func (t *T) GetTargets() (san.Targets, error) {
	a := t.array()
	data, err := a.GetISCSITargets()
	if err != nil {
		return nil, err
	}
	ports := make(san.Targets, 0)
	for _, d := range data {
		ports = append(ports, san.Target{
			Name: d.Name,
			Type: san.ISCSI,
		})
	}
	return ports, nil
}

func (t *T) DeleteDisk(name string) ([]pool.Disk, error) {
	disk := pool.Disk{}
	a := t.array()
	drvName := t.diskgroup() + "/" + name
	drvDisk, err := a.DelDisk(drvName)
	if err != nil {
		return []pool.Disk{}, err
	}
	disk.Driver = drvDisk
	disk.ID = a.DiskID(*drvDisk)
	if paths, err := a.DiskPaths(*drvDisk); err != nil {
		return []pool.Disk{disk}, err
	} else {
		disk.Paths = paths
	}
	return []pool.Disk{disk}, nil
}

func (t *T) CreateDisk(name string, size float64, paths san.Paths) ([]pool.Disk, error) {
	disk := pool.Disk{}
	if len(paths) == 0 {
		return []pool.Disk{}, errors.New("no mapping in request. cowardly refuse to create a disk that can not be mapped")
	}
	a := t.array()
	blocksize := fmt.Sprint(*t.blocksize())
	sparse := t.sparse()
	insecureTPC := t.insecureTPC()
	drvSize := sizeconv.ExactBSizeCompact(size)
	drvName := t.diskgroup() + "/" + name
	mapping := paths.Mapping()

	drvDisk, err := a.AddDisk(drvName, drvSize, blocksize, sparse, insecureTPC, mapping, nil)
	if err != nil {
		return []pool.Disk{}, err
	}
	disk.Driver = drvDisk
	disk.ID = a.DiskID(*drvDisk)
	if paths, err := a.DiskPaths(*drvDisk); err != nil {
		return []pool.Disk{disk}, err
	} else {
		disk.Paths = paths
	}
	return []pool.Disk{disk}, nil
}
