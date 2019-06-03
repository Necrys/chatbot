package bot

import (
  "../Config"
  "../CmdProcessor"
)

type Context struct {
    Admins  map[string]bool
    Waiting chan bool
    Debug   bool
    ChatsDb *KnownChatsDB
    CmdProc *cmdprocessor.CmdRegistry
}

func NewContext(cfg *config.Config) (*Context, error) {
    ctx := &Context { Admins  : make(map[string]bool),
                      Waiting : make(chan bool),
                      Debug   : cfg.Debug,
                      ChatsDb : NewKnownChatsDB() }

    ctx.ChatsDb.LoadFromFile()
    SetDebug( ctx.Debug )
    return ctx, nil
}

func (this* Context) IsAdmin(UserId string) (bool) {
    flag, ok := this.Admins[UserId]
    if ok == false {
        return false
    }

    return flag
}