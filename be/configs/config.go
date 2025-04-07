package configs

import "time"

type Config struct {
	Config struct {
		Env string `yaml:"env"`

		Server struct {
			Host    string `yaml:"host"`
			Port    string `yaml:"port"`
			Timeout struct {
				Server time.Duration `yaml:"server"`
				Read   time.Duration `yaml:"read"`
				Write  time.Duration `yaml:"write"`
				Idle   time.Duration `yaml:"idle"`
			} `yaml:"timeout"`
		} `yaml:"server"`

		Ws struct {
			Host    string `yaml:"host"`
			Port    string `yaml:"port"`
			Timeout struct {
				Server time.Duration `yaml:"server"`
				Read   time.Duration `yaml:"read"`
				Write  time.Duration `yaml:"write"`
				Idle   time.Duration `yaml:"idle"`
			} `yaml:"timeout"`
		} `yaml:"ws"`

		MySQL struct {
			Host     string `yaml:"host"`
			Port     string `yaml:"port"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
			Database string `yaml:"database"`
			ConnPool struct {
				MaxIdle   int           `yaml:"maxIdle"`
				ConnLimit int           `yaml:"connLimit"`
				IdleTime  time.Duration `yaml:"idleTime"`
			} `yaml:"connPool"`
		} `yaml:"mysql"`

		Redis struct {
			Host          string        `yaml:"host"`
			Port          string        `yaml:"port"`
			MinIddleConns int           `yaml:"minIddleCons"`
			PoolSize      int           `yaml:"poolSize"`
			PoolTimeout   time.Duration `yaml:"poolTimeout"`
			Password      string        `yaml:"password"`
			Db            int           `yaml:"db"`
		} `yaml:"redis"`
	} `yaml:"config"`
}
