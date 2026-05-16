package config

type Provider interface {
	SetSetting(key string, value string)
	GetSetting(key string) (value string, ok bool)
}

