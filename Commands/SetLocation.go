package commands

import (
  "../Bot"
  "../CmdProcessor"
  "fmt"
)

type SetLocation struct {
  botCtx *bot.Context
}

func NewSetLocation( inBotCtx *bot.Context ) ( *SetLocation ) {
  this := &SetLocation { botCtx: inBotCtx }
  return this
}

func ( this* SetLocation ) HandleCommand( cmd cmdprocessor.CommandCtxIf ) ( bool ) {
  err := this.botCtx.UserLocDb.SetUserLocation( cmd.User(), cmd.Args() )
  if err != nil {
    cmd.Reply( fmt.Sprintf( "Failed to set location: %+v", err ) )
    return true
  }

  cmd.Reply( fmt.Sprintf( "Location \"%s\" is set for user \"%s\"", cmd.Args(), cmd.User() ) )

  return true
}
