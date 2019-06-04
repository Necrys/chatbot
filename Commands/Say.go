package commands

import (
  "../CmdProcessor"
  "../Bot"
  "strings"
)

type Say struct {
  botCtx *bot.Context
}

func NewSay( ctx *bot.Context ) ( *Say ) {
  this := &Say { botCtx: ctx }
  return this
}

func ( this* Say ) HandleCommand( cmd cmdprocessor.CommandCtxIf ) ( bool ) {
  isadmin := this.botCtx.IsAdmin( cmd.UserId() )
  if isadmin == false {
    cmd.Reply( "You're not The Master" )
    return true
  }

  args := strings.Trim( cmd.Args(), " \n\t" )
  cmdLine := strings.SplitN( args, " ", 2 )
  if len( cmdLine ) == 0 {
    cmd.Reply( "No command line provided" )
    return true
  }

  cmd.SayToChat( cmdLine[ 1 ], cmdLine[ 0 ] )

  return true
}
