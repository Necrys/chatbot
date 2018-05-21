package main

import "./Bot"
import "./Config"
import "./CmdProcessor"
import "./Commands"
import "./Telegram"
import "log"

func main() {
    botCtx := bot.Context { Admins  : make(map[string]bool),
                            Waiting : make(chan bool) }

    log.Print("----- Start -----")

    log.Print("----- Load config -----")
    cfg, err := config.Read("config.json")
    if err != nil {
        log.Print("Failed to read configuration from file")
        return
    }
    
    // Create and registrate commands
    cmdStop := commands.NewCmdStop(&botCtx)
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
            log.Print("Failed to init telegram listener: ", err)
            return
        } else {
            tgListener.Start(cmdHandler)
        }
    }

    _ = <- botCtx.Waiting
    log.Print("----- Stop command received -----")

    if tgListener != nil {
        tgListener.Stop()
    }

    log.Print("----- Stop -----")
}
