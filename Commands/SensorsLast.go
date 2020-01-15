package commands

import "../CmdProcessor"
import "../Api"
import "../Common"
import "fmt"

type CmdSensorsLast struct {
}

func NewCmdSensorsLast() ( *CmdSensorsLast ) {
    this := &CmdSensorsLast {  }
    return this
}

func ( this* CmdSensorsLast ) HandleCommand( cmdCtx cmdprocessor.CommandCtxIf ) ( bool ) {
  if api.SensorsHistory.Len() == 0 {
    cmdCtx.Reply( "No data" )
    return true
  }

  // Prev() because we need to compensate Next() call in data appending procedure
  val := api.SensorsHistory.Prev().Value
  if val == nil {
    cmdCtx.Reply( "No data" )
    return true
  }
  
  data := val.(common.SensorData)
  cmdCtx.Reply( fmt.Sprintf(
    "%s\n    🌡 Температура: %.2f°C\n    💧 Относ.влажность: %.2f%%\n    ⏱ Атм.давление: %.2f мм рт.ст.\n",
    data.Timestamp.String(), data.Temperature, data.Humidity, data.Pressure ) )

  return true
}
