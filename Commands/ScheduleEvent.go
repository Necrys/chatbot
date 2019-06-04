package commands

import (
  "../Bot"
  "../CmdProcessor"
  "strings"
  "strconv"
  "fmt"
  "time"
  "log"
)

type ScheduleEvent struct {
  botCtx *bot.Context
}

func NewScheduleEvent( ctx *bot.Context ) ( *ScheduleEvent ) {
  this := &ScheduleEvent { botCtx: ctx }
  return this
}

func ( this* ScheduleEvent ) HandleCommand( cmd cmdprocessor.CommandCtxIf ) ( bool ) {
  isadmin := this.botCtx.IsAdmin( cmd.UserId() )
  if isadmin == false {
    cmd.Reply( "You're not The Master" )
    return true
  }

  args := strings.Trim( cmd.Args(), " \n\t" )

  var channel string
  isPeriodic := false
  period := uint64( 0 )
  var timePoint time.Time
  var cmdLine string

  tokens := strings.SplitN( args, " ", 3 )
  if len( tokens ) < 3 {
    cmd.Reply( "Bad syntax" )
    return true
  }

  channel = tokens[ 0 ]

  layout := "02.01.2006T15:04:05"

  // absolute value is intended as set by user in his local timezone
  // but it's parsed as UTC so we'll get the difference between UTC and user local TZ to fix the value
  userLocation, err := this.botCtx.UserLocDb.GetUserLocation( cmd.User() )
  if err != nil {
    // assume user is in bot's location
    userLocation, err = time.LoadLocation( "Local" )
    if err != nil {
      log.Panic( err )
      cmd.Reply( fmt.Sprintf( "Failed: %+v", err ) )
      return true
    }
  }

  localCurrentTime := time.Now()
  userCurrentTime := localCurrentTime.In( userLocation )

  timePoint, err = time.Parse( layout, tokens[ 1 ] )

  if err != nil {
    userEnteredTimeAffix := userCurrentTime.Format( "02.01.2006T" )
    userEnteredTime := userEnteredTimeAffix + tokens[ 1 ]

    timePoint, err = time.Parse( layout, userEnteredTime )
    if err != nil {
      fmt.Println( "Failed to parse userTime: ", err )
      cmd.Reply( "Bad time format" )
      return true
    }
  }

  _, offset := userCurrentTime.Zone()
  timePoint = timePoint.Add( time.Duration( -offset ) * time.Second )

  // check for period
  restTokens := strings.SplitN( tokens[ 2 ], " ", 3 )
  if len( restTokens ) == 3 && restTokens[ 0 ] == "-p" {
    var err error
    period, err = strconv.ParseUint( restTokens[ 1 ], 10, 64 )
    if err != nil {
      cmd.Reply( "Failed to parse period" )
      return true
    }

    isPeriodic = true
    cmdLine = restTokens[ 2 ]
  } else if len( restTokens ) == 0 {
    cmd.Reply( "No command line provided" )
    return true
  } else {
    cmdLine = tokens[ 2 ]
  }

  command := strings.SplitN( cmdLine, " ", 2 )[ 0 ]
  
  if this.botCtx.CmdProc.IsCommand( command ) == false {
    cmd.Reply( fmt.Sprintf( "Command \"%s\" is not recognized", command ) )
    return true
  }

  if this.botCtx.Debug {
    log.Printf( "ScheduleEvent: { channel: \"%+v\", cmdLine: \"%+v\", timePoint: %+v, isPeriodic: %+v, period: %+v }",
      channel, cmdLine, timePoint, isPeriodic, period )
  }

  id, err := this.botCtx.ScheduleEvent( channel, cmdLine, timePoint, isPeriodic, period )
  if err != nil {
    cmd.Reply( fmt.Sprintf( "%+v", err ) )
    return true
  }

  cmd.Reply( fmt.Sprintf( "Command scheduled, id: %v", id ) )
  return true
}
