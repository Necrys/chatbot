package bot

import (
  "time"
  "io"
  "io/ioutil"
  "encoding/json"
  "os"
  "log"
)

type ScheduledEvent struct {
  Id            uint64
  TargetChannel string
  CommandString string
  IsPeriodic    bool
  PeriodSeconds time.Duration
  TimePoint     time.Time
  deleted       bool
  listener      CommandListener
}

var schedulerDB = map[ uint64 ]ScheduledEvent{}
var schedulerIdCounter uint64
var schedulerFreeEventIds = []uint64{}

const SchedulerDBFileName = "scheduler_db.json"

func ( this *Context ) scheduleNextTimedEvent( e ScheduledEvent ) {
  time.AfterFunc( e.PeriodSeconds, func() {
    if _, ok := schedulerDB[ e.Id ]; ok {
      if schedulerDB[ e.Id ].deleted {
        delete( schedulerDB, e.Id )
        this.SaveScheduleDBToFile()
        schedulerFreeEventIds = append( schedulerFreeEventIds, e.Id )

        return
      }

      e.listener.PushMessage( e.TargetChannel, e.CommandString )
      if e.IsPeriodic {
        eventData := schedulerDB[ e.Id ]
        eventData.TimePoint = eventData.TimePoint.Add( eventData.PeriodSeconds )
        schedulerDB[ e.Id ] = eventData
        this.SaveScheduleDBToFile()
        this.scheduleNextTimedEvent( e )
      } else {
        // seems like a dead code
        delete( schedulerDB, e.Id )
        this.SaveScheduleDBToFile()
        schedulerFreeEventIds = append( schedulerFreeEventIds, e.Id )
      }
    }
  } )
}

func ( this *Context ) scheduleEvent( e ScheduledEvent ) {
  dur := e.TimePoint.Sub( time.Now() )
  time.AfterFunc( dur, func() {
    if _, ok := schedulerDB[ e.Id ]; ok {
      if schedulerDB[ e.Id ].deleted {
        delete( schedulerDB, e.Id )
        this.SaveScheduleDBToFile()
        schedulerFreeEventIds = append( schedulerFreeEventIds, e.Id )

        return
      }

      e.listener.PushMessage( e.TargetChannel, e.CommandString )
      if e.IsPeriodic {
        eventData := schedulerDB[ e.Id ]
        eventData.TimePoint = eventData.TimePoint.Add( eventData.PeriodSeconds )
        schedulerDB[ e.Id ] = eventData
        this.SaveScheduleDBToFile()
        this.scheduleNextTimedEvent( e )
      } else {
        delete( schedulerDB, e.Id )
        this.SaveScheduleDBToFile()
        schedulerFreeEventIds = append( schedulerFreeEventIds, e.Id )
      }
    }
  } )
}

func getNewSchedulerId() ( uint64 ) {
  var newId uint64
  if len( schedulerFreeEventIds ) > 0 {
    newId = schedulerFreeEventIds[ len( schedulerFreeEventIds ) - 1 ]
    schedulerFreeEventIds = schedulerFreeEventIds[ :len( schedulerFreeEventIds ) - 1 ]
  } else {
    newId = schedulerIdCounter
    schedulerIdCounter += 1
  }
  
  return newId
}

func ( this* Context ) ScheduleEvent( cid string, cmdLine string, inTimePoint time.Time, inIsPeriodic bool, inPeriod uint64 ) ( uint64, error ) {
  if this.Debug {
    log.Printf( "[ScheduleEvent] cid: %s, cmdLine: %s, inTimePoint: %+v, inIsPeriodic: %+v, inPeriod: %+v",
      cid, cmdLine, inTimePoint, inIsPeriodic, inPeriod )
  }

  service, err := this.ChatsDb.GetChatService( cid )
  if err != nil {
    return 0, err
  }

  serviceListener, err := GetListener( service )
  if err != nil {
    return 0, err
  }

  e := ScheduledEvent { Id            : getNewSchedulerId(),
                        TargetChannel : cid,
                        CommandString : cmdLine,
                        IsPeriodic    : inIsPeriodic,
                        PeriodSeconds : time.Duration( time.Second * time.Duration( inPeriod ) ),
                        TimePoint     : inTimePoint,
                        deleted       : false,
                        listener      : serviceListener,
                      }

  schedulerDB[ e.Id ] = e
  this.SaveScheduleDBToFile()

  this.scheduleEvent( e )

  return e.Id, nil
}

func ( this* Context ) DeleteEvent( id uint64 ) {
  Debug( "[DeleteEvent( %v )]", id )

  if e, ok := schedulerDB[ id ]; ok {
    e.deleted = true
    schedulerDB[ id ] = e
    Debug( "[DeleteEvent] Event %v marked as deleted", id )
  }
}

func ( this* Context ) SaveScheduleDBToFile() ( error ) {
  Debug( "[SaveScheduleDBToFile]" )

  // backup existing file
  if _, err := os.Stat( SchedulerDBFileName ); err == nil {
    backupFilePath := SchedulerDBFileName + ".bak"
    in, err := os.Open( SchedulerDBFileName )
    if err != nil {
      Debug( "[SaveScheduleDBToFile] failed to open \"%s\" file: %+v", SchedulerDBFileName, err )
      return err
    }

    defer in.Close()

    out, err := os.Create( backupFilePath )
    if err != nil {
      Debug( "[SaveScheduleDBToFile] failed to create backup file \"%s\" file: %+v",  backupFilePath, err )
      return err
    }
    defer out.Close()

    if _, err = io.Copy( out, in ); err != nil {
      Debug( "[SaveScheduleDBToFile] failed to copy data to backup file: %+v", err )
      return err
    }

    err = out.Sync()
    if err != nil {
      Debug( "[SaveScheduleDBToFile] failed to sync data backup file: %+v", err )
      return err
    }
  }
  
  // remove deleted events
  deletedEvents := []uint64{}
  for id, e := range schedulerDB {
    if e.deleted {
      deletedEvents = append( deletedEvents, id )
    }
  }

  for _, id := range deletedEvents {
    delete( schedulerDB, id )
    schedulerFreeEventIds = append( schedulerFreeEventIds, id )
  }

  // write json
  DBJson, err := json.Marshal( schedulerDB )
  if err != nil {
    Debug( "[SaveScheduleDBToFile] failed to marshal data: %+v", err )
    return err
  }
  
  err = ioutil.WriteFile( SchedulerDBFileName, DBJson, 0644 )
  if err != nil {
    Debug( "[SaveScheduleDBToFile] failed to write \"%s\" file: %+v", SchedulerDBFileName, err )
    return err
  }

  return nil
}

func ( this *Context ) LoadScheduleDBFromFile() ( error ) {
  Debug( "[LoadScheduleDBFromFile]" )

  file, err := ioutil.ReadFile( SchedulerDBFileName )
  if err != nil {
    Debug( "[LoadScheduleDBFromFile] failed to read file \"%s\": %+v", SchedulerDBFileName, err )
    return err
  }

  // reset existing map and id's
  schedulerDB = make( map[ uint64 ]ScheduledEvent )
  schedulerFreeEventIds = []uint64{}
  schedulerIdCounter = 0

  loadedDB := make( map[ uint64 ]ScheduledEvent )
  err = json.Unmarshal( file, &loadedDB )
  if err != nil {
    Debug( "[LoadScheduleDBFromFile] failed to unmarshal json data: %+v", err )
    return err
  }

  maxId := schedulerIdCounter

  if schedulerDB != nil {
    for eid, e := range loadedDB {
      Debug( "[LoadScheduleDBFromFile] loading event, eid: %+v, event: %+v", eid, e )
      service, err := this.ChatsDb.GetChatService( e.TargetChannel )
      if err != nil {
        Debug( "[LoadScheduleDBFromFile] failed to get chat \"%s\" service: %+v", e.TargetChannel, err )
        continue
      }

      serviceListener, err := GetListener( service )
      if err != nil {
        Debug( "[LoadScheduleDBFromFile] failed to get service \"%s\" listener: %+v", service, err )
        continue
      }

      loadedEvt := ScheduledEvent { Id            : eid,
                                    TargetChannel : e.TargetChannel,
                                    CommandString : e.CommandString,
                                    IsPeriodic    : e.IsPeriodic,
                                    PeriodSeconds : e.PeriodSeconds,
                                    TimePoint     : e.TimePoint,
                                    deleted       : false,
                                    listener      : serviceListener,
                                  }

      schedulerDB[ loadedEvt.Id ] = loadedEvt

      if maxId < eid {
        maxId = eid
      }

      this.scheduleEvent( loadedEvt )
    }

    // fix free id list
    usedIds := make ( map[ uint64 ] uint64 )
    for i := 0; i <= int( maxId ); i++ {
      usedIds[ uint64( i ) ] = uint64( i )
    }

    for eid, _ := range loadedDB {
      delete( usedIds, eid )
    }

    for eid, _ := range usedIds {
      schedulerFreeEventIds = append( schedulerFreeEventIds, eid )
    }

    schedulerIdCounter = maxId + 1
  }
  
  return nil
}

func GetActiveEvents() ( map[ uint64 ]ScheduledEvent ) {
  result := map[ uint64 ]ScheduledEvent{}

  for eid, e := range schedulerDB {
    if e.deleted == false {
      result[ eid ] = e
    }
  }

  return result
}
