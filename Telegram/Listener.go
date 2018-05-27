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
import "bytes"
import "time"
import "fmt"
import "strings"

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

func makeIndentString(level int) (string) {
    var indent bytes.Buffer
    for i := 0; i < level; i++ {
        indent.WriteString("  ")
    }

    return indent.String()
}

func dumpUser(user *tgbotapi.User, indentLevel int, name string) () {
    if user == nil {
        return
    }

    indent := makeIndentString(indentLevel)

    log.Printf("%v%v {", indent, name)
    log.Printf("%v  ID: %v",           indent, user.ID)
    log.Printf("%v  FirstName: %v",    indent, user.FirstName)
    log.Printf("%v  LastName: %v",     indent, user.LastName)
    log.Printf("%v  UserName: %v",     indent, user.UserName)
    log.Printf("%v  LanguageCode: %v", indent, user.LanguageCode)
    log.Printf("%v  IsBot: %v",        indent, user.IsBot)
    log.Printf("%v}", indent)
}

func dumpChat(chat *tgbotapi.Chat, indentLevel int, name string) () {
    if chat == nil {
        return
    }

    indent := makeIndentString(indentLevel)

    log.Printf("%v%v {", indent, name)
    log.Printf("%v  ID: %v",                  indent, chat.ID)
    log.Printf("%v  Type: %v",                indent, chat.Type)
    log.Printf("%v  Title: %v",               indent, chat.Title)
    log.Printf("%v  UserName: %v",            indent, chat.UserName)
    log.Printf("%v  FirstName: %v",           indent, chat.FirstName)
    log.Printf("%v  LastName: %v",            indent, chat.LastName)
    log.Printf("%v  AllMembersAreAdmins: %v", indent, chat.AllMembersAreAdmins)
    log.Printf("%v  Photo: %v",               indent, chat.Photo)
    log.Printf("%v  Description: %v",         indent, chat.Description)
    log.Printf("%v  InviteLink: %v",          indent, chat.InviteLink)
    log.Printf("%v}", indent)
}

func dumpUpdate(upd *tgbotapi.Update) () {
    log.Printf("========== New update, id: %v ==========", upd.UpdateID)
    if upd.Message != nil {
        log.Printf("Message: {")
        log.Printf("  MessageID: %v", upd.Message.MessageID)
        dumpUser(upd.Message.From, 1, "From")
        log.Printf("  Date: %v", upd.Message.Date)
        dumpChat(upd.Message.Chat, 1, "Chat")
        dumpUser(upd.Message.ForwardFrom, 1, "ForwardFrom")
        dumpChat(upd.Message.ForwardFromChat, 1, "ForwardFromChat")
        log.Printf("  ForwardFromMessageID", upd.Message.ForwardFromMessageID)
        log.Printf("  ForwardDate", upd.Message.ForwardDate)
        log.Printf("  ReplyToMessage")
        log.Printf("}")
    }
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
