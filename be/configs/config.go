package configs

type Config struct {
	Config struct {
		Env string `yaml:"env"`

		Server struct {
			Host    string `yaml:"host"`
			Port    string `yaml:"port"`
			Timeout struct {
				Server int `yaml:"server"`
				Read   int `yaml:"read"`
				Write  int `yaml:"write"`
				Idle   int `yaml:"idle"`
			} `yaml:"timeout"`
		} `yaml:"server"`

		Ws struct {
			Host    string `yaml:"host"`
			Port    string `yaml:"port"`
			Timeout struct {
				Server int `yaml:"server"`
				Read   int `yaml:"read"`
				Write  int `yaml:"write"`
				Idle   int `yaml:"idle"`
			} `yaml:"timeout"`
		} `yaml:"ws"`

		MySQL struct {
			Host     string `yaml:"host"`
			Port     string `yaml:"port"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
			Database string `yaml:"database"`
			ConnPool struct {
				MaxIdle   int `yaml:"maxIdle"`
				ConnLimit int `yaml:"connLimit"`
				IdleTime  int `yaml:"idleTime"`
			} `yaml:"connPool"`
		} `yaml:"mysql"`

		Redis struct {
			Host          string `yaml:"host"`
			Port          string `yaml:"port"`
			MinIddleConns int    `yaml:"minIddleCons"`
			PoolSize      int    `yaml:"poolSize"`
			PoolTimeout   int    `yaml:"poolTimeout"`
			Password      string `yaml:"password"`
			Db            int    `yaml:"db"`
		} `yaml:"redis"`

		Oauth2 struct {
			GoogleClientID     string   `yaml:"googleClientID"`
			GoogleClientSecret string   `yaml:"googleClientSecret"`
			Scopes             []string `yaml:"scopes"`
			RedirectEndpoint   string   `yaml:"redirectEndpoint"`
		} `yaml:"oauth2"`
	} `yaml:"config"`
}
