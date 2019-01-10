package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
	_ "github.com/infinivision/vectodb/cmd/vectodblite_cluster/docs" // docs is generated by Swag CLI, you have to import it.
	log "github.com/sirupsen/logrus"
	ginSwagger "github.com/swaggo/gin-swagger"                // gin-swagger middleware
	swaggerFiles "github.com/swaggo/gin-swagger/swaggerFiles" // swagger embed files
)

// @title VectoDBLite Cluster API
// @version 1.0
// @description This is a VectoDBLite cluster server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://github.com/infinivision/vectodb/issues
// @contact.email yuzhichang@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /api/v1

// report version and Git SHA, inspired by github.com/coreos/etcd/version/version.go
var (
	Version = "1.0-SNAPSHOT"
	// GitSHA and BuildTime will be set during build
	GitSHA    = "Not provided (use ./build.sh instead of go build)"
	BuildTime = "Not provided (use ./build.sh instead of go build)"
)

func parseConfig() (conf *ControllerConf) {
	conf = NewControllerConf()
	flag.StringVar(&conf.ListenAddr, "listen-addr", conf.ListenAddr, "Addr: listen address")
	flag.StringVar(&conf.EtcdAddr, "etcd-addr", conf.EtcdAddr, "Addr: etcd address")
	flag.StringVar(&conf.RedisAddr, "redis-addr", conf.RedisAddr, "Addr: redis address")
	flag.IntVar(&conf.Dim, "dim", conf.Dim, "VectoDBLite dimension")
	flag.Float64Var(&conf.DisThr, "distance-threshold", conf.DisThr, "VectoDBLite distance threshold")
	flag.IntVar(&conf.SizeLimit, "size-limit", conf.SizeLimit, "VectoDBLite size limit")
	flag.IntVar(&conf.BalanceInterval, "balance-interval", conf.BalanceInterval, "Time interval (in seconds) to balance the cluster load")

	flag.StringVar(&conf.EurekaAddr, "eureka-addr", conf.EurekaAddr, "eureka server address list, seperated by comma.")
	flag.StringVar(&conf.EurekaApp, "eureka-app", conf.EurekaApp, "VectoDBLite cluster service name which will be registered with eureka.")

	isDebug := flag.Bool("debug", false, "Set log level to debug")
	showVer := flag.Bool("version", false, "Show version and quit.")
	flag.Parse()
	if *showVer {
		fmt.Printf("vectodblite_cluster Version: %s\n", Version)
		fmt.Printf("Git SHA: %s\n", GitSHA)
		fmt.Printf("BuildTime: %s\n", BuildTime)
		fmt.Printf("Go Version: %s\n", runtime.Version())
		fmt.Printf("Go OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}
	if *isDebug {
		log.SetLevel(log.DebugLevel)
		gin.SetMode(gin.DebugMode)
	} else {
		log.SetLevel(log.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
	}
	return
}

func main() {
	conf := parseConfig()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctl := NewController(conf, ctx)
	r := gin.Default()
	r.POST("/api/v1/add", ctl.HandleAdd)
	r.POST("/api/v1/search", ctl.HandleSearch)
	r.POST("/mgmt/v1/acquire", ctl.HandleAcquire)
	r.POST("/mgmt/v1/release", ctl.HandleRelease)
	r.GET("/status", ctl.HandleStatus)
	r.GET("/health", ctl.HandleHealth)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(conf.ListenAddr)
}
