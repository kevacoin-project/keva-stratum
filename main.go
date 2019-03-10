package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"./pool"
	"./stratum"

	"github.com/goji/httpauth"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-version"
	"github.com/yvasiyarov/gorelic"
)

var cfg pool.Config

func startStratum() {
	if cfg.Threads > 0 {
		runtime.GOMAXPROCS(cfg.Threads)
		log.Printf("Running with %v threads", cfg.Threads)
	} else {
		n := runtime.NumCPU()
		runtime.GOMAXPROCS(n)
		log.Printf("Running with default %v threads", n)
	}

	s := stratum.NewStratum(&cfg)
	if cfg.Frontend.Enabled {
		go startFrontend(&cfg, s)
	}
	s.Listen()
}

func startFrontend(cfg *pool.Config, s *stratum.StratumServer) {
	r := mux.NewRouter()
	r.HandleFunc("/stats", s.StatsIndex)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./www/")))
	var err error
	if len(cfg.Frontend.Password) > 0 {
		auth := httpauth.SimpleBasicAuth(cfg.Frontend.Login, cfg.Frontend.Password)
		err = http.ListenAndServe(cfg.Frontend.Listen, auth(r))
	} else {
		err = http.ListenAndServe(cfg.Frontend.Listen, r)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func startNewrelic() {
	// Run NewRelic
	if cfg.NewrelicEnabled {
		nr := gorelic.NewAgent()
		nr.Verbose = cfg.NewrelicVerbose
		nr.NewrelicLicense = cfg.NewrelicKey
		nr.NewrelicName = cfg.NewrelicName
		nr.Run()
	}
}

func readConfig(cfg *pool.Config) {
	configFileName := "config.json"
	if len(os.Args) > 1 {
		configFileName = os.Args[1]
	}
	configFileName, _ = filepath.Abs(configFileName)
	log.Printf("Loading config: %v", configFileName)

	configFile, err := os.Open(configFileName)
	if err != nil {
		log.Fatal("File error: ", err.Error())
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&cfg); err != nil {
		log.Fatal("Config error: ", err.Error())
	}
}

func checkRedisVersion(info string) bool {
	var versionStr string
	parts := strings.Split(info, "\r\n")
	for _, line := range parts {
		if strings.Index(line, ":") != -1 {
			valParts := strings.Split(line, ":")
			if valParts[0] == "redis_version" {
				versionStr = valParts[1]
				break
			}
		}
	}
	if versionStr == "" {
		log.Printf("Could not detect redis version - must be super old or broken")
		return false
	}
	minVersion, _ := version.NewVersion("2.6")
	curVersion, err := version.NewVersion(versionStr)
	if err != nil {
		log.Printf("Could not check redis version: " + err.Error())
		return false
	}
	if curVersion.LessThan(minVersion) {
		log.Printf("You're using redis version %s the minimum required version is 2.6. Follow the damn usage instructions...", versionStr)
		return false
	}
	return true
}

func main() {
	info, err := stratum.RedisClient.Info().Result()
	if err != nil {
		log.Fatal("Cannot start redis server.")
	}
	if !checkRedisVersion(info) {
		os.Exit(1)
	}
	rand.Seed(time.Now().UTC().UnixNano())
	readConfig(&cfg)
	startNewrelic()
	startStratum()
}
