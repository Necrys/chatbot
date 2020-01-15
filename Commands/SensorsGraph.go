package commands

import (
  "../CmdProcessor"
  "../Api"
  "../Common"
	"github.com/wcharczuk/go-chart"
  "bytes"
  "fmt"
  "time"
  "github.com/wcharczuk/go-chart/drawing"
  "strings"
)

type CmdSensorsGraph struct {
}

func NewCmdSensorsGraph() ( *CmdSensorsGraph ) {
  this := &CmdSensorsGraph {}
  return this
}

type DataSource int

const (
  None        DataSource = 0
  Temperature DataSource = 1
  Humidity    DataSource = 2
  Pressure    DataSource = 3
)

func ( this* CmdSensorsGraph ) HandleCommand( cmdCtx cmdprocessor.CommandCtxIf ) ( bool ) {
  args := strings.Trim( cmdCtx.Args(), " \n\t" )

  tokens := strings.Split( args, " " )
  if len( tokens ) > 2 {
    cmdCtx.Reply( "There can be only 1 or 2 values" )
    return true
  }

  dataSrc := [2]DataSource{ Temperature, None }

  if len( args ) > 0 {
    for idx, arg := range tokens {
      lower := strings.ToLower(arg)
      if lower == "t" {
        dataSrc[idx] = Temperature
      } else if lower == "h" {
        dataSrc[idx] = Humidity
      } else if lower == "p" {
        dataSrc[idx] = Pressure
      } else {
        dataSrc[idx] = None
      }
    }
  }

  if dataSrc[0] == None && dataSrc[1] == None {
    cmdCtx.Reply( "Bad arguments (must be 't' - temperature, 'h' - humidity, 'p' - pressure)" )
    return true
  } else if dataSrc[0] == None {
    dataSrc[0] = dataSrc[1]
    dataSrc[1] = None
  }

  if api.SensorsHistory.Len() == 0 {
    cmdCtx.Reply( "No data" )
    return true
  }

  // grab data
  var timestampData []time.Time
  var dataAxis0 []float64
  var dataAxis0Units string
  var dataAxis0Name string
  var dataAxis1 []float64
  var dataAxis1Units string
  var dataAxis1Name string
  
  api.SensorsHistory.Do( func( p interface{} ) {
    if p != nil {
      data := p.( common.SensorData )

      timestampData = append( timestampData, data.Timestamp )

      if dataSrc[0] == Temperature {
        dataAxis0 = append( dataAxis0, data.Temperature )
      } else if dataSrc[0] == Humidity {
        dataAxis0 = append( dataAxis0, data.Humidity )
      } else if dataSrc[0] == Pressure {
        dataAxis0 = append( dataAxis0, data.Pressure )
      }

      if dataSrc[1] == Temperature {
        dataAxis1 = append( dataAxis1, data.Temperature )
      } else if dataSrc[1] == Humidity {
        dataAxis1 = append( dataAxis1, data.Humidity )
      } else if dataSrc[1] == Pressure {
        dataAxis1 = append( dataAxis1, data.Pressure )
      }      
    }
  })
  
  if dataSrc[0] == Temperature {
    dataAxis0Units = "°C"
    dataAxis0Name = "Temperature"
  } else if dataSrc[0] == Humidity {
    dataAxis0Units = "% RH"
    dataAxis0Name = "Humidity"
  } else if dataSrc[0] == Pressure {
    dataAxis0Units = "mmHg"
    dataAxis0Name = "Pressure"
  }

  if dataSrc[1] == Temperature {
    dataAxis1Units = "°C"
    dataAxis1Name = "Temperature"
  } else if dataSrc[1] == Humidity {
    dataAxis1Units = "% RH"
    dataAxis1Name = "Humidity"
  } else if dataSrc[1] == Pressure {
    dataAxis1Units = "mmHg"
    dataAxis1Name = "Pressure"
  }

  var graph chart.Chart
  
  if dataSrc[1] == None {
    graph = drawSingleGraph( timestampData, dataAxis0, dataAxis0Units, dataAxis0Name )
  } else {
    graph = drawDoubleGraph( timestampData, dataAxis0, dataAxis0Units, dataAxis0Name,
                             dataAxis1, dataAxis1Units, dataAxis1Name)
  }

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

  buffer := bytes.NewBuffer( []byte{} )
  err := graph.Render( chart.PNG, buffer )
  if err != nil {
    cmdCtx.Reply( fmt.Sprintf( "Graph rendering error: %v", err ) )
    return true
  }

  cmdCtx.UploadPNG( buffer )

  return true
}

func drawSingleGraph ( timeAxis []time.Time, dataAxis []float64, dataUnits string, dataName string ) chart.Chart {
	graph := chart.Chart {
    Background: chart.Style{
      Padding: chart.Box{
        Top:  20,
        Left: 20,
      },
			FillColor: drawing.ColorFromHex("01121B"),
		},
		Canvas: chart.Style{
			FillColor: drawing.ColorFromHex("01121B"),
		},
		XAxis: chart.XAxis {
			Name:         "Time",
			NameStyle:    chart.Shown(),
			Style:        chart.Shown(),
			ValueFormatter: func( v interface{} ) string {
				typed := v.( float64 )
        return time.Unix(0, int64(typed)).Format("01-02 15:04")
			},
		},
		YAxis: chart.YAxis{
			Name:      dataUnits,
			NameStyle: chart.Shown(),
			Style:     chart.Shown(), //enables / displays the y-axis
		},
		Series: []chart.Series{
			chart.TimeSeries {
        Name:    dataName,
				Style:   chart.Style{
					Hidden:      false,                             //note; if we set ANY other properties, we must set this to true.
					StrokeColor: drawing.ColorFromHex("BFFFA4"),               // will supercede defaults
					FillColor:   drawing.ColorFromHex("2B9100").WithAlpha(128), // will supercede defaults
				},
				XValues: timeAxis,
				YValues: dataAxis,
			},
		},
	}

  return graph
}

func drawDoubleGraph ( timeAxis []time.Time, dataAxis []float64, dataUnits string, dataName string,
                       secondaryDataAxis []float64, secondaryDataUnits string, secondaryDataName string, ) chart.Chart {
	graph := chart.Chart {
    Background: chart.Style{
      Padding: chart.Box{
        Top:  20,
        Left: 20,
      },
			FillColor: drawing.ColorFromHex("01121B"),
		},
		Canvas: chart.Style{
			FillColor: drawing.ColorFromHex("01121B"),
		},
		XAxis: chart.XAxis {
			Name:         "Time",
			NameStyle:    chart.Shown(),
			Style:        chart.Shown(),
			ValueFormatter: func( v interface{} ) string {
				typed := v.( float64 )
        return time.Unix(0, int64(typed)).Format("01-02 15:04")
			},
		},
		YAxis: chart.YAxis{
			Name:      dataUnits,
			NameStyle: chart.Shown(),
			Style:     chart.Shown(), //enables / displays the y-axis
		},
		YAxisSecondary: chart.YAxis{
			Name:      secondaryDataUnits,
			NameStyle: chart.Shown(),
			Style:     chart.Shown(), //enables / displays the secondary y-axis
		},
		Series: []chart.Series{
			chart.TimeSeries {
        Name:    dataName,
				Style:   chart.Style {
					Hidden:      false,                             //note; if we set ANY other properties, we must set this to true.
					StrokeColor: drawing.ColorFromHex("BFFFA4"),               // will supercede defaults
					FillColor:   drawing.ColorFromHex("2B9100").WithAlpha(128), // will supercede defaults
				},
				XValues: timeAxis,
				YValues: dataAxis,
			},
			chart.TimeSeries {
        Name:    secondaryDataName,
				Style:   chart.Style{
					Hidden:      false,                            //note; if we set ANY other properties, we must set this to true.
					StrokeColor: drawing.ColorFromHex("A8E3FF"),               // will supercede defaults
					FillColor:   drawing.ColorFromHex("04476B").WithAlpha(128), // will supercede defaults
				},
				YAxis:   chart.YAxisSecondary,
				XValues: timeAxis,
				YValues: secondaryDataAxis,
			},
		},
	}

  return graph
}