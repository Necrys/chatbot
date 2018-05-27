package slack

import "../Bot"
import "../Config"
import "../CmdProcessor"
import "../HistoryLogger"
import "github.com/nlopes/slack"
import "log"

type CmdId int

const (
    CmdStop = 0
)

type ListenerCmd struct {
    id CmdId
}

type Listener struct {
    bot     *bot.Context
    api     *slack.Client
    rtm     *slack.RTM
    control chan ListenerCmd
    logger  *history.ServiceLogger
}

func NewListener(cfg *config.Config, botCtx *bot.Context, logger *history.Logger) (*Listener, error) {
    this := &Listener { bot:     botCtx,
                        api:     slack.New(cfg.Slack.Token),
                        control: make(chan ListenerCmd) }

    this.logger, _ = logger.GetServiceLogger("Slack")
    this.rtm = this.api.NewRTM()

    go this.rtm.ManageConnection()

    return this, nil
}

func dumpEvent(event *slack.RTMEvent) {
    log.Printf("=== New RTMEvent ===")
    log.Printf("Type: %v", event.Type)
}

func (this *Listener) processEvent(event *slack.RTMEvent, cmdHandler *cmdprocessor.CmdRegistry) () {
    switch ev := event.Data.(type) {
    case *slack.MessageEvent:
        if this.bot.Debug == true {
            log.Printf("message event: \"%v\", from: \"%v\"", ev.Text, ev.User)
        }

        cmd, args := cmdprocessor.SplitCommandAndArgs(ev.Text, ev.User)
        cmdCtx := &CommandCtx { listener:  this,
                                message:   ev.Text,
                                command:   cmd,
                                args:      args,
                                userId:    ev.User,
                                userName:  ev.Username,
                                channelId: ev.Channel }

        cmdHandler.HandleCommand(cmdCtx)

    case *slack.ConnectingEvent:
        log.Printf("connecting, attempt: %v, count: %v", ev.Attempt, ev.ConnectionCount)

    case *slack.ConnectionErrorEvent:
        log.Printf("Error: %s\n", ev.Error())

    case *slack.LatencyReport:
        // do nothing

    default:
        if this.bot.Debug == true {
            dumpEvent(event)
        }
    }
}

func (this *Listener) listen(cmdHandler *cmdprocessor.CmdRegistry) () {
    log.Printf("slack.Listener: Start listener thread")

    isRunning := true

    for isRunning == true {
        select {
        case ev := <- this.rtm.IncomingEvents:
            this.processEvent(&ev, cmdHandler)

        case cmd := <- this.control:
            if cmd.id == CmdStop {
                log.Printf("slack.Listener: Stop command received")
                isRunning = false
            }
        }
    }
    
    log.Printf("slack.Listener: Exit listener thread")
}

func (this *Listener) Start(cmdHandler *cmdprocessor.CmdRegistry) () {
    go this.listen(cmdHandler)
}

func (this *Listener) Stop() () {
    this.control <- ListenerCmd { id: CmdStop }
}
