package bot

import (
  "io"
  "io/ioutil"
  "encoding/json"
  "os"
  "errors"
)

type ChatRecord struct {
  Name string
  Id   string
}

const DBFileName = "known_chats_db.json"

type KnownChatsDB struct {
  ServiceToChatsListMap map[string]map[string]string
  ChatToServiceMap      map[ string ] string
}

func NewKnownChatsDB() ( *KnownChatsDB ) {
  this := &KnownChatsDB{ ServiceToChatsListMap: map[string]map[string]string{},
                         ChatToServiceMap: map[ string ] string{} }
  return this
}

func ( this *KnownChatsDB ) LoadFromFile() ( error ) {
  file, err := ioutil.ReadFile( DBFileName )
  if err != nil {
    return err
  }

  // reset existing map
  this.ServiceToChatsListMap = make( map[ string ]map[ string ]string )
  err = json.Unmarshal( file, &this.ServiceToChatsListMap )
  if err != nil {
    return err
  }

  if this.ServiceToChatsListMap == nil {
    this.ServiceToChatsListMap = make( map[ string ]map[ string ]string )
  }
  
  for service, chatList := range this.ServiceToChatsListMap {
    for _, id := range chatList {
      this.ChatToServiceMap[ id ] = service
    }
  }

  return nil
}

func ( this *KnownChatsDB ) SaveToFile() ( error ) {
  // backup existing file
  if _, err := os.Stat( DBFileName ); err == nil {
    backupFilePath := DBFileName + ".bak"
    in, err := os.Open( DBFileName )
    if err != nil {
      return err
    }

    defer in.Close()

    out, err := os.Create( backupFilePath )
    if err != nil {
      return err
    }
    defer out.Close()

    if _, err = io.Copy( out, in ); err != nil {
      return err
    }

    err = out.Sync()
    if err != nil {
      return err
    }
  }

  // write json
  DBJson, _ := json.Marshal( this.ServiceToChatsListMap )
  err := ioutil.WriteFile( DBFileName, DBJson, 0644 )
  if err != nil {
    return err
  }

  return nil
}

func ( this *KnownChatsDB ) InsertChat( service string, name string, id string ) {
  if this.ServiceToChatsListMap[ service ] == nil {
    this.ServiceToChatsListMap[ service ] = map[ string ]string{}
  }

  this.ServiceToChatsListMap[ service ][ name ] = id
  this.ChatToServiceMap[ id ] = service
}

func ( this *KnownChatsDB ) GetChatService( id string ) ( string, error ) {
  if service, ok := this.ChatToServiceMap[ id ]; ok {
    return service, nil
  }

  return "", errors.New( "No chat found" )
}

func ( this *KnownChatsDB ) GetChatAndServiceByName( name string ) ( string, string, error ) {
  for svc, chats := range this.ServiceToChatsListMap {
    if chat, ok := chats[ name ]; ok {
      return chat, svc, nil
    }
  }

  return "", "", errors.New( "No chat found" )
}
