package object

import (
	"opensvc.com/opensvc/core/driver"
	"opensvc.com/opensvc/core/keyop"
	"opensvc.com/opensvc/core/keywords"
	"opensvc.com/opensvc/core/kind"
	"opensvc.com/opensvc/core/placement"
	"opensvc.com/opensvc/core/resource"
	"opensvc.com/opensvc/core/resourceid"
	"opensvc.com/opensvc/core/xconfig"
	"opensvc.com/opensvc/util/converters"
	"opensvc.com/opensvc/util/key"
)

var keywordStore = keywords.Store{
	{
		Section:   "DEFAULT",
		Option:    "hard_affinity",
		Inherit:   keywords.InheritHead,
		Converter: converters.ListLowercase,
		Aliases:   []string{"affinity"},
		Text:      "A whitespace separated list of services that must be started on the node to allow the monitor to start this service.",
		Example:   "svc1 svc2",
	},
	{
		Section:   "DEFAULT",
		Option:    "hard_anti_affinity",
		Inherit:   keywords.InheritHead,
		Converter: converters.ListLowercase,
		Aliases:   []string{"anti_affinity"},
		Text:      "A whitespace separated list of services that must not be started on the node to allow the monitor to start this service.",
		Example:   "svc1 svc2",
	},
	{
		Section:   "DEFAULT",
		Option:    "soft_affinity",
		Inherit:   keywords.InheritHead,
		Converter: converters.ListLowercase,
		Text:      "A whitespace separated list of services that must be started on the node to allow the monitor to start this service. If the local node is the only candidate ignore this constraint and allow start.",
		Example:   "svc1 svc2",
	},
	{
		Section:   "DEFAULT",
		Option:    "soft_anti_affinity",
		Inherit:   keywords.InheritHead,
		Converter: converters.ListLowercase,
		Text:      "A whitespace separated list of services that must not be started on the node to allow the monitor to start this service. If the local node is the only candidate ignore this constraint and allow start.",
		Example:   "svc1 svc2",
	},
	{
		Section:     "DEFAULT",
		Option:      "id",
		DefaultText: "<random uuid>",
		Scopable:    false,
		Text:        "A RFC 4122 random uuid generated by the agent. To use as reference in resources definitions instead of the service name, so the service can be renamed without affecting the resources.",
	},
	{
		Option: "comment",
		Text:   "Helps users understand the role of the service and resources, which is nice to on-call support people having to operate on a service they are not usually responsible for.",
	},
	{
		Section:   "DEFAULT",
		Option:    "disable",
		Scopable:  true,
		Converter: converters.Bool,
		Text:      "A disabled resource will be ignored on service startup and shutdown. Its status is always reported ``n/a``.\n\nSet in DEFAULT, the whole service is disabled. A disabled service does not honor start and stop actions. These actions immediately return success.\n\n:cmd:`om <path> disable` only sets :kw:`DEFAULT.disable`. As resources disabled state is not changed, :cmd:`om <path> enable` does not enable disabled resources.",
	},
	{
		Section:   "DEFAULT",
		Option:    "create_pg",
		Default:   "true",
		Scopable:  true,
		Converter: converters.Bool,
		Text:      "Use process containers when possible. Containers allow capping memory, swap and cpu usage per service. Lxc containers are naturally containerized, so skip containerization of their startapp.",
	},
	{
		Option:   "pg_cpus",
		Attr:     "PG.Cpus",
		Scopable: true,
		Inherit:  keywords.InheritLeaf,
		Depends:  keyop.ParseList("create_pg=true"),
		Text:     "Allow service process to bind only the specified cpus. Cpus are specified as list or range : 0,1,2 or 0-2",
		Example:  "0-2",
	},
	{
		Option:   "pg_mems",
		Attr:     "PG.Mems",
		Scopable: true,
		Inherit:  keywords.InheritLeaf,
		Text:     "Allow service process to bind only the specified memory nodes. Memory nodes are specified as list or range : 0,1,2 or 0-2",
		Example:  "0-2",
	},
	{
		Option:    "pg_cpu_shares",
		Attr:      "PG.CpuShares",
		Scopable:  true,
		Converter: converters.Size,
		Inherit:   keywords.InheritLeaf,
		Text:      "Kernel default value is used, which usually is 1024 shares. In a cpu-bound situation, ensure the service does not use more than its share of cpu ressource. The actual percentile depends on shares allowed to other services.",
		Example:   "512",
	},
	{
		Option:   "pg_cpu_quota",
		Attr:     "PG.CpuQuota",
		Scopable: true,
		Inherit:  keywords.InheritLeaf,
		Text:     "The percent ratio of one core to allocate to the process group if % is specified, else the absolute value to set in the process group parameter. For example, on Linux cgroups, ``-1`` means unlimited, and a positive absolute value means the number of microseconds to allocate each period. ``50%@all`` means 50% of all cores, and ``50%@2`` means 50% of two cores.",
		Example:  "50%@all",
	},
	{
		Option:   "pg_mem_oom_control",
		Attr:     "PG.MemOOMControl",
		Scopable: true,
		Inherit:  keywords.InheritLeaf,
		Text:     "A flag (0 or 1) that enables or disables the Out of Memory killer for a cgroup. If enabled (0), tasks that attempt to consume more memory than they are allowed are immediately killed by the OOM killer. The OOM killer is enabled by default in every cgroup using the memory subsystem; to disable it, write 1.",
		Example:  "1",
	},
	{
		Option:    "pg_mem_limit",
		Attr:      "PG.MemLimit",
		Scopable:  true,
		Converter: converters.Size,
		Inherit:   keywords.InheritLeaf,
		Text:      "Ensures the service does not use more than specified memory (in bytes). The Out-Of-Memory killer get triggered in case of tresspassing.",
		Example:   "512m",
	},
	{
		Option:    "pg_vmem_limit",
		Attr:      "PG.VMemLimit",
		Scopable:  true,
		Converter: converters.Size,
		Inherit:   keywords.InheritLeaf,
		Text:      "Ensures the service does not use more than specified memory+swap (in bytes). The Out-Of-Memory killer get triggered in case of tresspassing. The specified value must be greater than :kw:`pg_mem_limit`.",
		Example:   "1g",
	},
	{
		Option:   "pg_mem_swappiness",
		Attr:     "PG.MemSwappiness",
		Scopable: true,
		Inherit:  keywords.InheritLeaf,
		Text:     "Set a swappiness percentile value for the process group.",
		Example:  "40",
	},
	{
		Option:   "pg_blkio_weight",
		Attr:     "PG.BlkioWeight",
		Scopable: true,
		Inherit:  keywords.InheritLeaf,
		Text:     "Block IO relative weight. Value: between 10 and 1000. Kernel default: 1000.",
		Example:  "50",
	},
	{
		Option:    "stat_timeout",
		Converter: converters.Duration,
		Scopable:  true,
		Text:      "The maximum wait time for a stat call to respond. When expired, the resource status is degraded is to warn, which might cause a TOC if the resource is monitored.",
	},
	{
		Section:     "DEFAULT",
		Option:      "nodes",
		Scopable:    true,
		Kind:        kind.Or(kind.Svc, kind.Vol),
		Inherit:     keywords.InheritHead,
		Converter:   xconfig.NodesConverter,
		Text:        "A node selector expression specifying the list of cluster nodes hosting service instances.",
		DefaultText: "The lowercased hostname of the evaluating node.",
		Example:     "n1 n*",
	},
	{
		Section:   "DEFAULT",
		Option:    "nodes",
		Scopable:  true,
		Kind:      kind.Or(kind.Cfg, kind.Sec, kind.Usr, kind.Nscfg),
		Inherit:   keywords.InheritHead,
		Converter: xconfig.NodesConverter,
		Text:      "A node selector expression specifying the list of cluster nodes hosting service instances.",
		Default:   "*",
	},
	{
		Section:   "DEFAULT",
		Option:    "drpnodes",
		Scopable:  true,
		Inherit:   keywords.InheritHead,
		Converter: xconfig.OtherNodesConverter,
		Text:      "The backup node where the service is activated in a DRP situation. This node is also a data synchronization target for :c-res:`sync` resources.",
		Example:   "n1 n2",
	},
	{
		Section:   "DEFAULT",
		Option:    "encapnodes",
		Inherit:   keywords.InheritHead,
		Converter: xconfig.OtherNodesConverter,
		Text:      "The list of `containers` handled by this service and with an OpenSVC agent installed to handle the encapsulated resources. With this parameter set, parameters can be scoped with the ``@encapnodes`` suffix.",
		Example:   "n1 n2",
	},
	{
		Section:    "DEFAULT",
		Option:     "monitor_action",
		Inherit:    keywords.InheritHead,
		Scopable:   true,
		Candidates: []string{"reboot", "crash", "freezestop", "switch"},
		Text:       "The action to take when a monitored resource is not up nor standby up, and if the resource restart procedure has failed.",
		Example:    "reboot",
	},
	{
		Section:  "DEFAULT",
		Option:   "pre_monitor_action",
		Inherit:  keywords.InheritHead,
		Scopable: true,
		Text:     "A script to execute before the :kw:`monitor_action`. For example, if the :kw:`monitor_action` is set to ``freezestop``, the script can decide to crash the server if it detects a situation were the freezestop can not succeed (ex. fs can not be umounted with a dead storage array).",
		Example:  "/bin/true",
	},

	{
		Section: "DEFAULT",
		Option:  "app",
		Default: "default",
		Text:    "Used to identify who is responsible for this service, who is billable and provides a most useful filtering key. Better keep it a short code.",
	},
	{
		Section:     "DEFAULT",
		Option:      "env",
		Aliases:     []string{"service_type"},
		Inherit:     keywords.InheritHead,
		DefaultText: "Same as the node env",
		Candidates:  validEnvs,
		Text:        "A non-PRD service can not be brought up on a PRD node, but a PRD service can be startup on a non-PRD node (in a DRP situation). The default value is the node :kw:`env`.",
	},
	{
		Section:   "DEFAULT",
		Option:    "stonith",
		Inherit:   keywords.InheritHead,
		Converter: converters.Bool,
		Default:   "false",
		Depends:   keyop.ParseList("topology=failover"),
		Text:      "Stonith the node previously running the service if stale upon start by the daemon monitor.",
	},
	{
		Section:    "DEFAULT",
		Option:     "constraints",
		Inherit:    keywords.InheritHead,
		Scopable:   true,
		Deprecated: "2.1",
		Example:    "$(\"{nodename}\"==\"n2.opensvc.com\")",
		Depends:    keyop.ParseList("orchestrate=ha"),
		Text:       "An expression evaluating as a boolean, constraining the service instance placement by the daemon monitor to nodes with the constraints evaluated as True.\n\nThe constraints are not honored by manual start operations. The constraints value is embedded in the json status.\n\nSupported comparison operators are ``==``, ``!=``, ``>``, ``>=``, ``<=``, ``in (e1, e2)``, ``in [e1, e2]``.\n\nSupported arithmetic operators are ``*``, ``+``, ``-``, ``/``, ``**``, ``//``, ``%``.\n\nSupported binary operators are ``&``, ``|``, ``^``.\n\nThe negation operator is ``not``.\n\nSupported boolean operators are ``and``, ``or``.\n\nReferences are allowed.\n\nStrings, and references evaluating as strings, containing dots must be quoted.",
	},
	{
		Section:    "DEFAULT",
		Option:     "placement",
		Scopable:   false,
		Inherit:    keywords.InheritHead,
		Default:    "nodes order",
		Candidates: placement.PolicyNames(),
		Text: `Set a service instances placement policy:

* none        no placement policy. a policy for dummy, observe-only, services.
* nodes order the left-most available node is allowed to start a service instance when necessary.
* load avg    the least loaded node takes precedences.
* shift       shift the nodes order ranking by the service prefix converter to an integer.
* spread      a spread policy tends to perfect leveling with many services.
* score       the highest scoring node takes precedence (the score is a composite indice of load, mem and swap).
`,
	},
	{
		Section:    "DEFAULT",
		Option:     "topology",
		Scopable:   false,
		Default:    "failover",
		Inherit:    keywords.InheritHead,
		Aliases:    []string{"cluster_type"},
		Candidates: []string{"failover", "flex"},
		Text:       "``failover`` the service is allowed to be up on one node at a time. ``flex`` the service can be up on :kw:`flex_target` nodes, where :kw:`flex_target` must be in the [flex_min, flex_max] range.",
	},
	{
		Section:     "DEFAULT",
		Option:      "flex_primary",
		Scopable:    true,
		Inherit:     keywords.InheritHead,
		Converter:   converters.ListLowercase,
		Depends:     keyop.ParseList("topology=flex"),
		DefaultText: "first node of the nodes parameter.",
		Text:        "The node in charge of syncing the other nodes. :opt:`--cluster` actions on the flex_primary are executed on all peer nodes (ie, not drpnodes).",
	},
	{
		Section:   "DEFAULT",
		Option:    "shared",
		Scopable:  true,
		Default:   "true",
		Converter: converters.Bool,
		Text:      "Set to ``true`` to skip the resource on provision and unprovision actions if the action has already been done by a peer. Shared resources, like vg built on SAN disks must be provisioned once. All resources depending on a shared resource must also be flagged as shared.",
	},
	{
		Section:   "DEFAULT",
		Option:    "check_carrier",
		Scopable:  true,
		Default:   "true",
		Converter: converters.Bool,
		Text:      "Activate the link carrier check. Set to false if ipdev is a backend bridge or switch.",
	},
	{
		Section:   "DEFAULT",
		Option:    "flex_min",
		Aliases:   []string{"flex_min_nodes"},
		Default:   "1",
		Inherit:   keywords.InheritHead,
		Converter: converters.Int,
		Depends:   keyop.ParseList("topology=flex"),
		Text:      "Minimum number of up instances in the cluster. Below this number the aggregated service status is degraded to warn..",
	},
	{
		Section:     "DEFAULT",
		Option:      "flex_max",
		Aliases:     []string{"flex_max_nodes"},
		DefaultText: "Number of svc nodes",
		Inherit:     keywords.InheritHead,
		Converter:   converters.Int,
		Depends:     keyop.ParseList("topology=flex"),
		Text:        "Maximum number of up instances in the cluster. Above this number the aggregated service status is degraded to warn. ``0`` means unlimited.",
	},
	{
		Section:     "DEFAULT",
		Option:      "flex_target",
		DefaultText: "The value of flex_min",
		Inherit:     keywords.InheritHead,
		Converter:   converters.Int,
		Depends:     keyop.ParseList("topology=flex"),
		Text:        "Optimal number of up instances in the cluster. The value must be between :kw:`flex_min` and :kw:`flex_max`. If ``orchestrate=ha``, the monitor ensures the :kw:`flex_target` is met.",
	},
	{
		Section:   "DEFAULT",
		Option:    "parents",
		Inherit:   keywords.InheritHead,
		Converter: converters.ListLowercase,
		Text:      "List of services or instances expressed as ``<path>[@<nodename>]`` that must be ``avail up`` before allowing this service to be started by the daemon monitor. Whitespace separated.",
	},
	{
		Section:   "DEFAULT",
		Option:    "children",
		Default:   "",
		Inherit:   keywords.InheritHead,
		Converter: converters.ListLowercase,
		Text:      "List of services that must be ``avail down`` before allowing this service to be stopped by the daemon monitor. Whitespace separated.",
	},
	{
		Section:   "DEFAULT",
		Option:    "slaves",
		Default:   "",
		Inherit:   keywords.InheritHead,
		Converter: converters.ListLowercase,
		Text:      "List of services to propagate the :c-action:`start` and :c-action:`stop` actions to.",
	},
	{
		Section:    "DEFAULT",
		Option:     "orchestrate",
		Inherit:    keywords.InheritHead,
		Default:    "no",
		Candidates: []string{"no", "ha", "start"},
		Text:       "If set to ``no``, disable service orchestration by the OpenSVC daemon monitor, including service start on boot. If set to ``start`` failover services won't failover automatically, though the service instance on the natural placement leader is started if another instance is not already up. Flex services won't restart the :kw:`flex_target` number of up instances. Resource restart is still active whatever the :kw:`orchestrate` value.",
	},
	{
		Section:   "DEFAULT",
		Option:    "priority",
		Default:   "50",
		Scopable:  false,
		Inherit:   keywords.InheritHead,
		Converter: converters.Int,
		Text:      "A scheduling priority (the smaller the more priority) used by the monitor thread to trigger actions for the top priority services, so that the :kw:`node.max_parallel` constraint doesn't prevent prior services to start first. The priority setting is dropped from a service configuration injected via the api by a user not granted the prioritizer role.",
	},
	{
		Section:   "subset",
		Option:    "parallel",
		Scopable:  true,
		Converter: converters.Bool,
		Text:      "If set to ``true``, actions are executed in parallel amongst the subset member resources.",
	},

	// Secrets
	{
		Section:  "DEFAULT",
		Option:   "cn",
		Scopable: true,
		Text:     "Certificate Signing Request Common Name.",
		Example:  "test.opensvc.com",
		Kind:     kind.Or(kind.Sec),
	},
	{
		Section:  "DEFAULT",
		Option:   "c",
		Scopable: true,
		Text:     "Certificate Signing Request Country.",
		Example:  "FR",
		Kind:     kind.Or(kind.Sec),
	},
	{
		Section:  "DEFAULT",
		Option:   "st",
		Scopable: true,
		Text:     "Certificate Signing Request State.",
		Example:  "Oise",
		Kind:     kind.Or(kind.Sec),
	},
	{
		Section:  "DEFAULT",
		Option:   "l",
		Scopable: true,
		Text:     "Certificate Signing Request Location.",
		Example:  "Gouvieux",
		Kind:     kind.Or(kind.Sec),
	},
	{
		Section:  "DEFAULT",
		Option:   "o",
		Scopable: true,
		Text:     "Certificate Signing Request Organization.",
		Example:  "OpenSVC",
		Kind:     kind.Or(kind.Sec),
	},
	{
		Section:  "DEFAULT",
		Option:   "ou",
		Scopable: true,
		Text:     "Certificate Signing Request Organizational Unit.",
		Example:  "Lab",
		Kind:     kind.Or(kind.Sec),
	},
	{
		Section:  "DEFAULT",
		Option:   "email",
		Scopable: true,
		Text:     "Certificate Signing Request Email.",
		Example:  "test@opensvc.com",
		Kind:     kind.Or(kind.Sec),
	},
	{
		Section:   "DEFAULT",
		Option:    "alt_names",
		Converter: converters.List,
		Scopable:  true,
		Text:      "Certificate Signing Request Alternative Domain Names.",
		Example:   "www.opensvc.com opensvc.com",
		Kind:      kind.Or(kind.Sec),
	},
	{
		Section:   "DEFAULT",
		Option:    "bits",
		Converter: converters.Size,
		Scopable:  true,
		Text:      "Certificate Private Key Length.",
		Default:   "4kib",
		Example:   "8192",
		Kind:      kind.Or(kind.Sec),
	},

	// Usr
	{
		Section:   "DEFAULT",
		Option:    "grant",
		Scopable:  true,
		Kind:      kind.Or(kind.Usr),
		Inherit:   keywords.InheritHead,
		Converter: converters.ListLowercase,
		Text: `Grant roles on namespaces to the user.

A whitespace-separated list of root|squatter|prioritizer|blacklistadmin|<role selector>:<namespace selector>, 
where role selector is a comma-separated list of role in admin,operator,guest 
and the namespace selector is a glob pattern applied to existing namespaces.

The root role is required to add resource triggers and non-containerized resources other 
than (container.docker, container.podman task.docker, task.podman and volume). 
The squatter role is required to create a new namespace. 
The admin role is required to create, deploy and delete objects. 
The guest role is required to list and read objects configurations and status.
`,
		Example: "admin:test* guest:*",
	},

	{
		Section:   "DEFAULT",
		Option:    "rollback",
		Scopable:  true,
		Default:   "true",
		Converter: converters.Bool,
		Text:      "If set to ``false``, the default 'rollback on action error' behaviour is inhibited, leaving the service in its half-started state. The daemon also refuses to takeover a service if rollback is disabled and a peer instance is 'start failed'.",
	},
	{
		Section:   "DEFAULT",
		Option:    "validity",
		Converter: converters.Duration,
		Scopable:  true,
		Text:      "Certificate Validity duration.",
		Default:   "1y",
		Example:   "10y",
		Kind:      kind.Or(kind.Sec),
	},
	{
		Section:  "DEFAULT",
		Option:   "ca",
		Scopable: true,
		Text:     "The name of secret containing a certificate to use as a Certificate Authority. This secret must be in the same namespace.",
		Example:  "ca",
		Kind:     kind.Or(kind.Sec),
	},
	{
		Section:  "DEFAULT",
		Option:   "monitor_schedule",
		Scopable: true,
		Text:     "The object's monitored resources status evaluation schedule. See ``usr/share/doc/schedule`` for the schedule syntax.",
		Default:  "@5m",
	},
	{
		Section:  "DEFAULT",
		Option:   "resinfo_schedule",
		Scopable: true,
		Text:     "The object's key-val table emit schedule. See ``usr/share/doc/schedule`` for the schedule syntax.",
		Default:  "@60m",
	},
	{
		Section:  "DEFAULT",
		Option:   "status_schedule",
		Scopable: true,
		Text:     "The object's status evaluation schedule. See ``usr/share/doc/schedule`` for the schedule syntax.",
		Default:  "@10m",
	},
	{
		Section:  "DEFAULT",
		Option:   "comp_schedule",
		Scopable: true,
		Text:     "The object's compliance run schedule. See ``usr/share/doc/schedule`` for the schedule syntax.",
		Default:  "~00:00-06:00",
	},
	{
		Section:  "DEFAULT",
		Option:   "sync_schedule",
		Scopable: true,
		Text:     "The object's sync default schedule. See ``usr/share/doc/schedule`` for the schedule syntax.",
		Default:  "04:00-06:00",
	},
	{
		Section:  "DEFAULT",
		Option:   "run_schedule",
		Scopable: true,
		Text:     "The object's tasks run action default schedule. See ``usr/share/doc/schedule`` for the schedule syntax.",
	},
	{
		Option:    "timeout",
		Attr:      "Timeout",
		Converter: converters.Duration,
		Scopable:  true,
		Default:   "1h",
		Text:      "Wait for <duration> before declaring the any state-changing action a failure. A per-action <action>_timeout can override this value.",
		Example:   "2h",
	},
	{
		Option:    "start_timeout",
		Attr:      "StartTimeout",
		Converter: converters.Duration,
		Scopable:  true,
		Text:      "Wait for <duration> before declaring the action a failure. Takes precedence over :kw:`timeout`.",
		Example:   "1m30s",
	},
	{
		Option:    "stop_timeout",
		Attr:      "StopTimeout",
		Converter: converters.Duration,
		Scopable:  true,
		Text:      "Wait for <duration> before declaring the action a failure. Takes precedence over :kw:`timeout`.",
		Example:   "1m30s",
	},
	{
		Option:    "provision_timeout",
		Attr:      "ProvisionTimeout",
		Converter: converters.Duration,
		Scopable:  true,
		Text:      "Wait for <duration> before declaring the action a failure. Takes precedence over :kw:`timeout`.",
		Example:   "1m30s",
	},
	{
		Option:    "unprovision_timeout",
		Attr:      "UnprovisionTimeout",
		Converter: converters.Duration,
		Scopable:  true,
		Text:      "Wait for <duration> before declaring the action a failure. Takes precedence over :kw:`timeout`.",
		Example:   "1m30s",
	},
	{
		Option:    "run_timeout",
		Attr:      "RunTimeout",
		Converter: converters.Duration,
		Scopable:  true,
		Text:      "Wait for <duration> before declaring the action a failure. Takes precedence over :kw:`timeout`.",
		Example:   "1m30s",
	},
	{
		Option:    "sync_timeout",
		Attr:      "SyncTimeout",
		Converter: converters.Duration,
		Scopable:  true,
		Text:      "Wait for <duration> before declaring the action a failure. Takes precedence over :kw:`timeout`.",
		Example:   "1m30s",
	},
	{
		Option:     "access",
		Attr:       "Access",
		Kind:       kind.Or(kind.Vol),
		Inherit:    keywords.InheritHead,
		Default:    "rwo",
		Candidates: []string{"rwo", "roo", "rwx", "rox"},
		Scopable:   true,
		Text:       "The access mode of the volume.\n``rwo`` is Read Write Once,\n``roo`` is Read Only Once,\n``rwx`` is Read Write Many,\n``rox`` is Read Only Many.\n``rox`` and ``rwx`` modes are served by flex volume services.",
	},
	{
		Option:    "size",
		Attr:      "Size",
		Inherit:   keywords.InheritHead,
		Kind:      kind.Or(kind.Vol),
		Scopable:  true,
		Converter: converters.Size,
		Text:      "The size used by this volume in its pool.",
	},
	{
		Option:   "pool",
		Attr:     "Pool",
		Inherit:  keywords.InheritHead,
		Kind:     kind.Or(kind.Vol),
		Scopable: true,
		Text:     "The name of the pool this volume was allocated from.",
	},
	{
		Option: "type",
		Text:   "The resource driver name.",
	},
}

func driverIDFromRID(t Configurer, section string) (driver.ID, error) {
	sectionTypeKey := key.T{
		Section: section,
		Option:  "type",
	}
	sectionType := t.Config().Get(sectionTypeKey)
	rid, err := resourceid.Parse(section)
	if err != nil {
		return driver.ID{}, err
	}
	did := driver.ID{
		Group: rid.DriverGroup(),
		Name:  sectionType,
	}
	return did, nil
}

func keywordLookup(store keywords.Store, k key.T, kd kind.T, sectionType string) keywords.Keyword {
	switch k.Section {
	case "data", "env":
		return keywords.Keyword{
			Option:   "*", // trick IsZero()
			Scopable: true,
			Required: false,
		}
	}
	driverGroup := driver.GroupUnknown
	rid, err := resourceid.Parse(k.Section)
	if err == nil {
		driverGroup = rid.DriverGroup()
	}

	if kw := store.Lookup(k, kd, sectionType); !kw.IsZero() {
		// base keyword
		return kw
	}

	for _, i := range driver.ListWithGroup(driverGroup) {
		allocator, ok := i.(func() resource.Driver)
		if !ok {
			continue
		}
		kws := allocator().Manifest().Keywords
		if kws == nil {
			continue
		}
		store := keywords.Store(kws)
		if kw := store.Lookup(k, kd, sectionType); !kw.IsZero() {
			return kw
		}
	}
	return keywords.Keyword{}
}
