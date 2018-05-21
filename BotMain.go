package main

import "./Bot"
import "./Config"
import "./CmdProcessor"
import "./Commands"
import "./Telegram"
import "log"

type AppContext struct {
    waiting chan bool
    admins  map[string]bool
}

type CmdStop struct {
    app *AppContext
}

type CmdGoAdmin struct {
}

func main() {
    botCtx := bot.Context { make(map[string]bool) }
    appCtx := AppContext { waiting: make(chan bool) }

    log.Print("----- Start -----")

    log.Print("----- Load config -----")
    cfg, err := config.Read("config.json")
    if err != nil {
        log.Print("Failed to read configuration from file")
    }
    
    cmdStop := &CmdStop { app : &appCtx }
    cmdGoAdmin := commands.NewCmdGoAdmin(cfg, &botCtx)
    cmds := map[string]cmdprocessor.CommandProcIf {
        "stop":    cmdStop,
        "goadmin": cmdGoAdmin,
    }
    
    // Append some basic commands to the config so it'll be registered always.
    // Not a good solution but will work at this point
    cfg.Commands = append(cfg.Commands, "stop")
    cfg.Commands = append(cfg.Commands, "goadmin")
    
    cmdHandler, err := cmdprocessor.NewCmdRegistry(cfg, cmds)
    if err != nil {
        log.Print("Failed to create command registry")
        return
    }
    
    var tgListener *telegram.Listener = nil
    if cfg.Telegram.Token != "" {
        tgListener, err = telegram.NewListener(cfg)
        if err != nil {
            log.Print("Failed to read configuration from file: ", err)
        } else {
            tgListener.Start(cmdHandler)
        }
    }

    _ = <- appCtx.waiting
    log.Print("----- Stop command received -----")

    if tgListener != nil {
        tgListener.Stop()
    }

    log.Print("----- Stop -----")
}
/*
func (this *CmdStop) Name() (string) {
    return "stop"
}
*/
func (this* CmdStop) HandleCommand(cmd cmdprocessor.CommandCtxIf) (bool) {
    if cmd.Message() != "stop" {
        return false
    }

    cmd.Reply("Ok, I'll stop")

    this.app.waiting <- false
    return true
}
