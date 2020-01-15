package api

import (
  "../Bot"
  "../Config"
  "fmt"
  "net/http"
)

func RunAPIHandler( cfg *config.Config, ctrl *bot.SmartHomeController ) {
  http.Handle( "/bme280/", newBME280APIHandler( cfg, ctrl ) )
  hostAddr := fmt.Sprintf( ":%d", cfg.API.Port )
  go http.ListenAndServe( hostAddr, nil )
}
