package api

import (
  "../Config"
  "fmt"
  "net/http"
)

func RunAPIHandler( cfg *config.Config ) {
  http.Handle( "/bme280/", newBME280APIHandler( cfg ) )
  hostAddr := fmt.Sprintf( ":%d", cfg.API.Port )
  go http.ListenAndServe( hostAddr, nil )
}
