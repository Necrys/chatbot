package config

import "encoding/json"
import "io/ioutil"

type CfgProxySettings struct {
    Server string
    User   string
    Pass   string
}

type CfgTelegramSettings struct {
    Token         string
    ProxySettings CfgProxySettings
}

type CfgAdmin struct {
    UserId   string
    Password string
}

type Config struct {
    Telegram CfgTelegramSettings
    Commands []string
    Admins   []CfgAdmin
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
