package commands

import (
  "../Bot"
  "../Config"
  "../CmdProcessor"
  "fmt"
  "crypto/md5"
  "encoding/hex"
)

type CmdGoAdmin struct {
  botCtx     *bot.Context
  AdminsList map[ string ]string
}

func NewCmdGoAdmin( cfg *config.Config, inBotCtx *bot.Context ) ( *CmdGoAdmin ) {
  this := &CmdGoAdmin { botCtx:     inBotCtx, 
                        AdminsList: make( map[ string ]string ) }
  
  for _, v := range cfg.Admins {
    this.AdminsList[ v.UserId ] = v.Password
  }
  
  return this
}

func (this* CmdGoAdmin) HandleCommand(cmdCtx cmdprocessor.CommandCtxIf) (bool) {
  pass, ok := this.AdminsList[ cmdCtx.UserId() ]
  if ok != true {
    cmdCtx.Reply( fmt.Sprintf( "You're not The Master, @%s", cmdCtx.User() ) )
    return true
  }

  hasher := md5.New()
  hasher.Write( []byte( cmdCtx.Args() ) )
  passwordMD5 := hex.EncodeToString( hasher.Sum( nil ) )

  if pass != passwordMD5 {
    cmdCtx.Reply( fmt.Sprintf( "Invalid password, @%s", cmdCtx.User() ) )
  } else {
    this.botCtx.Admins[ cmdCtx.UserId() ] = true
    cmdCtx.Reply( fmt.Sprintf( "I serve You, @%s", cmdCtx.User() ) )
  }

  return true
}
