package config

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"time"
)

type (
	Opts struct {
		// logger
		Logger struct {
			Debug   bool `           long:"debug"        env:"DEBUG"    description:"debug mode"`
			Verbose bool `short:"v"  long:"verbose"      env:"VERBOSE"  description:"verbose mode"`
			LogJson bool `           long:"log.json"     env:"LOG_JSON" description:"Switch log output to json format"`
		}

		// DeadMansSnitch settings
		DeadMansSnitch struct {
			Token string `long:"deadmanssnitch.token"  env:"DEADMANSSNITCH_TOKEN"   description:"DeadMansSnitch access token" required:"true" json:"-"`
		}

		// general options
		ServerBind string        `long:"bind"          env:"SERVER_BIND"   description:"Server address"                default:":8080"`
		ScrapeTime time.Duration `long:"scrape.time"   env:"SCRAPE_TIME"   description:"Scrape time (time.duration)"   default:"5m"`
	}
)

func (o *Opts) GetJson() []byte {
	jsonBytes, err := json.Marshal(o)
	if err != nil {
		log.Panic(err)
	}
	return jsonBytes
}
