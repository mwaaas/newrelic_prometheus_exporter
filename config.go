package main

type Configuration struct {
	Port       string
	Debug      bool
	AppName    string
	AppNameEnv string
}

var Config = Configuration{}

func (c *Configuration) getGlobalLabel() map[string]string {
	return map[string]string{
		"AppName":    c.AppName,
		"AppNameEnv": c.AppNameEnv,
	}
}
