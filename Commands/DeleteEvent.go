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
  isadmin := this.botCtx.IsAdmin( cmd.UserId() )
  if isadmin == false {
    cmd.Reply( "You're not The Master" )
    return true
  }

  args := strings.Trim( cmd.Args(), " \n\t" )

  eid, err := strconv.ParseUint( args, 10, 64 )
  if err != nil {
    cmd.Reply( "Failed to parse event Id" )
    return true
  }
  
  this.botCtx.DeleteEvent( eid )

  return true
}
