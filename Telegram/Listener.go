package telegram

import "../Config"
import "../CmdProcessor"
import "github.com/go-telegram-bot-api/telegram-bot-api"
import "golang.org/x/net/proxy"
import "errors"
import "net/http"
import "log"

type CmdId int

const (
    CmdStop = 0
)

type ListenerCmd struct {
    id CmdId
}

type Listener struct {
    api     *tgbotapi.BotAPI
    control chan ListenerCmd
}

func NewListener(cfg *config.Config) (*Listener, error) {
    this := &Listener { api:     nil,
                        control: make(chan ListenerCmd) }

    if cfg.Telegram.ProxySettings.Server != "" {
        auth := proxy.Auth { User     : cfg.Telegram.ProxySettings.User,
                             Password : cfg.Telegram.ProxySettings.Password }
        dialer, err := proxy.SOCKS5("tcp", cfg.Telegram.ProxySettings.Server, &auth, proxy.Direct)
        if err != nil {
            return nil, errors.New("Failed to init SOCKS5 proxy dialer")
        }

        httpTransport := &http.Transport {}
        httpTransport.Dial = dialer.Dial
        httpClient := &http.Client { Transport: httpTransport }

        api, err := tgbotapi.NewBotAPIWithClient(cfg.Telegram.Token, httpClient)
        if err != nil {
            return nil, errors.New("Failed to init Telegram api API with proxy")
        }
        this.api = api
    } else {
        api, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
        if err != nil {
            return nil, errors.New("Failed to init Telegram api API")
        }
        this.api = api
    }

    return this, nil
}

func (this* Listener) listen(cmdHandler *cmdprocessor.CmdRegistry) () {
    log.Printf("telegram.Listener: Start listener thread")
    isRunning := true

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates, err := this.api.GetUpdatesChan(u)

    if err != nil {
        return
    }

    for isRunning == true {
        select {
        case update := <- updates:
            if update.Message == nil {
                continue
            }

            log.Printf("update has come (%s)", update.Message.Text)

            cmd, args := cmdprocessor.SplitCommandAndArgs(update.Message.Text, this.api.Self.UserName)

            cmdCtx := &CommandCtx { listener: this,
                                    user:     update.Message.From.UserName,
                                    msg:      update.Message.Text,
                                    mid:      update.Message.MessageID,
                                    cid:      update.Message.Chat.ID,
                                    command:  cmd,
                                    args:     args }

            cmdHandler.HandleCommand(cmdCtx)

        case cmd := <- this.control:
            if cmd.id == CmdStop {
                log.Printf("telegram.Listener: Stop command received")
                isRunning = false
            }
        }
    }

    log.Printf("telegram.Listener: Exit listener thread")
}

func (this* Listener) Start(cmdHandler *cmdprocessor.CmdRegistry) () {
    go this.listen(cmdHandler)
}

func (this* Listener) Stop() () {
    this.control <- ListenerCmd { id: CmdStop }
}
