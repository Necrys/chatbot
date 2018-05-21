package commands

import "../Bot"
import "../Config"
import "../CmdProcessor"
import "strings"
import "fmt"

type CmdGoAdmin struct {
    botCtx     *bot.Context
    AdminsList map[string]string
}

func NewCmdGoAdmin(cfg *config.Config, inBotCtx *bot.Context) (*CmdGoAdmin){
    this := &CmdGoAdmin { botCtx:     inBotCtx, 
                          AdminsList: make(map[string]string) }
    
    for _, v := range cfg.Admins {
        this.AdminsList[v.UserId] = v.Password
    }
    
    return this
}

func (this* CmdGoAdmin) HandleCommand(cmd cmdprocessor.CommandCtxIf) (bool) {
    tokens := strings.Split(cmd.Message(), " ")
    if len(tokens) == 0 {
        return false
    }
    if strings.ToLower(tokens[0]) != "goadmin" {
        return false
    }
    
    var passArg string = ""
    if len(tokens) > 1 {
        passArg = tokens[1]
    }
    
    pass, ok := this.AdminsList[cmd.User()]
    if ok != true {
        cmd.Reply(fmt.Sprintf("You're not The Master, ", cmd.User()))
        return true
    }

    if pass != passArg {
        cmd.Reply(fmt.Sprintf("Invalid password, ", cmd.User()))
    } else {
        this.botCtx.Admins[cmd.User()] = true
        cmd.Reply(fmt.Sprintf("I serve You, ", cmd.User()))
    }

    return true
}
