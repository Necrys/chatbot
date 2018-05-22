package commands

import "../Bot"
import "../CmdProcessor"
import "fmt"

type CmdNoAdmin struct {
    botCtx *bot.Context
}

func NewCmdNoAdmin(inBotCtx *bot.Context) (*CmdNoAdmin) {
    this := &CmdNoAdmin { botCtx: inBotCtx }
    return this
}

func (this* CmdNoAdmin) HandleCommand(cmd cmdprocessor.CommandCtxIf) (bool) {
    this.botCtx.Admins[cmd.UserId()] = false
    cmd.Reply(fmt.Sprintf("You're no more The Master, @%s", cmd.User()))

    return true
}
