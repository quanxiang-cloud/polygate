package config

import (
	"crypto/tls"
	"io/ioutil"
	"time"

	"github.com/Shopify/sarama"
	"github.com/quanxiang-cloud/cabin/logger"

	"gopkg.in/yaml.v2"
)

// Conf 全局配置文件
var Conf *Config

// DefaultPath 默认配置路径
var DefaultPath = "./configs/config.yml"

// Config 配置文件
type Config struct {
	Port            string            `yaml:"port"`
	Model           string            `yaml:"model"`
	Schema          string            `yaml:"schema"`
	APIFilterConfig APIFilterConfig   `yaml:"apiFilter"`
	RedrectService  map[string]string `yaml:"redirectService"`
	Gate            Gate              `yaml:"gate"`
	Remotes         RemotesConfig     `yaml:"remotes"`
	Proxy           Proxy             `yaml:"proxy"`
	Log             *logger.Config    `yaml:"log"`

	Kafka   Kafka   `yaml:"kafka"`
	Handler Handler `yaml:"handler"`
}

// NewConfig 获取配置配置
func NewConfig(path string) (*Config, error) {
	if path == "" {
		path = DefaultPath
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(file, &Conf)
	if err != nil {
		return nil, err
	}

	return Conf, nil
}

// GateLimitRate config
type GateLimitRate struct {
	Enable        bool `yaml:"enable" validate:"required"`
	RatePerSecond int  `yaml:"ratePerSecond" validate:"required,min=1,max=100000"`
}

// // GateAPIBlock config
// type GateAPIBlock struct {
// 	Enable        bool  `yaml:"enable" validate:"required"`
// 	MaxAllowError int64 `yaml:"maxAllowError" validate:"required,min=1,max=100"`
// 	BlockSeconds  int64 `yaml:"blockSeconds" validate:"required,min=10,max=3600"`
// 	APITimeoutMS  int64 `yaml:"apiTimeoutMS" validate:"required,min=1"`
// }

// GateIPBlock config
type GateIPBlock struct {
	Enable bool     `yaml:"enable" validate:"required"`
	White  []string `yaml:"white"`
	Black  []string `yaml:"black"`
}

// Gate config
type Gate struct {
	// APIBlock  GateAPIBlock  `yaml:"apiBlock"`
	LimitRate GateLimitRate `yaml:"limitRate"`
	IPBlock   GateIPBlock   `yaml:"ipBlock"`
}

// HTTPClientConfig http client
type HTTPClientConfig struct {
	Addr         string        `yaml:"addr"`
	MaxIdleConns int           `yaml:"maxIdleConns"`
	Timeout      time.Duration `yaml:"timeout"`
}

// RemotesConfig presents URLs of auth config
type RemotesConfig struct {
	OauthToken HTTPClientConfig `yaml:"oauthToken"`
	OauthKey   HTTPClientConfig `yaml:"oauthKey"`
	Goalie     HTTPClientConfig `yaml:"goalie"`
}

// Proxy proxy
type Proxy struct {
	Timeout               time.Duration `yaml:"timeout"`
	KeepAlive             time.Duration `yaml:"keepAlive"`
	MaxIdleConns          int           `yaml:"maxIdleConns"`
	IdleConnTimeout       time.Duration `yaml:"idleConnTimeout"`
	TLSHandshakeTimeout   time.Duration `yaml:"tlsHandshakeTimeout"`
	ExpectContinueTimeout time.Duration `yaml:"expectContinueTimeout"`
}

// APIFilterConfig is the api black list
type APIFilterConfig struct {
	White []string
}

// Handler handler
type Handler struct {
	Topic          string `yaml:"topic"`
	Group          string `yaml:"group"`
	NumOfProcessor int    `yaml:"numOfProcessor"`
	Buffer         int    `yaml:"buffer"`
}

// Kafka kafka config
type Kafka struct {
	Sarama sarama.Config

	Broker []string `yaml:"broker"`
	TLS    *tls.Config
}

func pre(conf Kafka) *sarama.Config {
	config := sarama.NewConfig()

	// TLS
	config.Net.TLS.Enable = conf.Sarama.Net.TLS.Enable
	config.Net.TLS.Config = conf.Sarama.Net.TLS.Config

	return config
}
