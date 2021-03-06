package api

import (
  "../Config"
  "../Bot"
  "../Common"
  "container/ring"
  "log"
  "net/http"
  "strconv"
  "time"
)

type bme280APIHandler struct {
  homeCtrl *bot.SmartHomeController
}

var SensorsHistory *ring.Ring

func newBME280APIHandler ( cfg *config.Config, hc *bot.SmartHomeController ) *bme280APIHandler {
  bot.Debug( "newBME280APIHandler %v+", hc )

  handler := &bme280APIHandler { homeCtrl : hc }
  SensorsHistory = ring.New( cfg.HomeMon.HistorySize )
  return handler
}

func ( this *bme280APIHandler ) ServeHTTP ( w http.ResponseWriter, r *http.Request ) {
  // parse arguments
  data := common.SensorData{ time.Now(), 0.0, 0.0, 0.0 }
  var err error

  times, ok := r.URL.Query()["ts"]

  if !ok || len( times[0] ) < 1 {
    log.Println( "Temperature is missing" )
  } else {
    data.Timestamp, err = time.Parse( "02012006150405MST", times[0] )
    if err != nil {
      log.Printf( "Error parsing timestamp: %v\n", err )
    }
  }
  
  temps, ok := r.URL.Query()["t"]

  if !ok || len( temps[0] ) < 1 {
    log.Println( "Temperature is missing" )
  } else {
    data.Temperature, err = strconv.ParseFloat( string( temps[0] ), 64 )
    if err != nil {
      log.Printf( "Error parsing temperature: %v\n", err )
    }
  }

  hums, ok := r.URL.Query()["h"]

  if !ok || len( hums[0] ) < 1 {
    log.Println( "Humidity is missing" )
  } else {
    data.Humidity, err = strconv.ParseFloat( string( hums[0] ), 64 )
    if err != nil {
      log.Printf( "Error parsing humidity: %v\n", err )
    }
  }

  press, ok := r.URL.Query()["p"]

  if !ok || len( press[0] ) < 1 {
    log.Println( "Pressure is missing" )
  } else {
    data.Pressure, err = strconv.ParseFloat( string( press[0] ), 64 )
    if err != nil {
      log.Printf( "Error parsing pressure: %v\n", err )
    }
  }

  SensorsHistory.Value = data
  SensorsHistory = SensorsHistory.Next()

  this.homeCtrl.UpdateAtmoData( data )
}
