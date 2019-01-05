package config

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
    HistorySize int
}

type Config struct {
    Telegram CfgTelegramSettings
    Slack    CfgSlackSettings
    Commands []string
    Admins   []CfgAdmin
    Debug    bool
    History  []CfgHistory
    API      APISettings
    HomeMon  HomeMonitorSettings
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
