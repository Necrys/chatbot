package telegram

import "../Config"
import "../CmdProcessor"
import "../Bot"
import "../HistoryLogger"
import "github.com/go-telegram-bot-api/telegram-bot-api"
import "golang.org/x/net/proxy"
import "errors"
import "net/http"
import "log"
import "time"
import "fmt"
import "strings"
import "strconv"

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
    bot     *bot.Context
    logger  *history.ServiceLogger
    handler *cmdprocessor.CmdRegistry
}

func NewListener(cfg *config.Config, botCtx *bot.Context, logger *history.Logger) (*Listener, error) {
    this := &Listener { api:     nil,
                        control: make(chan ListenerCmd),
                        bot:     botCtx }

    this.logger, _ = logger.GetServiceLogger("Telegram")

    if cfg.Telegram.ProxySettings.Server != "" {
        auth := proxy.Auth { User     : cfg.Telegram.ProxySettings.User,
                             Password : cfg.Telegram.ProxySettings.Password }
        dialer, err := proxy.SOCKS5("tcp", cfg.Telegram.ProxySettings.Server, &auth, proxy.Direct)
        if err != nil {
            return nil, errors.New("Failed to init SOCKS5 proxy dialer")
        }

        httpTransport := &http.Transport { DisableKeepAlives: false }
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

func (this* Listener) logChatMessage(upd *tgbotapi.Update) () {
    if this.logger == nil {
        return
    }

    chatTitle := strings.Replace(upd.Message.Chat.Title, " ", "_", -1)
    var chatId string = ""
    if len(chatTitle) == 0 {
        chatId = fmt.Sprintf("%x", (uint64)(upd.Message.Chat.ID))
    } else {
        chatId = fmt.Sprintf("%x_%v", (uint64)(upd.Message.Chat.ID), chatTitle)
    }

    this.logger.Printf(chatId, "%v    %v (%v %v): %v",
        time.Unix((int64)(upd.Message.Date), 0), upd.Message.From.UserName,
        upd.Message.From.FirstName, upd.Message.From.LastName, upd.Message.Text)
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
            if this.bot.Debug == true {
                dumpUpdate(&update)
            }

            if update.Message == nil {
                continue
            }

            this.logChatMessage(&update)
            var chatName string
            if len( update.Message.Chat.Title ) == 0 {
              chatName = update.Message.From.UserName
            } else {
              chatName = update.Message.Chat.Title
            }
            this.bot.ChatsDb.InsertChat( "telegram", chatName, strconv.FormatInt( update.Message.Chat.ID, 16 ) )

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
    this.handler = cmdHandler
    go this.listen(cmdHandler)
}

func (this* Listener) Stop() () {
    this.control <- ListenerCmd { id: CmdStop }
}

func ( this* Listener ) PushMessage( channelId string, cmdLine string ) () {
  cmd, args := cmdprocessor.SplitCommandAndArgs( cmdLine, this.api.Self.UserName )

  iChanId, err := strconv.ParseInt( channelId, 16, 64 )
  if err != nil {
    log.Printf( "telegram.PushMessage: can't parse channel Id" )
    return
  }
  
  cmdCtx := &CommandCtx { listener: this,
                          user:     "__thisbot__",
                          msg:      cmdLine,
                          mid:      0,
                          cid:      iChanId,
                          command:  cmd,
                          args:     args }

  this.handler.HandleCommand( cmdCtx )
}
