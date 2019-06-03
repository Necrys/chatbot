package bot

import (
  "errors"
)

type CommandListener interface {
  PushMessage( string, string )
}

var serviceListenersMap = map[ string ] CommandListener{}

func AddListener( serviceName string, listener CommandListener ) {
  serviceListenersMap[ serviceName ] = listener
}

func GetListener( serviceName string ) ( CommandListener, error ) {
  if listener, ok := serviceListenersMap[ serviceName ]; ok {
    return listener, nil
  }

  return nil, errors.New( "No listener found" )
}
