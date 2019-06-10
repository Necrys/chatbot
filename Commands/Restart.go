package commands

import (
  "../Bot"
  "../CmdProcessor"
  "log"
  "fmt"
  "os/exec"
)

type CmdRestart struct {
  botCtx     *bot.Context
}

func NewCmdRestart( inBotCtx *bot.Context ) ( *CmdRestart ) {
  this := &CmdRestart { botCtx : inBotCtx }
  return this
}

func ( this* CmdRestart ) HandleCommand( c cmdprocessor.CommandCtxIf ) ( bool ) {
  isadmin := this.botCtx.IsAdmin( c.UserId() )
  if isadmin == false {
    c.Reply("You're not The Master")
    return true
  }

  c.Reply( "Restarting..." )

  err := exec.Command( "./restart.sh" ).Start()
  if err != nil {
    log.Panic( err )
    c.Reply( fmt.Sprintf( "%v", err ) )
    return true
  }

  return true
}
