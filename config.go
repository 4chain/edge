package echogy

import (
	"context"
	"encoding/json"
	"github.com/echogy-io/echogy/pkg/auth"
)

type Config struct {
	LogLevel    string            `json:"logLevel"`
	LogFile     string            `json:"logFile"` // Path to log file
	EnablePProf bool              `json:"pprof"`
	HttpAddr    string            `json:"httpAddr"`
	SSHAddr     string            `json:"SSHAddr"`
	Domain      string            `json:"domain"`
	PrivateKey  string            `json:"privateKey"`
	Auth        *auth.DefaultAuth `json:"auth"`
	//Supabase    SupabaseConfig    `json:"supabase"`
}

type ContextKey struct {
	name string
}

var (
	ConfigKey = &ContextKey{name: "configKey"}
)

func WithConfig(ctx context.Context, conf []byte) context.Context {
	var c Config
	err := json.Unmarshal(conf, &c)
	if nil != err {
		panic(err)
	}
	return context.WithValue(ctx, ConfigKey, &c)
}

func ContextConfig(ctx context.Context) *Config {
	return ctx.Value(ConfigKey).(*Config)
}

func ContextAuth(ctx context.Context) auth.Auth {
	a := ContextConfig(ctx).Auth
	return auth.Update(a)
}
