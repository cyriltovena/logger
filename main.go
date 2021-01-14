package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/grafana/loki-client-go/loki"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/model"

	gofakeit "github.com/brianvoe/gofakeit/v6"
)

const (
	RFC3164    = "Jan 02 15:04:05"
	RFC3164Log = "<%d>%s %s %s[%d]: %s"
)

var (
	apiURL      = flag.String("url", "http://localhost:3100/loki/api/v1/push", "send log via loki api using the provided url (e.g http://localhost:3100/api/prom/push)")
	logPerSec   = flag.Int64("logps", 500, "The total amount of log per second to generate.(default 500)")
	randomPanic = flag.Duration("panic-after", 0, "generate random panics")
	useJson     = flag.Bool("json", false, "use json payload(default false)")
)

func init() {
	http.Handle("/metrics", promhttp.Handler())
	flag.Parse()
}

func main() {
	go func() {
		_ = http.ListenAndServe(":2112", nil)
	}()
	if *randomPanic != 0 {
		go func() {
			time.Sleep(time.Duration(rand.Int63n(int64(*randomPanic))))
			thedeepstack()
		}()
	}
	host, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	logViaAPI(*apiURL, host)
}

func logViaAPI(apiURL string, hostname string) {
	cfg, err := loki.NewDefaultConfig(apiURL)
	if err != nil {
		panic(err)
	}
	c, err := loki.New(cfg)
	if err != nil {
		panic(err)
	}

	defer c.Stop()

	ticker := time.NewTicker(time.Second / time.Duration(*logPerSec))
	defer ticker.Stop()
	for {
		<-ticker.C
		_ = c.Handle(
			model.LabelSet{
				"hostname":  model.LabelValue(hostname),
				"component": randComponent(),
				"service":   randService(),
				"app":       "logger",
			}, time.Now(), NewRFC3164Log(time.Now()))
	}
}

// NewRFC3164Log creates a log string with syslog (RFC3164) format
func NewRFC3164Log(t time.Time) string {
	return fmt.Sprintf(
		RFC3164Log,
		gofakeit.Number(0, 191),
		t.Format(RFC3164),
		strings.ToLower(gofakeit.Username()),
		gofakeit.Word(),
		gofakeit.Number(1, 10000),
		gofakeit.HackerPhrase(),
	)
}

func thedeepstack() {
	reallydeep()
}

func reallydeep() {
	deepdown()
}

func deepdown() {
	switch rand.Intn(2) {
	case 1:
		i := 0
		fmt.Print(10 / i)
	default:
		panic("File read error: open /proc/diskstats: too many open files")
	}
}

func randComponent() model.LabelValue {
	return components[rand.Intn(5)]
}

func randService() model.LabelValue {
	return services[rand.Intn(6)]
}

var components = []model.LabelValue{
	"devopsend",
	"fullstackend",
	"frontend",
	"everything-else",
	"backend",
}

var services = []model.LabelValue{
	"potatoes-cart",
	"phishing",
	"stateless-database",
	"random-policies-generator",
	"cookie-jar",
	"distributed-unicorn",
}
