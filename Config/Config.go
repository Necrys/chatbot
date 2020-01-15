package config

import "../Common"
import "encoding/json"
import "io/ioutil"

type CfgProxySettings struct {
    Server   string
    User     string
    Password string
}

type CfgTelegramSettings struct {
    Token         string
    ProxySettings CfgProxySettings
}

type CfgSlackSettings struct {
    Token string
}

type CfgAdmin struct {
    UserId   string
    Password string
}

type CfgHistory struct {
    Service string
    Dir     string
}

type APISettings struct {
    Port uint16
}

type HomeMonitorSettings struct {
    HistorySize           int
    TemperatureThresholds common.Span
    HumidityThresholds    common.Span
    PressureThresholds    common.Span
    Subscribers           []string
}

type Config struct {
    Telegram       CfgTelegramSettings
    Slack          CfgSlackSettings
    Commands       []string
    CommandAliases map[string]string
    Admins         []CfgAdmin
    Debug          bool
    History        []CfgHistory
    API            APISettings
    HomeMon        HomeMonitorSettings
    CalendUrl      string
}

func Read(cfgPath string) (*Config, error) {
    file, err := ioutil.ReadFile(cfgPath)
    if err != nil {
        return nil, err
    }

    var cfg Config
    err = json.Unmarshal(file, &cfg)
    if err != nil {
        return nil, err
    }

    return &cfg, nil
}

func ( this *Config ) Write( cfgPath string ) ( error ) {
    str, err := json.MarshalIndent( this, "", "  " )
    if err != nil {
        return err
    }

    err = ioutil.WriteFile( cfgPath, str, 0644 )
    if err != nil {
        return err
    }

    return nil
}