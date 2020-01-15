package main

import (
  "./Bot"
  "./Config"
  "./CmdProcessor"
  "./Commands"
  "./Telegram"
  "./Slack"
  "./HistoryLogger"
  "./Api"
  "log"
  "time"
  "os"
  "os/signal"
  "syscall"
  "fmt"
)

var (
	Version   = "undefined"
  Commit    = "undefined"
	BuildTime = "undefined"
	GitHash   = "undefined"
)

const cfg_path = "config.json"

func main() {
    log.Print("----- Start -----")
    log.Print( fmt.Sprintf( "----- Version: %v.%v ( %v, %v ) -----", Version, Commit, GitHash, BuildTime ) )

    log.Print("----- Load config -----")
    cfg, err := config.Read( cfg_path )
    if err != nil {
        log.Print("Failed to read configuration from file")
        return
    }
    
    botCtx, err := bot.NewContext(cfg)
    if err != nil {
        log.Print("Failed to create bot context")
        return
    }

    sigs := make( chan os.Signal, 1 )
    signal.Notify( sigs, syscall.SIGINT, syscall.SIGTERM )

    go func() {
      sig := <-sigs
      log.Println()
      log.Println( sig )
      botCtx.Waiting <- false 
    } ()

    // periodically save dbs to disk
    ticker := time.NewTicker( 3 * time.Hour )
    go func() {
      for {
        <-ticker.C
        botCtx.ChatsDb.SaveToFile()
        botCtx.SaveScheduleDBToFile()
      }
    }()

    defer botCtx.ChatsDb.SaveToFile()

    // run HTTP API handler
    api.RunAPIHandler( cfg, botCtx.HomeCtrl )

    // Create and registrate commands
    cmds := map[string]cmdprocessor.CommandProcIf {
        "stop":    commands.NewCmdStop(botCtx),
        "restart": commands.NewCmdRestart(botCtx),
        "goadmin": commands.NewCmdGoAdmin(cfg, botCtx),
        "noadmin": commands.NewCmdNoAdmin(botCtx),
        "roll":    commands.NewCmdRoll(),
        "sensors": commands.NewCmdSensorsLast(),
        "sensorshistory": commands.NewCmdSensorsGraph(),
        "getchats": commands.NewGetKnownChats( botCtx ),
        "say": commands.NewSay( botCtx ),
        "schedule": commands.NewScheduleEvent( botCtx ),
        "deleteevent": commands.NewDeleteEvent( botCtx ),
        "getevents": commands.NewGetActiveEvents( botCtx ),
        "setlocation": commands.NewSetLocation( botCtx ),
        "calend": commands.NewCalend( cfg ),
        "showkeyboard": commands.NewCmdShowKeyboard(),
    }

    // Append some basic commands to the config so it'll be registered always.
    // Not a good solution but will work at this point
    cfg.Commands = append(cfg.Commands, "stop")
    cfg.Commands = append(cfg.Commands, "restart")
    cfg.Commands = append(cfg.Commands, "goadmin")
    cfg.Commands = append(cfg.Commands, "noadmin")

    botCtx.CmdProc, err = cmdprocessor.NewCmdRegistry(cfg, cmds)
    if err != nil {
        log.Print("Failed to create command registry")
        return
    }

    history, err := history.NewLogger(cfg)
    if err != nil {
        log.Print("Failed to init history logger")
    }

    var tgListener *telegram.Listener = nil
    if cfg.Telegram.Token != "" {
        tgListener, err = telegram.NewListener(cfg, botCtx, history)
        if err != nil {
            log.Print("Failed to init telegram listener: ", err)
            return
        } else {
            tgListener.Start( botCtx.CmdProc )
            bot.AddListener( "telegram", tgListener )
        }
    }

    var slackListener *slack.Listener = nil
    if cfg.Slack.Token != "" {
        slackListener, err = slack.NewListener(cfg, botCtx, history)
        if err != nil {
            log.Print("Failed to init slack listener: ", err)
            return
        } else {
            slackListener.Start( botCtx.CmdProc )
            bot.AddListener( "slack", slackListener )
        }
    }

    err = botCtx.LoadScheduleDBFromFile()
    if err != nil {
      log.Print( err )
    }

    helloUsers := []string{}
    for _, v := range cfg.Admins {
      helloUsers = append( helloUsers, v.UserId )
    }

    botCtx.SayHello( fmt.Sprintf( "Bot, version %s.%s ( %s, %s ) started", Version, Commit, GitHash, BuildTime ), helloUsers )

    <- botCtx.Waiting
    log.Print("----- Stop command received -----")

    if tgListener != nil {
        tgListener.Stop()
    }

    if slackListener != nil {
        slackListener.Stop()
    }

    ticker.Stop()
    cfg.Write( cfg_path )

    log.Print("----- Stop -----")
}
