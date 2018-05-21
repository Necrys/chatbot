package telegram

import "../Config"
import "../CmdProcessor"
import "github.com/go-telegram-bot-api/telegram-bot-api"
import "golang.org/x/net/proxy"
import "errors"
import "net/http"

type Listener struct {
    bot       *tgbotapi.BotAPI
    isRunning bool
}

type CommandCtx struct {
    listener *Listener
    user     string
    msg      string
    mid      int
    cid      int64
}

func NewListener(cfg *config.Config) (*Listener, error) {
    this := &Listener { bot: nil }

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

        bot, err := tgbotapi.NewBotAPIWithClient(cfg.Telegram.Token, httpClient)
        if err != nil {
            return nil, errors.New("Failed to init Telegram bot API with proxy")
        }
        this.bot = bot
    } else {
        bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
        if err != nil {
            return nil, errors.New("Failed to init Telegram bot API")
        }
        this.bot = bot
    }

    return this, nil
}

func (this* Listener) listen(cmdHandler *cmdprocessor.CmdRegistry) () {
    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60
    this.isRunning = true

    updates, err := this.bot.GetUpdatesChan(u)

    if err != nil {
        return
    }

    for this.isRunning {
        for update := range updates {
            if update.Message == nil {
              continue
            }

            //msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
            //msg.ReplyToMessageID = update.Message.MessageID
            //this.bot.Send(msg)
            
            cmd := &CommandCtx { listener: this,
                                 user:     update.Message.From.UserName,
                                 msg:      update.Message.Text,
                                 mid:      update.Message.MessageID,
                                 cid:      update.Message.Chat.ID }

            cmdHandler.HandleCommand(cmd)
        }
    }
}

func (this* Listener) Start(cmdHandler *cmdprocessor.CmdRegistry) () {
    go this.listen(cmdHandler)
    //this.listen(cmdHandler)
}

func (this* Listener) Stop() () {
    this.isRunning = false
}

func (this* CommandCtx) Message() (string) {
    return this.msg
}

func (this* CommandCtx) Reply(text string) () {
    msg := tgbotapi.NewMessage(this.cid, text)
    msg.ReplyToMessageID = this.mid
    this.listener.bot.Send(msg)
}

func (this* CommandCtx) User() (string) {
    return this.user
}