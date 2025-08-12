package config

type ParamsProvider interface {
	GetConsoleKeyShort() string
	GetConsoleKeyFull() string
	GetConfigType() string
}

type Params struct {
	consoleKeyShort string
	consoleKeyFull  string
	configType      string
}

func (p *Params) GetConsoleKeyShort() string {
	return p.consoleKeyShort
}

func (p *Params) GetConsoleKeyFull() string {
	return p.consoleKeyFull
}

func (p *Params) GetConfigType() string {
	return p.configType
}

// NewConfigParams возвращает параметры для конфига
func NewConfigParams(consoleKeyShort, consoleKeyFull, configType string) ParamsProvider {
	return &Params{
		consoleKeyShort: consoleKeyShort,
		consoleKeyFull:  consoleKeyFull,
		configType:      configType,
	}
}
