package commands

import (
  "../Bot"
  "../CmdProcessor"
  "strings"
  "strconv"
)

type DeleteEvent struct {
  botCtx *bot.Context
}

func NewDeleteEvent( ctx *bot.Context ) ( *DeleteEvent ) {
  this := &DeleteEvent { botCtx: ctx }
  return this
}

func ( this* DeleteEvent ) HandleCommand( cmd cmdprocessor.CommandCtxIf ) ( bool ) {
  args := strings.Trim( cmd.Args(), " \n\t" )

  eid, err := strconv.ParseUint( args, 10, 64 )
  if err != nil {
    cmd.Reply( "Failed to parse event Id" )
    return true
  }
  
  this.botCtx.DeleteEvent( eid )

  return true
}
