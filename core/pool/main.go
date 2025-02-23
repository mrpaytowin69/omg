package pool

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"opensvc.com/opensvc/core/driver"
	"opensvc.com/opensvc/core/keyop"
	"opensvc.com/opensvc/core/nodesinfo"
	"opensvc.com/opensvc/core/path"
	"opensvc.com/opensvc/core/rawconfig"
	"opensvc.com/opensvc/core/volaccess"
	"opensvc.com/opensvc/util/key"
	"opensvc.com/opensvc/util/render/tree"
	"opensvc.com/opensvc/util/san"
	"opensvc.com/opensvc/util/sizeconv"
)

type (
	T struct {
		driver string
		name   string
		config Config
	}

	StatusUsage struct {
		// Free unit is KiB
		Free float64 `json:"free"`
		// Used unit is KiB
		Used float64 `json:"used"`
		// Size unit is KiB
		Size float64 `json:"size"`
	}

	Status struct {
		Type         string         `json:"type"`
		Name         string         `json:"name"`
		Capabilities []string       `json:"capabilities"`
		Head         string         `json:"head"`
		Errors       []string       `json:"errors"`
		Volumes      []VolumeStatus `json:"volumes"`
		StatusUsage
	}
	StatusList   []Status
	Capabilities []string

	VolumeStatus struct {
		Path     path.T   `json:"path"`
		Children []path.T `json:"children"`
		Orphan   bool     `json:"orphan"`
		// Size unit is B
		Size float64 `json:"size"`
	}
	VolumeStatusList []VolumeStatus

	Config interface {
		GetString(key.T) string
		GetStringStrict(key.T) (string, error)
		GetStrings(key.T) []string
		GetBool(k key.T) bool
		GetSize(k key.T) *int64
	}
	Pooler interface {
		SetName(string)
		SetDriver(string)
		Name() string
		Type() string
		Head() string
		Mappings() map[string]string
		Capabilities() []string
		Usage() (StatusUsage, error)
		SetConfig(Config)
		Config() Config
		Separator() string
	}
	ArrayPooler interface {
		Pooler
		GetTargets() (san.Targets, error)
		CreateDisk(name string, size float64, paths san.Paths) ([]Disk, error)
		DeleteDisk(name string) ([]Disk, error)
	}
	Translater interface {
		Translate(name string, size float64, shared bool) ([]string, error)
	}
	BlkTranslater interface {
		BlkTranslate(name string, size float64, shared bool) ([]string, error)
	}
	volumer interface {
		FQDN() string
		Set(context.Context, ...keyop.T) error
	}

	Disk struct {
		// ID is the created disk wwid
		ID string

		// Paths is the subset of requested san path actually setup for this disk
		Paths san.Paths

		// Driver is a driver-specific dataset
		Driver any
	}
)

func NewStatus() Status {
	t := Status{}
	t.Volumes = make([]VolumeStatus, 0)
	t.Errors = make([]string, 0)
	return t
}

func sectionName(poolName string) string {
	return "pool#" + poolName
}

func cKey(poolName string, option string) key.T {
	section := sectionName(poolName)
	return key.New(section, option)
}

func cString(config Config, poolName string, option string) string {
	key := cKey(poolName, option)
	return config.GetString(key)
}

func New(name string, config Config) Pooler {
	poolType := cString(config, name, "type")
	fn := Driver(poolType)
	if fn == nil {
		return nil
	}
	t := fn()
	t.SetName(name)
	t.SetDriver(poolType)
	t.SetConfig(config)
	return t.(Pooler)
}

func (t *T) Mappings() map[string]string {
	s := cString(t.config, t.name, "mappings")
	m := make(map[string]string)
	for _, e := range strings.Fields(s) {
		l := strings.SplitN(e, ":", 2)
		if len(l) < 2 {
			continue
		}
		m[l[0]] = l[1]
	}
	return m
}

func Driver(t string) func() Pooler {
	did := driver.NewID(driver.GroupPool, t)
	i := driver.Get(did)
	if i == nil {
		return nil
	}
	if drv, ok := i.(func() Pooler); ok {
		return drv
	}
	return nil
}

// Separator is the string to use as the separator between
// name and hostname in the array-side disk name. Some array
// have a restricted characterset for such names, so better
// let the pool driver decide.
func (t T) Separator() string {
	return "-"
}

func (t T) Name() string {
	return t.name
}

func (t *T) SetName(name string) {
	t.name = name
}

func (t *T) SetDriver(driver string) {
	t.driver = driver
}

func (t T) Type() string {
	return t.driver
}

func (t *T) Config() Config {
	return t.config
}

func (t *T) SetConfig(c Config) {
	t.config = c
}

func GetStatus(t Pooler, withUsage bool) Status {
	data := NewStatus()
	data.Type = t.Type()
	data.Name = t.Name()
	data.Capabilities = t.Capabilities()
	data.Head = t.Head()
	if withUsage {
		if usage, err := t.Usage(); err != nil {
			data.Errors = append(data.Errors, err.Error())
		} else {
			data.Free = usage.Free
			data.Used = usage.Used
			data.Size = usage.Size
		}
	}
	return data
}

func pKey(p Pooler, s string) key.T {
	return pk(p.Name(), s)
}

func pk(poolName, s string) key.T {
	return key.New("pool#"+poolName, s)
}

func (t *T) GetStrings(s string) []string {
	k := pk(t.name, s)
	return t.Config().GetStrings(k)
}

func (t *T) GetString(s string) string {
	k := pk(t.name, s)
	return t.Config().GetString(k)
}

func (t *T) GetBool(s string) bool {
	k := pk(t.name, s)
	return t.Config().GetBool(k)
}

func (t *T) GetSize(s string) *int64 {
	k := pk(t.name, s)
	return t.Config().GetSize(k)
}

func (t *T) MkfsOptions() string {
	return t.GetString("mkfs_opt")
}

func (t *T) MkblkOptions() string {
	return t.GetString("mkblk_opt")
}

func (t *T) FSType() string {
	return t.GetString("fs_type")
}

func (t *T) MntOptions() string {
	return t.GetString("mnt_opt")
}

func (t *T) AddFS(name string, shared bool, fsIndex int, diskIndex int, onDisk string) []string {
	data := make([]string, 0)
	fsType := t.FSType()
	switch fsType {
	case "zfs":
		data = append(data, []string{
			fmt.Sprintf("disk#%d.type=zpool", diskIndex),
			fmt.Sprintf("disk#%d.name=%s", diskIndex, name),
			fmt.Sprintf("disk#%d.vdev={%s.exposed_devs[0]}", diskIndex, onDisk),
			fmt.Sprintf("disk#%d.shared=%t", diskIndex, shared),
			fmt.Sprintf("fs#%d.type=zfs", fsIndex),
			fmt.Sprintf("fs#%d.dev=%s/root", fsIndex, name),
			fmt.Sprintf("fs#%d.mnt=%s", fsIndex, MountPointFromName(name)),
			fmt.Sprintf("fs#%d.shared=%t", fsIndex, shared),
		}...)
	case "":
		panic("fsType should not be empty at this point")
	default:
		data = append(data, []string{
			fmt.Sprintf("fs#%d.type=%s", fsIndex, fsType),
			fmt.Sprintf("fs#%d.dev={%s.exposed_devs[0]}", fsIndex, onDisk),
			fmt.Sprintf("fs#%d.mnt=%s", fsIndex, MountPointFromName(name)),
			fmt.Sprintf("fs#%d.shared=%t", fsIndex, shared),
		}...)
	}
	if opts := t.MkfsOptions(); opts != "" {
		data = append(data, fmt.Sprintf("fs#%d.mkfs_opt=%s", fsIndex, opts))
	}
	if opts := t.MntOptions(); opts != "" {
		data = append(data, fmt.Sprintf("fs#%d.mnt_opt=%s", fsIndex, opts))
	}
	return data
}

func MountPointFromName(name string) string {
	return filepath.Join(filepath.FromSlash("/srv"), name)
}

func baseKeywords(p Pooler, size float64, acs volaccess.T) []string {
	return []string{
		fmt.Sprintf("pool=%s", p.Name()),
		fmt.Sprintf("size=%s", sizeconv.ExactBSizeCompact(size)),
		fmt.Sprintf("access=%s", acs),
	}
}

func flexKeywords(acs volaccess.T) []string {
	if acs.IsOnce() {
		return []string{}
	}
	return []string{
		"topology=flex",
		"flex_min=0",
	}
}

func nodeKeywords(nodes []string) []string {
	if len(nodes) <= 0 {
		return []string{}
	}
	return []string{
		"nodes=" + strings.Join(nodes, " "),
	}
}

func statusScheduleKeywords(p Pooler) []string {
	statusSchedule := p.Config().GetString(pKey(p, "status_schedule"))
	if statusSchedule == "" {
		return []string{}
	}
	return []string{
		"status_schedule=" + statusSchedule,
	}
}

func syncKeywords() []string {
	if true {
		return []string{}
	}
	return []string{
		"sync#i0.disable=true",
	}
}

func ConfigureVolume(p Pooler, vol volumer, size float64, format bool, acs volaccess.T, shared bool, nodes []string, env []string) error {
	name := vol.FQDN()
	kws, err := translate(p, name, size, format, shared)
	if err != nil {
		return err
	}
	kws = append(kws, env...)
	kws = append(kws, baseKeywords(p, size, acs)...)
	kws = append(kws, flexKeywords(acs)...)
	kws = append(kws, nodeKeywords(nodes)...)
	kws = append(kws, statusScheduleKeywords(p)...)
	kws = append(kws, syncKeywords()...)
	if err := vol.Set(context.Background(), keyop.ParseOps(kws)...); err != nil {
		return err
	}
	return nil
}

func translate(p Pooler, name string, size float64, format bool, shared bool) ([]string, error) {
	var kws []string
	var err error
	switch format {
	case true:
		o, ok := p.(Translater)
		if !ok {
			return nil, fmt.Errorf("pool %s does not support formatted volumes", p.Name())
		}
		if kws, err = o.Translate(name, size, shared); err != nil {
			return nil, err
		}
	case false:
		o, ok := p.(BlkTranslater)
		if !ok {
			return nil, fmt.Errorf("pool %s does not support block volumes", p.Name())
		}
		if kws, err = o.BlkTranslate(name, size, shared); err != nil {
			return nil, err
		}
	}
	return kws, nil
}

func NewStatusList() StatusList {
	l := make([]Status, 0)
	return StatusList(l)
}

func (t StatusList) Len() int {
	return len(t)
}

func (t StatusList) Less(i, j int) bool {
	return t[i].Name < t[j].Name
}

func (t StatusList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t StatusList) Add(p Pooler, withUsage bool) StatusList {
	s := GetStatus(p, withUsage)
	l := []Status(t)
	l = append(l, s)
	return StatusList(l)
}

func (t StatusList) Render(verbose bool) string {
	nt := t
	if !verbose {
		for i, _ := range nt {
			nt[i].Volumes = []VolumeStatus{}
		}
	}
	return t.Tree().Render()
}

// Tree returns a tree loaded with the type instance.
func (t StatusList) Tree() *tree.Tree {
	tree := tree.New()
	t.LoadTreeNode(tree.Head())
	return tree
}

// LoadTreeNode add the tree nodes representing the type instance into another.
func (t StatusList) LoadTreeNode(head *tree.Node) {
	head.AddColumn().AddText("name").SetColor(rawconfig.Color.Bold)
	head.AddColumn().AddText("type").SetColor(rawconfig.Color.Bold)
	head.AddColumn().AddText("caps").SetColor(rawconfig.Color.Bold)
	head.AddColumn().AddText("head").SetColor(rawconfig.Color.Bold)
	head.AddColumn().AddText("vols").SetColor(rawconfig.Color.Bold)
	head.AddColumn().AddText("size").SetColor(rawconfig.Color.Bold)
	head.AddColumn().AddText("used").SetColor(rawconfig.Color.Bold)
	head.AddColumn().AddText("free").SetColor(rawconfig.Color.Bold)
	sort.Sort(t)
	for _, data := range t {
		n := head.AddNode()
		data.LoadTreeNode(n)
	}
}

// LoadTreeNode add the tree nodes representing the type instance into another.
func (t Status) LoadTreeNode(head *tree.Node) {
	head.AddColumn().AddText(t.Name).SetColor(rawconfig.Color.Primary)
	head.AddColumn().AddText(t.Type)
	head.AddColumn().AddText(strings.Join(t.Capabilities, ","))
	head.AddColumn().AddText(t.Head)
	head.AddColumn().AddText(fmt.Sprint(len(t.Volumes)))
	if t.Size == 0 {
		head.AddColumn().AddText("-")
		head.AddColumn().AddText("-")
		head.AddColumn().AddText("-")
	} else {
		head.AddColumn().AddText(sizeconv.BSizeCompact(t.Size * sizeconv.KiB))
		head.AddColumn().AddText(sizeconv.BSizeCompact(t.Used * sizeconv.KiB))
		head.AddColumn().AddText(sizeconv.BSizeCompact(t.Free * sizeconv.KiB))
	}
	if len(t.Volumes) > 0 {
		n := head.AddNode()
		VolumeStatusList(t.Volumes).LoadTreeNode(n)
	}
}

func (t VolumeStatusList) Len() int {
	return len(t)
}

func (t VolumeStatusList) Less(i, j int) bool {
	return t[i].Path.String() < t[j].Path.String()
}

func (t VolumeStatusList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// LoadTreeNode add the tree nodes representing the type instance into another.
func (t VolumeStatusList) LoadTreeNode(head *tree.Node) {
	head.AddColumn().AddText("volume").SetColor(rawconfig.Color.Bold)
	head.AddColumn().AddText("").SetColor(rawconfig.Color.Bold)
	head.AddColumn().AddText("children").SetColor(rawconfig.Color.Bold)
	head.AddColumn().AddText("orphan").SetColor(rawconfig.Color.Bold)
	head.AddColumn().AddText("").SetColor(rawconfig.Color.Bold)
	head.AddColumn().AddText("").SetColor(rawconfig.Color.Bold)
	head.AddColumn().AddText("").SetColor(rawconfig.Color.Bold)
	head.AddColumn().AddText("").SetColor(rawconfig.Color.Bold)
	sort.Sort(t)
	for _, data := range t {
		n := head.AddNode()
		data.LoadTreeNode(n)
	}
}

// LoadTreeNode add the tree nodes representing the type instance into another.
func (t VolumeStatus) LoadTreeNode(head *tree.Node) {
	head.AddColumn().AddText(t.Path.String())
	head.AddColumn().AddText("")
	head.AddColumn().AddText(path.L(t.Children).String())
	head.AddColumn().AddText(strconv.FormatBool(t.Orphan))
	head.AddColumn().AddText("")
	head.AddColumn().AddText(sizeconv.BSizeCompact(t.Size))
	head.AddColumn().AddText("")
	head.AddColumn().AddText("")
}

func HasAccess(p Pooler, acs volaccess.T) bool {
	return HasCapability(p, acs.String())
}

func HasCapability(p Pooler, s string) bool {
	for _, capa := range p.Capabilities() {
		if capa == s {
			return true
		}
	}
	return false

}

func (t Status) HasAccess(acs volaccess.T) bool {
	return t.HasCapability(acs.String())
}

func (t Status) HasCapability(s string) bool {
	for _, capa := range t.Capabilities {
		if capa == s {
			return true
		}
	}
	return false

}

func GetMapping(p ArrayPooler, nodes []string) (san.Paths, error) {
	targets, err := p.GetTargets()
	if err != nil {
		return san.Paths{}, err
	}
	nodesInfo, err := nodesinfo.Get()
	if err != nil {
		return san.Paths{}, err
	}
	filteredPaths := make(san.Paths, 0)
	for _, node := range nodes {
		nodeInfo, ok := nodesInfo[node]
		if !ok {
			continue
		}
		for _, target := range targets {
			filteredPaths = append(filteredPaths, nodeInfo.Paths.WithTargetName(target.Name)...)
		}
	}
	return filteredPaths, nil
}
