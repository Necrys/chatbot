package bot

import (
  "../Config"
  "../CmdProcessor"
  "log"
  "fmt"
)

type Context struct {
    Admins    map[string]bool
    Waiting   chan bool
    Debug     bool
    ChatsDb   *KnownChatsDB
    UserLocDb *UserLocationsDB
    CmdProc   *cmdprocessor.CmdRegistry
    HomeCtrl  *SmartHomeController
}

func NewContext(cfg *config.Config) (*Context, error) {
  log.Print( "NewContext" )
  SetDebug( cfg.Debug )
  ctx := &Context { Admins    : make(map[string]bool),
                    Waiting   : make(chan bool),
                    Debug     : cfg.Debug,
                    ChatsDb   : NewKnownChatsDB(),
                    UserLocDb : NewUserLocationsDB() }
  
  ctx.HomeCtrl = NewSmartHomeController( cfg, ctx )

  ctx.ChatsDb.LoadFromFile()
  ctx.UserLocDb.LoadFromFile()

  ctx.Admins[ "__thisbot__" ] = true;

  return ctx, nil
}

func (this* Context) IsAdmin(UserId string) (bool) {
    flag, ok := this.Admins[UserId]
    if ok == false {
        return false
    }

    return flag
}

func ( this* Context ) SayHello( s string, users []string ) () {
  for _, v := range users {
    ch, svc, err := this.ChatsDb.GetChatAndServiceByName( v )
    if err == nil {
      l, err := GetListener( svc )
      if err != nil {
        log.Panic( fmt.Sprintf( "%v", err ) )
        continue
      }

      l.PushMessage( ch, fmt.Sprintf( "say %s %s", ch, s ) )
    }
  }
}
