package bot

import(
  "../Config"
  "../Common"
  "fmt"
  "log"
)

type SmartHomeController struct {
  TemperatureSpan  common.Span
  TemperatureAlarm bool
  HumiditySpan     common.Span
  HumidityAlarm    bool
  PressureSpan     common.Span
  PressureAlarm    bool
  Subscribers      []string
  bot              *Context
}

func NewSmartHomeController( cfg *config.Config, b *Context ) ( *SmartHomeController ) {
  Debug( "NewSmartHomeController" )
  ctrl := &SmartHomeController { TemperatureSpan  : cfg.HomeMon.TemperatureThresholds,
                                 TemperatureAlarm : false,
                                 HumiditySpan     : cfg.HomeMon.HumidityThresholds,
                                 HumidityAlarm    : false,
                                 PressureSpan     : cfg.HomeMon.PressureThresholds,
                                 PressureAlarm    : false,
                                 Subscribers      : cfg.HomeMon.Subscribers,
                                 bot              : b }

  return ctrl
}

func ( this *SmartHomeController ) UpdateAtmoData( data common.SensorData ) {
  Debug( "UpdateAtmoData ( t: %.2f, h: %.2f )", data.Temperature, data.Humidity )
  var notifications []string

  if data.Temperature < this.TemperatureSpan.Min && this.TemperatureAlarm == false {
    notifications = append( notifications, fmt.Sprintf( "Температура упала ниже %.2f°C", this.TemperatureSpan.Min ) )
    this.TemperatureAlarm = true
    Debug( "t < tMin" )
  } else if data.Temperature > this.TemperatureSpan.Max && this.TemperatureAlarm == false {
    notifications = append( notifications, fmt.Sprintf( "Температура поднялась выше %.2f°C", this.TemperatureSpan.Max ) )
    this.TemperatureAlarm = true
    Debug( "t > tMax" )
  } else if this.TemperatureAlarm == true && data.Temperature <= this.TemperatureSpan.Max && data.Temperature >= this.TemperatureSpan.Min {
    this.TemperatureAlarm = false
    notifications = append( notifications, fmt.Sprintf( "Температура нормализовалась" ) )
    Debug( "t in range" )
  }

  if data.Humidity < this.HumiditySpan.Min && this.HumidityAlarm == false {
    notifications = append( notifications, fmt.Sprintf( "Влажность упала ниже %.2f%%", this.HumiditySpan.Min ) )
    this.HumidityAlarm = true
    Debug( "h < hMin" )
  } else if data.Humidity > this.HumiditySpan.Max && this.HumidityAlarm == false {
    notifications = append( notifications, fmt.Sprintf( "Влажность поднялась выше %.2f%%", this.HumiditySpan.Max ) )
    this.HumidityAlarm = true
    Debug( "h > hMax" )
  } else if this.HumidityAlarm == true && data.Humidity <= this.HumiditySpan.Max && data.Humidity >= this.HumiditySpan.Min {
    this.HumidityAlarm = false
    notifications = append( notifications, fmt.Sprintf( "Влажность нормализовалась" ) )
    Debug( "h in range" )
  }

  if len( notifications ) > 0 {
    for _, cid := range this.Subscribers {
      service, err := this.bot.ChatsDb.GetChatService( cid )
      if err != nil {
        log.Printf( "%v+", err )
        return
      }

      serviceListener, err := GetListener( service )
      if err != nil {
        log.Printf( "%v+", err )
        return
      }

      for _, notif := range notifications {
        serviceListener.PushMessage( cid, fmt.Sprintf( "Say %s %s", cid, notif ) )
      }
    }
  }
}
