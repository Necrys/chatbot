package commands

import (
  "../Bot"
  "../CmdProcessor"
  "strings"
  "fmt"
)

type GetActiveEvents struct {
  botCtx *bot.Context
}

func NewGetActiveEvents( ctx *bot.Context ) ( *GetActiveEvents ) {
  this := &GetActiveEvents { botCtx : ctx }
  return this
}

func ( this* GetActiveEvents ) HandleCommand( cmd cmdprocessor.CommandCtxIf ) ( bool ) {
  isadmin := this.botCtx.IsAdmin( cmd.UserId() )
  if isadmin == false {
    cmd.Reply( "You're not The Master" )
    return true
  }

  eventsMap := bot.GetActiveEvents()
  var str strings.Builder
  for _, e := range eventsMap {
    str.WriteString( fmt.Sprintf( "%2d: channel: %s, next time: %+v, command: %s, periodic: %v, period: %+v\n", e.Id, e.TargetChannel, e.TimePoint, e.CommandString, e.IsPeriodic, e.PeriodSeconds.Seconds() ) )
  }

  cmd.Reply( str.String() )
  return true
}
