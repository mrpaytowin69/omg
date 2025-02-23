package rescontainerlxc

import (
	"opensvc.com/opensvc/core/driver"
	"opensvc.com/opensvc/core/keywords"
	"opensvc.com/opensvc/core/manifest"
	"opensvc.com/opensvc/drivers/rescontainer"
	"opensvc.com/opensvc/util/converters"
)

var (
	drvID = driver.NewID(driver.GroupContainer, "lxc")
)

func init() {
	driver.Register(drvID, New)
}

// Manifest exposes to the core the input expected by the driver.
func (t T) Manifest() *manifest.T {
	m := manifest.New(drvID, t)
	m.AddContext([]manifest.Context{
		{
			Key:  "path",
			Attr: "Path",
			Ref:  "object.path",
		},
		{
			Key:  "object_id",
			Attr: "ObjectID",
			Ref:  "object.id",
		},
		{
			Key:  "nodes",
			Attr: "Nodes",
			Ref:  "object.nodes",
		},
		{
			Key:  "dns",
			Attr: "DNS",
			Ref:  "node.dns",
		},
	}...)
	m.AddKeyword([]keywords.Keyword{
		{
			Option:   "data_dir",
			Aliases:  []string{"container_data_dir"},
			Attr:     "DataDir",
			Scopable: true,
			Text:     "If this keyword is set, the service configures a resource-private container data store. This setup is allows stateful service relocalization.",
			Example:  "/srv/svc1/data/containers",
		},
		{
			Option:       "rootfs",
			Attr:         "RootDir",
			Provisioning: true,
			Text:         "Sets the root fs directory of the container",
			Example:      "/srv/svc1/data/containers",
		},
		{
			Option:       "cf",
			Attr:         "ConfigFile",
			Provisioning: true,
			Text:         "Defines a lxc configuration file in a non-standard location.",
			Example:      "/srv/svc1/config",
		},
		{
			Option:       "template",
			Attr:         "Template",
			Provisioning: true,
			Text:         "Sets the url of the template unpacked into the container root fs or the name of the template passed to :cmd:`lxc-create`.",
			Example:      "ubuntu",
		},
		{
			Option:       "template_options",
			Attr:         "TemplateOptions",
			Provisioning: true,
			Converter:    converters.Shlex,
			Text:         "The arguments to pass through :cmd:`lxc-create` to the per-template script.",
			Example:      "--release focal",
		},
		{
			Option:       "create_secrets_environment",
			Attr:         "CreateSecretsEnvironment",
			Provisioning: true,
			Scopable:     true,
			Converter:    converters.Shlex,
			Text:         "Set variables in the :cmd:`lxc-create` execution environment. A whitespace separated list of ``<var>=<secret name>/<key path>``. A shell expression spliter is applied, so double quotes can be around ``<secret name>/<key path>`` only or whole ``<var>=<secret name>/<key path>``. Variables are uppercased.",
			Example:      "CRT=cert1/server.crt PEM=cert1/server.pem",
		},
		{
			Option:       "create_configs_environment",
			Attr:         "CreateConfigsEnvironment",
			Provisioning: true,
			Scopable:     true,
			Converter:    converters.Shlex,
			Text:         "Set variables in the :cmd:`lxc-create` execution environment. The whitespace separated list of ``<var>=<config name>/<key path>``. A shell expression spliter is applied, so double quotes can be around ``<config name>/<key path>`` only or whole ``<var>=<config name>/<key path>``. Variables are uppercased.",
			Example:      "CRT=cert1/server.crt PEM=cert1/server.pem",
		},
		{
			Option:       "create_environment",
			Attr:         "CreateEnvironment",
			Provisioning: true,
			Scopable:     true,
			Converter:    converters.Shlex,
			Text:         "Set variables in the :cmd:`lxc-create` execution environment. The whitespace separated list of ``<var>=<config name>/<key path>``. A shell expression spliter is applied, so double quotes can be around ``<config name>/<key path>`` only or whole ``<var>=<config name>/<key path>``. Variables are uppercased.",
			Example:      "FOO=bar BAR=baz",
		},
		rescontainer.KWRCmd,
		rescontainer.KWName,
		rescontainer.KWHostname,
		rescontainer.KWStartTimeout,
		rescontainer.KWStopTimeout,
		rescontainer.KWSCSIReserv,
		rescontainer.KWPromoteRW,
		rescontainer.KWNoPreemptAbort,
		rescontainer.KWOsvcRootPath,
		rescontainer.KWGuestOS,
	}...)
	return m
}
