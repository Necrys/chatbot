package bot

import (
  "time"
  "errors"
  "io"
  "io/ioutil"
  "encoding/json"
  "os"
  "log"
)

const UserLocationDBFileName = "user_locations_db.json"

type UserLocPair struct {
  Location string
  loc      *time.Location
}

type UserLocationsDB struct {
  UserLocationsMap map[ string ]UserLocPair
}

func NewUserLocationsDB() ( *UserLocationsDB ) {
  db := &UserLocationsDB{ UserLocationsMap : make( map[ string ]UserLocPair ) }
  return db
}

func ( this *UserLocationsDB ) SetUserLocation( user string, location string ) ( error ) {
  loc, err := time.LoadLocation( location )
  if err != nil {
    return err
  }

  this.UserLocationsMap[ user ] = UserLocPair{ location, loc }
  this.SaveToFile()

  return nil
}

func ( this *UserLocationsDB ) GetUserLocation( user string ) ( *time.Location, error ) {
  if loc, ok := this.UserLocationsMap[ user ]; ok {
    return loc.loc, nil
  }

  return nil, errors.New( "User not found" )
}

func ( this *UserLocationsDB ) LoadFromFile() ( error ) {
  file, err := ioutil.ReadFile( UserLocationDBFileName )
  if err != nil {
    return err
  }

  loadedDb := make( map[ string ]UserLocPair )
  err = json.Unmarshal( file, &loadedDb )
  if err != nil {
    return err
  }

  // reset existing map
  this.UserLocationsMap = make( map[ string ]UserLocPair )
  for u, l := range loadedDb {
    loc, err := time.LoadLocation( l.Location )
    if err != nil {
      log.Panic( "Unable to find location ", l.Location )
      continue
    }

    l.loc = loc
    this.UserLocationsMap[ u ] = l
  }

  return nil
}

func ( this *UserLocationsDB ) SaveToFile() ( error ) {
  // backup existing file
  if _, err := os.Stat( UserLocationDBFileName ); err == nil {
    backupFilePath := UserLocationDBFileName + ".bak"
    in, err := os.Open( UserLocationDBFileName )
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
  DBJson, _ := json.Marshal( this.UserLocationsMap )
  err := ioutil.WriteFile( UserLocationDBFileName, DBJson, 0644 )
  if err != nil {
    return err
  }

  return nil
}