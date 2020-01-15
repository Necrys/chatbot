package bot

import(
  "../Config"
  "../Common"
  "fmt"
)

type SmartHomeController struct {
  TemperatureSpan common.Span
  HumiditySpan    common.Span
  PressureSpan    common.Span
  Subscribers     []string
  bot             *Context
}

func NewSmartHomeController( cfg *config.Config, b *Context ) ( *SmartHomeController ) {
  ctrl := &SmartHomeController { TemperatureSpan : cfg.HomeMon.TemperatureThresholds,
                                 HumiditySpan    : cfg.HomeMon.HumidityThresholds,
                                 PressureSpan    : cfg.HomeMon.PressureThresholds,
                                 Subscribers     : cfg.HomeMon.Subscribers,
                                 bot             : b }

  return ctrl
}

func ( this *SmartHomeController ) UpdateAtmoData( data common.SensorData ) {
  var notifications []string

  if data.Temperature < this.TemperatureSpan.Min {
    notifications = append( notifications, fmt.Sprintf( "Температура упала ниже %.2f°C", this.TemperatureSpan.Min ) )
  } else if data.Temperature > this.TemperatureSpan.Max {
    notifications = append( notifications, fmt.Sprintf( "Температура поднялась выше %.2f°C", this.TemperatureSpan.Max ) )
  }

  if data.Humidity < this.HumiditySpan.Min {
    notifications = append( notifications, fmt.Sprintf( "Влажность упала ниже %.2f%%", this.HumiditySpan.Min ) )
  } else if data.Humidity > this.HumiditySpan.Max {
    notifications = append( notifications, fmt.Sprintf( "Влажность поднялась выше %.2f%%", this.HumiditySpan.Max ) )
  }

  if len( notifications ) > 0 {
    for _, cid := range this.Subscribers {
      service, err := this.bot.ChatsDb.GetChatService( cid )
      if err != nil {
        return
      }

      serviceListener, err := GetListener( service )
      if err != nil {
        return
      }

      for _, notif := range notifications {
        serviceListener.PushMessage( cid, fmt.Sprintf( "Say %s %s", cid, notif ) )
      }
    }
  }
}
