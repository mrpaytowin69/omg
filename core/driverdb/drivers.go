package driverdb

import (
	_ "opensvc.com/opensvc/drivers/arrayfreenas"
	_ "opensvc.com/opensvc/drivers/pooldirectory"
	_ "opensvc.com/opensvc/drivers/poolfreenas"
	_ "opensvc.com/opensvc/drivers/poolshm"
	_ "opensvc.com/opensvc/drivers/poolvirtual"
	_ "opensvc.com/opensvc/drivers/poolzpool"
	_ "opensvc.com/opensvc/drivers/resappforking"
	_ "opensvc.com/opensvc/drivers/resappsimple"
	_ "opensvc.com/opensvc/drivers/rescertificatetls"
	_ "opensvc.com/opensvc/drivers/resdiskdisk"
	_ "opensvc.com/opensvc/drivers/resdiskloop"
	_ "opensvc.com/opensvc/drivers/resdisklv"
	_ "opensvc.com/opensvc/drivers/resdiskmd"
	_ "opensvc.com/opensvc/drivers/resdiskraw"
	_ "opensvc.com/opensvc/drivers/resdiskvg"
	_ "opensvc.com/opensvc/drivers/resexposeenvoy"
	_ "opensvc.com/opensvc/drivers/resfsdir"
	_ "opensvc.com/opensvc/drivers/resfsflag"
	_ "opensvc.com/opensvc/drivers/resfshost"
	_ "opensvc.com/opensvc/drivers/resfszfs"
	_ "opensvc.com/opensvc/drivers/resiphost"
	_ "opensvc.com/opensvc/drivers/resiproute"
	_ "opensvc.com/opensvc/drivers/resrouteenvoy"
	_ "opensvc.com/opensvc/drivers/ressharenfs"
	_ "opensvc.com/opensvc/drivers/restaskhost"
	_ "opensvc.com/opensvc/drivers/resvhostenvoy"
	_ "opensvc.com/opensvc/drivers/resvol"
)

func Load() {
	return
}
