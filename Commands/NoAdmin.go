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
    tokens := PrepareMessage(cmd.Message())
    if len(tokens) == 0 {
        return false
    }
    if tokens[0] != "noadmin" {
        return false
    }
    
    this.botCtx.Admins[cmd.User()] = false
    cmd.Reply(fmt.Sprintf("You're not The Master, @%s", cmd.User()))

    return true
}
