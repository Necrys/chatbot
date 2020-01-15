package slack

import "log"
import "bytes"

type CommandCtx struct {
    listener  *Listener
    message   string
    command   string
    args      string
    userId    string
    userName  string
    channelId string
}

func ( this* CommandCtx ) SayToChat( text string, cid string ) () {
  this.ReplyTo( text, cid, false )
}

func ( this* CommandCtx ) Reply( msg string ) () {
  this.ReplyTo( msg, this.channelId, true )
}

func ( this* CommandCtx ) ReplyNoCitation( msg string ) () {
  this.ReplyTo( msg, this.channelId, false )
}

func ( this* CommandCtx ) ReplyTo( text string, cid string, useCitation bool ) () {
  if this.listener.bot.Debug == true {
    log.Printf( "reply to %v: %v", cid, text )
  }

  this.listener.rtm.SendMessage( this.listener.rtm.NewOutgoingMessage( text, cid ) )
}

func (this* CommandCtx) UploadPNG( buffer *bytes.Buffer ) () {
  // not implemented
}

func (this* CommandCtx) Message() (string) {
    return this.message
}

func (this* CommandCtx) UserId() (string) {
    return this.userId
}

func (this* CommandCtx) User() (string) {
    return this.userName
}

func (this* CommandCtx) Command() (string) {
    return this.command
}

func (this* CommandCtx) Args() (string) {
    return this.args
}

func (this* CommandCtx) ShowKeyboard( [][]string ) () {
}

func (this* CommandCtx) HideKeyboard() () {
}

func ( this* CommandCtx ) HideUserCommand() {
}
