package commands

// https://www.calend.ru/img/export/informer_names.png

import (
  "../CmdProcessor"
  "../Config"
  "net/http"
  "fmt"
  "log"
  "io/ioutil"
  "bytes"
)

type Calend struct {
  url string
}

func NewCalend( cfg *config.Config ) ( *Calend ) {
  this := &Calend { url: cfg.CalendUrl }
  return this
}

func ( this* Calend ) HandleCommand( c cmdprocessor.CommandCtxIf ) ( bool ) {
  response, err := http.Get( this.url )
  if err != nil {
    log.Panic( err )
    c.Reply( fmt.Sprintf( "%v", err ) )
    return true
  }
  defer response.Body.Close()

  body, err := ioutil.ReadAll( response.Body )
  if err != nil {
    log.Panic( err )
    c.Reply( fmt.Sprintf( "%v", err ) )
    return true
  }

  buff := bytes.NewBuffer( body )
  
  c.UploadPNG( buff )

  return true
}