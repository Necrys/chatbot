package main

import "./Bot"
import "./Config"
import "./CmdProcessor"
import "./Commands"
import "./Telegram"
import "./Slack"
import "log"

func main() {
    log.Print("----- Start -----")

    log.Print("----- Load config -----")
    cfg, err := config.Read("config.json")
    if err != nil {
        log.Print("Failed to read configuration from file")
        return
    }
    
    botCtx, err := bot.NewContext(cfg)
    if err != nil {
        log.Print("Failed to create bot context")
        return
    }

    // Create and registrate commands
    cmds := map[string]cmdprocessor.CommandProcIf {
        "stop":    commands.NewCmdStop(botCtx),
        "goadmin": commands.NewCmdGoAdmin(cfg, botCtx),
        "noadmin": commands.NewCmdNoAdmin(botCtx),
        "roll":    commands.NewCmdRoll(),
    }
    
    // Append some basic commands to the config so it'll be registered always.
    // Not a good solution but will work at this point
    cfg.Commands = append(cfg.Commands, "stop")
    cfg.Commands = append(cfg.Commands, "goadmin")
    cfg.Commands = append(cfg.Commands, "noadmin")
    
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

    var slackListener *slack.Listener = nil
    if cfg.Slack.Token != "" {
        slackListener, err = slack.NewListener(cfg, botCtx)
        if err != nil {
            log.Print("Failed to init slack listener: ", err)
            return
        } else {
            slackListener.Start(cmdHandler)
        }
    }

    _ = <- botCtx.Waiting
    log.Print("----- Stop command received -----")

    if tgListener != nil {
        tgListener.Stop()
    }

    if slackListener != nil {
        slackListener.Stop()
    }

    log.Print("----- Stop -----")
}
