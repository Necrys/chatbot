package commands

import (
  "../Bot"
  "../CmdProcessor"
  "strings"
  "fmt"
)

type GetActiveEvents struct {
}

func NewGetActiveEvents() ( *GetActiveEvents ) {
  this := &GetActiveEvents {}
  return this
}

func ( this* GetActiveEvents ) HandleCommand( cmd cmdprocessor.CommandCtxIf ) ( bool ) {
  eventsMap := bot.GetActiveEvents()
  var str strings.Builder
  for _, e := range eventsMap {
    str.WriteString( fmt.Sprintf( "%2d: channel: %s, next time: %+v, command: %s, periodic: %v, period: %+v\n", e.Id, e.TargetChannel, e.TimePoint, e.CommandString, e.IsPeriodic, e.PeriodSeconds.Seconds() ) )
  }

  cmd.Reply( str.String() )
  return true
}
