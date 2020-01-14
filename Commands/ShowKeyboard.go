package commands

import (
  "../CmdProcessor"
  "strings"
  "strconv"
)

type CmdShowKeyboard struct {
}

func NewCmdShowKeyboard() ( *CmdShowKeyboard ) {
  this := &CmdShowKeyboard {}
  return this
}

func ( this* CmdShowKeyboard ) HandleCommand( cmdCtx cmdprocessor.CommandCtxIf ) ( bool ) {
  showKeyboard := true

  args := strings.Trim( cmdCtx.Args(), " \n\t" )
  tokens := strings.Split( args, " " )
  if len( tokens ) > 1 {
    cmdCtx.Reply( "Ожидается только один параметр (возможные значения: 0, 1)" )
    return true
  }

  if len( tokens ) == 1 {
    value, err := strconv.ParseUint( tokens[ 0 ], 10, 64 )
    if err != nil {
      cmdCtx.Reply( "Плохое значение параметра (возможные значения: 0, 1)" )
      return true
    }

    if value == 0 {
      showKeyboard = false
    } else {
      showKeyboard = true
    }
  }

  if showKeyboard == true {
    var homeKb [][]string
    var row = []string{ "Сейчас", "Температура", "Влажность", "Давление", }
    homeKb = append( homeKb, row )

    cmdCtx.ShowKeyboard( homeKb )
  } else {
    cmdCtx.HideKeyboard()  
  }

  return true
}
