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

func (this* CommandCtx) Reply(msg string) () {
    if this.listener.bot.Debug == true {
        log.Printf("reply to %v: %v", this.channelId, msg)
    }

    this.listener.rtm.SendMessage(this.listener.rtm.NewOutgoingMessage(msg, this.channelId))
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
