package history

import "../Config"
import "fmt"
import "os"
import "bufio"
import "errors"
import "time"
import "log"

type ChannelLog struct {
    file *bufio.Writer
    date time.Time
    fail bool // just to avoid log spam
}

type ServiceLogger struct {
    channelLogs map[string]*ChannelLog
    logDir      string
}

type Logger struct {
    serviceLogs map[string]*ServiceLogger
}


func checkCreateDir(path string) (error) {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        err = os.MkdirAll(path, 0755)
        return err
    }

    return nil
}

func NewLogger (cfg *config.Config) (*Logger, error) {
    if cfg == nil {
        return nil, errors.New("No config provided")
    }

    this := &Logger { serviceLogs: make(map[string]*ServiceLogger) }

    for _, v := range cfg.History {
        // try to create logs dir
        if err := checkCreateDir(v.Dir); err != nil {
            log.Printf("Failed to create service history log directory \"%v\", error: %v", v.Dir, err)
        }

        this.serviceLogs[v.Service] = &ServiceLogger {
            channelLogs: make(map[string]*ChannelLog),
            logDir:      v.Dir }
    }

    return this, nil
}

func (this *Logger) GetServiceLogger (service string) (*ServiceLogger, error) {
    logger, ok := this.serviceLogs[service]
    if ok != true {
        return nil, errors.New(fmt.Sprintf("No history configured for \"%v\" service", service))
    }

    return logger, nil
}

func (this* ServiceLogger) Printf (channel string, format string, args ...interface{}) () {
    writer, ok := this.channelLogs[channel]
    if ok != true {
        this.channelLogs[channel] = &ChannelLog {
            file: nil,
            date: time.Now().Local(),
            fail: false }

        writer = this.channelLogs[channel]

        // try to create channel dir
        channelDir := this.logDir + "/" + channel
        if err := checkCreateDir(channelDir); err != nil {
            log.Printf("Failed to create channel history log directory \"%v\", error: %v", channelDir, err)
            writer.fail = true
            return
        }

        fullLogPath := channelDir + "/" + fmt.Sprintf("%d-%.2d-%.2d", writer.date.Year(), writer.date.Month(), writer.date.Day()) + ".log"
        log.Printf("Create history log: %v", fullLogPath)
        f, err := os.Create(fullLogPath)
        if err != nil {
            log.Printf("Failed to create channel history log file \"%v\", error: %v", fullLogPath, err)
            writer.fail = true
            return
        }

        writer.file = bufio.NewWriter(f)
    }

    // avoid log spamming
    if writer.fail == true {
        return
    }

    date := time.Now().Local()
    if writer.date.Day() != date.Day() || writer.date.Month() != date.Month() || writer.date.Year() != date.Year() {
        channelDir := this.logDir + "/" + channel
        fullLogPath := channelDir + "/" + fmt.Sprintf("%d-%.2d-%.2d", writer.date.Year(), writer.date.Month(), writer.date.Day()) + ".log"
        f, err := os.Create(fullLogPath)
        if err != nil {
            log.Printf("Failed to create channel history log file \"%v\", error: %v", fullLogPath, err)
            writer.fail = true
            return
        }

        writer.date = date
        writer.file = bufio.NewWriter(f)
    }

    fmt.Fprintf(writer.file, format, args...)
    fmt.Fprintf(writer.file, "\n")
    writer.file.Flush()
}
