package telegram

import "github.com/go-telegram-bot-api/telegram-bot-api"

import "log"
import "encoding/json"

func dumpUpdate( update *tgbotapi.Update ) () {
  str, _ := json.MarshalIndent( update, "", "  " )
  log.Print( string( str ) )
}
