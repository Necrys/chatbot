package cmdprocessor

import "../Config"
import "errors"

type CommandCtxIf interface {
    Reply(string) ()
    Message() (string)
    User() (string)
}

type CommandProcIf interface {
//    Name() (string)
    HandleCommand(CommandCtxIf) (bool)
}

type CmdRegistry struct {
    commands []CommandProcIf
}

func NewCmdRegistry(cfg *config.Config, commands map[string]CommandProcIf) (*CmdRegistry, error) {
    if len(commands) == 0 {
        return nil, errors.New("No command processors passed")
    }

    this := &CmdRegistry{}

    for _, v := range cfg.Commands {
        cmd, ok := commands[v]
        if ok == true {
            this.commands = append(this.commands, cmd)
            delete(commands, v)
        }
    }

    return this, nil
}

func (this *CmdRegistry) HandleCommand(cmd CommandCtxIf) (bool) {
    for _, cmdProc := range this.commands {
        result := cmdProc.HandleCommand(cmd)
        if result == true {
            return true
        }
    }

    return false
}