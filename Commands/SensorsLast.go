package commands

import "../CmdProcessor"
import "../Api"
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

  val := api.SensorsHistory.Value
  if val == nil {
    cmdCtx.Reply( "No data" )
    return true
  }
  
  data := val.(api.SensorData)
  cmdCtx.Reply( fmt.Sprintf(
    "%s\n    temperature: %.2f°C\n    humidity: %.2f%%\n    pressure: %.2f mmHg\n",
    data.Timestamp.String(), data.Temperature, data.Humidity, data.Pressure ) )

  return true
}
