package cmdprocessor

import "../Config"
import "bytes"
import "errors"
import "strings"

type CommandCtxIf interface {
    SayToChat( string, string ) ()
    Reply(string) ()
    ReplyTo(string, string, bool) ()
    UploadPNG( *bytes.Buffer ) ()
    Message() (string)
    User() (string)
    UserId() (string)
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
func SplitCommandAndArgs(message string, botName string) (string, string) {
    tokens := strings.SplitN(strings.Trim(message, " \n\t"), " ", 2)

    if len(tokens) == 0 {
        return "", ""
    }

    tokens[0] = strings.ToLower(strings.Trim(tokens[0], "/"))

    // check if it's a direct command to bot
    // if not, command will be ignored
    if strings.Contains(tokens[0], "@") == true {
        cmd := strings.Split(tokens[0], "@")
        if cmd[1] == strings.ToLower( botName ) {
            tokens[0] = cmd[0]
        } else {
            tokens[0] = ""
        }
    }

    if len(tokens) == 1 {
        return tokens[0], ""
    }
    
    tokens[1] = strings.Trim(tokens[1], " \n\t")

    return tokens[0], tokens[1]
}

func ( this *CmdRegistry ) IsCommand( name string ) ( bool ) {
  _, ok := this.commands[ name ]
  return ok
}
