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
			Env      string `yaml:"env"`
			Host     string `yaml:"host"`
			Port     string `yaml:"port"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
			Database string `yaml:"database"`
			Region   string `yaml:"region"`
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

		Storage struct {
			ProjectID  string `yaml:"projectId"`
			BucketName string `yaml:"bucketName"`
		} `yaml:"storage"`

		Cron struct {
			Time string `yaml:"time"`
		} `yaml:"cron"`

		Midtrans struct {
			ServerKey   string `yaml:"serverKey"`
			Environment string `yaml:"environment"`
		} `yaml:"midtrans"`

		AWS struct {
			AccessKeyID     string `yaml:"accessKeyID"`
			SecretAccessKey string `yaml:"secretAccessKey"`
			Region          string `yaml:"region"`
		} `yaml:"aws"`

		Kafka struct {
			Dial      string `yaml:"dial"`
			Host      string `yaml:"host"`
			Port      string `yaml:"port"`
			Topic     string `yaml:"topic"`
			Partition int    `yaml:"partition"`
		} `yaml:"kafka"`

		ElasticSearch struct {
			Host      string `yaml:"host"`
			Port      string `yaml:"port"`
			Transport struct {
				MaxIdleConnsPerHost int `yaml:"maxIdleConnsPerHost"`
				Timeout             int `yaml:"timeout"`
				KeepAlive           int `yaml:"keepAlive"`
			} `yaml:"transport"`
			MaxRetries int `yaml:"maxRetries"`
		} `yaml:"elasticsearch"`

		Email struct {
			Host   string `yaml:"host"`
			Port   int    `yaml:"port"`
			Sender struct {
				Name     string `yaml:"name"`
				Email    string `yaml:"email"`
				Password string `yaml:"password"`
			} `yaml:"sender"`
		} `yaml:"email"`
	} `yaml:"config"`
}
