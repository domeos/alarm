package g

import (
	"log"
	"runtime"
)

const (
	FALCON_ALARM_VERSION = "2.0.2"
	DOMEOS_VERSION = "0.2"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
