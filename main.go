package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"github.com/webdevops/deadmanssnitch-exporter/config"
	"github.com/webdevops/deadmanssnitch-exporter/dmsclient"
)

const (
	author = "webdevops.io"

	CollectorErrorThreshold = 10
)

var (
	argparser *flags.Parser
	opts      config.Opts

	DmsClient            *dmsclient.Client
	collectorGeneralList map[string]*CollectorGeneral

	// Git version information
	gitCommit = "<unknown>"
	gitTag    = "<unknown>"
)

func main() {
	initArgparser()

	log.Infof("starting deadmanssnitch-exporter v%s (%s; %s; by %v)", gitTag, gitCommit, runtime.Version(), author)
	log.Info(string(opts.GetJson()))

	log.Infof("init DeadMansSnitch client")
	initDeadMansSnitch()

	log.Infof("starting metrics collection")
	initMetricCollector()

	log.Infof("starting http server on %s", opts.Server.Bind)
	startHTTPServer()
}

// init argparser and parse/validate arguments
func initArgparser() {
	argparser = flags.NewParser(&opts, flags.Default)
	_, err := argparser.Parse()

	// check if there is an parse error
	if err != nil {
		var flagsErr *flags.Error
		if ok := errors.As(err, &flagsErr); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Println()
			argparser.WriteHelp(os.Stdout)
			os.Exit(1)
		}
	}

	// verbose level
	if opts.Logger.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	// debug level
	if opts.Logger.Debug {
		log.SetReportCaller(true)
		log.SetLevel(log.TraceLevel)
		log.SetFormatter(&log.TextFormatter{
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				s := strings.Split(f.Function, ".")
				funcName := s[len(s)-1]
				return funcName, fmt.Sprintf("%s:%d", path.Base(f.File), f.Line)
			},
		})
	}

	// json log format
	if opts.Logger.LogJson {
		log.SetReportCaller(true)
		log.SetFormatter(&log.JSONFormatter{
			DisableTimestamp: true,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				s := strings.Split(f.Function, ".")
				funcName := s[len(s)-1]
				return funcName, fmt.Sprintf("%s:%d", path.Base(f.File), f.Line)
			},
		})
	}
}

// Init and build DeadMansSnitch client
func initDeadMansSnitch() {
	DmsClient = dmsclient.NewClient(opts.DeadMansSnitch.Token)
	DmsClient.SetUserAgent(fmt.Sprintf("deadmanssnitch-exporter/%s", gitTag))
}

func initMetricCollector() {
	var collectorName string
	collectorGeneralList = map[string]*CollectorGeneral{}

	collectorName = "Snitch"
	if opts.Scrape.Time.Seconds() > 0 {
		collectorGeneralList[collectorName] = NewCollectorGeneral(collectorName, &MetricsCollectorSnitch{})
		collectorGeneralList[collectorName].Run(opts.Scrape.Time)
	} else {
		log.WithField("collector", collectorName).Infof("collector disabled")
	}

	collectorName = "Collector"
	collectorGeneralList[collectorName] = NewCollectorGeneral(collectorName, &MetricsCollectorCollector{})
	collectorGeneralList[collectorName].Run(time.Duration(10 * time.Second))
	collectorGeneralList[collectorName].SetIsHidden(true)
}

// start and handle prometheus handler
func startHTTPServer() {
	mux := http.NewServeMux()

	// healthz
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, "Ok"); err != nil {
			log.Error(err)
		}
	})

	// readyz
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, "Ok"); err != nil {
			log.Error(err)
		}
	})

	mux.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:         opts.Server.Bind,
		Handler:      mux,
		ReadTimeout:  opts.Server.ReadTimeout,
		WriteTimeout: opts.Server.WriteTimeout,
	}
	log.Fatal(srv.ListenAndServe())
}
