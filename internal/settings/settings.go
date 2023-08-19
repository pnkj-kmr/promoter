package settings

import (
	"fmt"
	"time"

	jsondb "github.com/pnkj-kmr/simple-json-db"
	zrl "github.com/pnkj-kmr/zap-rotate-logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	LeaderRole Role   = "leader"
	BrokerRole Role   = "broker"
	promoter   string = "PROMOTER"
	version    string = "v1.0.0"
)

var (
	C    Conf
	APPS map[string]App
	DB   jsondb.DB
	L    *zap.Logger
	AL   ActiveLeader
)

type (
	Role string

	// Conf represent self configuration
	Conf struct {
		Bind        string        `mapstructure:"bind"`
		MyId        int32         `mapstructure:"my_id"`
		ClusterId   string        `mapstructure:"cluster_id"`
		Priority    int32         `mapstructure:"priority"`
		Role        []Role        `mapstructure:"role"`
		LeaderNodes []string      `mapstructure:"leader_nodes"`
		RefreshRate time.Duration `mapstructure:"refresh_rate"`
		Timeout     time.Duration `mapstructure:"timeout"`
	}

	// App - application detail from conf file
	App struct {
		AppId       int32  `mapstructure:"app_id"`
		Priority    int32  `mapstructure:"priority"`
		Persist     int32  `mapstructure:"persist"`
		Start       string `mapstructure:"start"`
		Stop        string `mapstructure:"stop"`
		Status      string `mapstructure:"status"`
		StatusMatch string `mapstructure:"status_match"`
	}
)

func Init(f string, debug bool) {
	// starting active leader flag
	AL = &_actve_leader{}
	// loading the application configuration
	viper.SetConfigName(f)
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		L.Fatal("unable to load application configuration", zap.Error(err))
	}
	var _c struct {
		Conf Conf           `mapstructure:"config"`
		Apps map[string]App `mapstructure:"applications"`
	}
	if err := viper.Unmarshal(&_c); err != nil {
		L.Fatal("unable to read application configuration", zap.Error(err))
	}
	// assigning the values
	C = _c.Conf
	APPS = _c.Apps
	// initialing the logger along with database instance
	L = zrl.New(zrl.WithFileName(fmt.Sprintf("app-%d", C.MyId)), zrl.WithDebug(debug))
	// branding here
	logo_print()
	_db, err := jsondb.New(fmt.Sprintf("db-%d", C.MyId), &jsondb.Options{Logger: L})
	if err != nil {
		L.Fatal("unable to create database", zap.Error(err))
	}
	DB = _db
}
