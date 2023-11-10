package library

import (
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

type Env struct {
	viper.Viper
}
// GetStringWithDefault if key not found in env file use default value
func (e *Env) GetStringWithDefault(key, def string) string {
	if val := e.Get(key); val != nil {
		return cast.ToString(val)
	}
	return def
}


// GetIntWithDefault if key not found in env file use default value
func (e *Env) GetIntWithDefault(key string, def int) int {
	if val := e.Get(key); val != nil {
		return cast.ToInt(val)
	}
	return def
}


// GetBoolWithDefault if key not found in env file use default value
func (e *Env) GetBoolWithDefault(key string, def bool) bool {
	if val := e.Get(key); val != nil {
		return cast.ToBool(val)
	}
	return def
}


func NewEnv(fileName, fileType, filePath string) (env *Env,err error) {
	v := viper.New()
	env = &Env{
		*v,
	}
	env.SetConfigName(fileName)
	env.SetConfigType(fileType)
	env.SetConfigFile(filePath + fileName)
	err = env.ReadInConfig()
	return env, err
}

func NewEnvFromConfig(cfg *EnvConfig) (env *Env, err error) {
	v := viper.New()
	env = &Env{
		*v,
	}
	env.SetConfigName(cfg.FileName)
	env.SetConfigType(cfg.FileType)
	env.SetConfigFile(cfg.FilePath + cfg.FileName)
	err = env.ReadInConfig()
	return env, err
}
