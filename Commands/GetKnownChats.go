package commands

import (
  "../Bot"
  "../CmdProcessor"
)

type GetKnownChats struct {
  botCtx *bot.Context
}

func NewGetKnownChats( inBotCtx *bot.Context ) ( *GetKnownChats ) {
  this := &GetKnownChats { botCtx: inBotCtx }
  return this
}

func ( this* GetKnownChats ) HandleCommand( cmd cmdprocessor.CommandCtxIf ) ( bool ) {
  isadmin := this.botCtx.IsAdmin(cmd.UserId())
  if isadmin == false {
    cmd.Reply("You're not The Master")
  } else {
    var replyText string
    for service, chats := range this.botCtx.ChatsDb.ServiceToChatsListMap {
      replyText = replyText + service + "\n"
      for name, id := range chats {
        replyText = replyText + "  \"" + name + "\": \"" + id + "\"\n"
      }
    }

    cmd.Reply( replyText )
  }

  return true
}
