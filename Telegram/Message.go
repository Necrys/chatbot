package telegram

import "github.com/go-telegram-bot-api/telegram-bot-api"
import "bytes"

type CommandCtx struct {
    listener *Listener
    user     string
    msg      string
    mid      int
    cid      int64
    command  string
    args     string
}

func (this* CommandCtx) Message() (string) {
    return this.msg
}

func (this* CommandCtx) Reply(text string) () {
    msg := tgbotapi.NewMessage(this.cid, text)
    msg.ReplyToMessageID = this.mid
    this.listener.api.Send(msg)
}

func (this* CommandCtx) UploadPNG( buffer *bytes.Buffer ) () {
  b := tgbotapi.FileBytes{ Name: "image.jpg", Bytes: buffer.Bytes() }
  msg := tgbotapi.NewPhotoUpload(this.cid, b)
  msg.ReplyToMessageID = this.mid
  this.listener.api.Send(msg)
}

func (this* CommandCtx) User() (string) {
    return this.user
}

func (this* CommandCtx) UserId() (string) {
    return this.user
}

func (this* CommandCtx) Command() (string) {
    return this.command
}

func (this* CommandCtx) Args() (string) {
    return this.args
}
