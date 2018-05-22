package commands

import "../Bot"
import "../CmdProcessor"

type CmdStop struct {
    botCtx *bot.Context
}

func NewCmdStop(inBotCtx *bot.Context) (*CmdStop) {
    this := &CmdStop { botCtx: inBotCtx }
    return this
}

func (this* CmdStop) HandleCommand(cmd cmdprocessor.CommandCtxIf) (bool) {
    isadmin := this.botCtx.IsAdmin(cmd.UserId())
    if isadmin == false {
        cmd.Reply("You're not The Master")
    } else {
        cmd.Reply("Ok, I'll stop")
        this.botCtx.Waiting <- false
    }

    return true
}
