package cmdprocessor

import "../Config"
import "errors"
import "strings"

type CommandCtxIf interface {
    Reply(string) ()
    Message() (string)
    User() (string)
    Command() (string)
    Args() (string)
}

type CommandProcIf interface {
    HandleCommand(CommandCtxIf) (bool)
}

type CmdRegistry struct {
    commands map[string]CommandProcIf
}

func NewCmdRegistry(cfg *config.Config, commands map[string]CommandProcIf) (*CmdRegistry, error) {
    if len(commands) == 0 {
        return nil, errors.New("No command processors passed")
    }

    this := &CmdRegistry{ make(map[string]CommandProcIf) }

    for _, v := range cfg.Commands {
        cmd, ok := commands[v]
        if ok == true {
            this.commands[v] = cmd
            delete(commands, v)
        }
    }

    return this, nil
}

func (this *CmdRegistry) HandleCommand(cmd CommandCtxIf) (bool) {
    cmdProc, ok := this.commands[cmd.Command()]
    if ok != true {
        return false
    }

    go cmdProc.HandleCommand(cmd)

    return true
}

// split first token and the rest of the message
// convert first token to lower case
func SplitCommandAndArgs(message string) (string, string) {
    tokens := strings.SplitN(strings.Trim(message, " \n\t"), " ", 2)

    if len(tokens) == 0 {
        return "", ""
    }

    tokens[0] = strings.ToLower(tokens[0])

    if len(tokens) == 1 {
        return tokens[0], ""
    }
    
    tokens[1] = strings.Trim(tokens[1], " \n\t")

    return tokens[0], tokens[1]
}
