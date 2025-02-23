package driver

type (
	Registry map[ID]any
)

var (
	registry = NewRegistry()
)

func NewRegistry() Registry {
	return make(Registry)
}

func Register(id ID, allocator any) {
	registry[id] = allocator
}

func Exists(id ID) bool {
	return Get(id) != nil
}

func Get(id ID) any {
	allocator, ok := registry[id]
	if !ok {
		// <group>.<name> driver not found, ... try <group>
		// used for example by the volume driver, whose
		// type keyword is not pointing a resource sub driver
		// but a pool driver.
		id.Name = ""
		return GetStrict(id)
	}
	return allocator
}

func GetStrict(id ID) any {
	allocator, _ := registry[id]
	return allocator
}

func List() IDs {
	l := make(IDs, len(registry))
	i := 0
	for did, _ := range registry {
		l[i] = did
		i = i + 1
	}
	return l
}

func ListWithGroup(group Group) Registry {
	m := NewRegistry()
	for _, did := range List() {
		if did.Group != group {
			continue
		}
		m[did] = Get(did)
	}
	return m
}
