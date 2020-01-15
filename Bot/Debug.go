package bot

import (
  "log"
)

var debugOn bool

func Debug( fmt string, args ...interface{} ) {
  if debugOn {
    log.Printf( fmt, args... )
  }
}

func SetDebug( flag bool ) {
  debugOn = flag
}
