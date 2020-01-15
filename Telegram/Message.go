package telegram

import (
  "github.com/go-telegram-bot-api/telegram-bot-api"
  "bytes"
  "strconv"
  "fmt"
)

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

func (this* CommandCtx) SayToChat( text string, cid string ) () {
  this.ReplyTo( text, cid, false )
}

func (this* CommandCtx) Reply( text string ) () {
  this.ReplyTo( text, strconv.FormatInt( this.cid, 16 ), true )
}

func (this* CommandCtx) ReplyNoCitation( text string ) () {
  this.ReplyTo( text, strconv.FormatInt( this.cid, 16 ), false )
}

func (this* CommandCtx) ReplyTo( text string, cid string, useCitation bool ) () {
  channelId, err := strconv.ParseInt( cid, 16, 64 )
  if err != nil {
    return
  }

  msg := tgbotapi.NewMessage( channelId, text )
  if this.mid != 0 && useCitation == true {
    msg.ReplyToMessageID = this.mid
  }
  msg.ParseMode = tgbotapi.ModeMarkdown
  this.listener.api.Send(msg)
}

func (this* CommandCtx) UploadPNG( buffer *bytes.Buffer, useCitation bool ) () {
  b := tgbotapi.FileBytes{ Name: "image.jpg", Bytes: buffer.Bytes() }
  msg := tgbotapi.NewPhotoUpload(this.cid, b)
  if this.mid != 0 && useCitation == true {
    msg.ReplyToMessageID = this.mid
  }
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

func ( this* CommandCtx ) ShowKeyboard( keyboard [][]string ) () {
  msg := tgbotapi.NewMessage( this.cid, this.msg )

  kb := tgbotapi.NewReplyKeyboard()
  for _, row := range keyboard {
    kbRow := tgbotapi.NewKeyboardButtonRow()
    for _, key := range row {
      kbRow = append( kbRow, tgbotapi.NewKeyboardButton( key ) )
    }
    kb.Keyboard = append( kb.Keyboard, kbRow )
  }

  msg.ReplyMarkup = kb

  if _, err := this.listener.api.Send( msg ); err != nil {
    this.Reply( fmt.Sprintf( "Ошибка: %v", err ) )
  }
}

func ( this* CommandCtx ) HideKeyboard() () {
  msg := tgbotapi.NewMessage( this.cid, this.msg )
  msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard( true )

  if _, err := this.listener.api.Send( msg ); err != nil {
    this.Reply( fmt.Sprintf( "Ошибка: %v", err ) )
  }
}

func ( this* CommandCtx ) HideUserCommand() {
	deleteMessageConfig := tgbotapi.DeleteMessageConfig {
		ChatID:    this.cid,
		MessageID: this.mid,
	}

  this.listener.api.DeleteMessage( deleteMessageConfig )
}
