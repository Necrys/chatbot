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

  layouts := []string{ "02.01.2006T15:04:05",
                       "15:04:05" }

  timeParsed := false
  for i, layout := range layouts {
    var err error
    timePoint, err = time.Parse( layout, tokens[ 1 ] )
    if err == nil {
      timeParsed = true
      
      // second layout without date was used, set local date
      if i == 1 {
        t := time.Now()

        loc, err := time.LoadLocation( "UTC" )
        if err != nil {
          log.Panic( err )
          cmd.Reply( fmt.Sprintf( "LoadLocation UTC failed: %+v", err ) )
          return true
        }

        timePoint = time.Date( t.Year(), t.Month(), t.Day(), timePoint.Hour(), timePoint.Minute(), timePoint.Second(), 0, loc )
      }
      
      break
    }
  }

  if timeParsed != true {
    cmd.Reply( "Bad time format" )
    return true
  }

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

  _, tzOffset := time.Now().Zone()
  timePoint = timePoint.Add( time.Second * time.Duration( -tzOffset ) )

  id, err := this.botCtx.ScheduleEvent( channel, cmdLine, timePoint, isPeriodic, period )
  if err != nil {
    cmd.Reply( fmt.Sprintf( "%+v", err ) )
    return true
  }

  cmd.Reply( fmt.Sprintf( "Command scheduled, id: %v", id ) )
  return true
}
