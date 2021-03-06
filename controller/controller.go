package controller

import (
	"flag"
	"github.com/braintree/manners"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/potix/gox/configure"
	"github.com/potix/gox/downloader"
	"github.com/potix/gox/filter"
	"github.com/potix/gox/router"
	"os"
	"path/filepath"
)

type Controller struct {
	configPath    string
	goxHome       string
	availablePath string
	enablePath    string
	address       string
	port          string
	downloader    *downloader.Downloader
	filter        *filter.Filter
	router        *router.Router
	manners       *manners.GracefulServer
	engine        *gin.Engine
}

func newController(configPath string, goxHome string, availablePath string, enablePath string,
	address string, port string, downloader *Downloader, filter *Filter, router *Router) *Controller {
	return &Controller{
		goxHome:       config.GoxHome,
		availablePath: availablePath,
		enablePath:    enablePath,
		address:       config.ControllerAddress,
		port:          config.ControllerPort,
		downloader:    downloader,
		filter:        filter,
		router:        router,
	}
}

func (c *controller) start() (err error) {
	staticPath = filepath.Join(goxHome, "www", "static")
	if _, err := os.Stat(staticPath); err != nil {
		return err
	}
	templatePath = filepath.Join(goxHome, "www", "template")
	if _, err := os.Stat(templatePath); err != nil {
		return err
	}
	c.engine = gin.Default()
	c.engine.Static("/asset", staticPath)
	c.engine.LoadHTMLGlob(filepath.Join(templatePath, "*"))
	c.engine.GET("/", c.getRoot)
	c.engine.POST("/gox", c.postGox)
	c.engine.GET("/gox/config", c.getGoxConfig)
	c.engine.PUT("/gox/config", c.putGoxConfig)
	c.engine.GET("/plugins", c.getPlugins)
	c.engine.POST("/plugins", c.postPlugins)
	c.engine.DELETE("/plugins/:pluginName", c.deletePlugin)
	c.engine.POST("/plugins/:pluginName", c.postPlugin)
	c.engine.GET("/plugins/:pluginName/help", c.getPluginHelp)
	c.engine.GET("/plugins/:pluginName/commands", c.getPluginCommands)
	c.engine.GET("/plugins/:pluginName/events", c.getPluginEvents)
	c.engine.GET("/plugins/:pluginName/config", c.getPluginConfig)
	c.engine.PUT("/plugins/:pluginName/config", c.putPluginConfig)
	c.manners = manners.NewWithServer(&http.Server{
		Addr:    ":8080",
		Handler: c.engine,
	})
	err = c.manners.ListenAndServe()
	if err != nil {
		return err
	}
}

func makePluginPath(goxHoma stringe) (availablePath string, enablePath string, err error) {
	err = os.MkdirAll(goxHome, 0755)
	if err != nil {
		return nil, nil, err
	}
	availablePath = filepath.Join(goxHome, "plugins", "available")
	err = os.MkdirAll(availablePath, 0755)
	if err != nil {
		return nil, nil, err
	}
	enablePath = filepath.Join(goxHome, "plugins", "enable")
	err = os.MkdirAll(enalePath, 0755)
	if err != nil {
		return nil, nil, err
	}
	return availablePath, enablePath, err
}

func Run() (err error) {
	configPath := flag.String("c", "/etc/gox.conf", "config file papth")
	flag.Parse()
	config, err := configre.Load(configPath)
	if err != nil {
		return err
	}
	availablePath, enablePath, err := makePluginPath(config.GoxHome)
	if err != nil {
		return err
	}
	downloader := downloader.NewDownloader(availablePath)
	filter := filter.NewFilter(enablePath, config.FilterRules)
	hookFunc := filter.GetHookFunc()
	router := router.NewRouter(enablePath, hookFunc)
	controller := newController(configPath, config.GoxHome, availablePath, enablePath,
		config.ControllerAddress, config.ControllerPort, downloader, filter, router)
	err = contoller.start()
	if err != nil {
		return err
	}
}
