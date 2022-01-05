package global

import (
	"crypto/md5"
	"os"
	"path/filepath"
	"time"

	"github.com/jkstack/natpass/code/utils"
	"github.com/lwch/runtime"
	"github.com/lwch/yaml"
)

// Configure server configure
type Configure struct {
	Listen       uint16
	Enc          [md5.Size]byte
	TLSKey       string
	TLSCrt       string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	LogDir       string
	LogSize      utils.Bytes
	LogRotate    int
}

// LoadConf load configure file
func LoadConf(dir string) *Configure {
	var cfg struct {
		Listen uint16 `yaml:"listen"`
		Secret string `yaml:"secret"`
		Link   struct {
			ReadTimeout  time.Duration `yaml:"read_timeout"`
			WriteTimeout time.Duration `yaml:"write_timeout"`
		} `yaml:"link"`
		Log struct {
			Dir    string      `yaml:"dir"`
			Size   utils.Bytes `yaml:"size"`
			Rotate int         `yaml:"rotate"`
		} `yaml:"log"`
		TLS struct {
			Key string `yaml:"key"`
			Crt string `yaml:"crt"`
		} `yaml:"tls"`
	}
	runtime.Assert(yaml.Decode(dir, &cfg))
	if !filepath.IsAbs(cfg.Log.Dir) {
		dir, err := os.Executable()
		runtime.Assert(err)
		cfg.Log.Dir = filepath.Join(filepath.Dir(dir), cfg.Log.Dir)
	}
	return &Configure{
		Listen:       cfg.Listen,
		Enc:          md5.Sum([]byte(cfg.Secret)),
		TLSKey:       cfg.TLS.Key,
		TLSCrt:       cfg.TLS.Crt,
		ReadTimeout:  cfg.Link.ReadTimeout,
		WriteTimeout: cfg.Link.WriteTimeout,
		LogDir:       cfg.Log.Dir,
		LogSize:      cfg.Log.Size,
		LogRotate:    cfg.Log.Rotate,
	}
}
