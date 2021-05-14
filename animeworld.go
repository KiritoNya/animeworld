package animeworld

import (
	"errors"
	"github.com/tebeka/selenium"
	"io"
	"os"
	"fmt"
)

const BaseUrl string = "https://www.animeworld.tv"
const ArchiveUrl string = BaseUrl + "/az-list/"
const SiteApi string = "/api/"
const episodeInfoApi string = BaseUrl + SiteApi + "episode/info?id="
const AnimeworldBizApi string = "https://www.animeworld.biz/api/source/"

type ServiceConfig struct {
	SeleniumPath    string
	GeckoDriverPath string
	Port            int
	Verbose         bool
	Output          io.Writer
}

type WebDriverConfig struct {
	Capabilities selenium.Capabilities
	UrlService   string
}

var service *selenium.Service
var serviceActive bool

func NewService(cfg *ServiceConfig) (err error) {

	if serviceActive {
		return errors.New("Service already exist")
	}

	opts := []selenium.ServiceOption{
		selenium.StartFrameBuffer(),               // Start an X frame buffer for the browser to run in.
		selenium.GeckoDriver(cfg.GeckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
	}

	if cfg.Output != nil {
		opts = append(opts, selenium.Output(os.Stderr))
	}

	selenium.SetDebug(cfg.Verbose)

	fmt.Println("[ ] Creating selenium service...")
	
	service, err = selenium.NewSeleniumService(cfg.SeleniumPath, cfg.Port, opts...)
	if err != nil {
		return err
	}

	serviceActive = true
	
	fmt.Println("[V] Creating selenium service...")

	return nil
}

func NewDefaultService() error {

	cfg := ServiceConfig{
		SeleniumPath:    "vendor/selenium-server.jar",
		GeckoDriverPath: "vendor/geckodriver",
		Port:            4444,
		Verbose:         false,
		Output:          nil,
	}

	err := NewService(&cfg)
	if err != nil {
		return err
	}

	return nil
}

func NewWebDriver(url string, cfg *WebDriverConfig) (selenium.WebDriver, error) {

	if !serviceActive {
		return nil, errors.New("Service is inactive, read the documentation.")
	}

	wd, err := selenium.NewRemote(cfg.Capabilities, cfg.UrlService)
	if err != nil {
		return nil, err
	}

	if err = wd.Get(url); err != nil {
		return nil, err
	}

	return wd, err
}

func NewDefaultWebDriver(url string) (selenium.WebDriver, error) {

	cfg := WebDriverConfig{
		Capabilities: selenium.Capabilities{"browserName": "firefox"},
		UrlService:   "http://localhost:4444/wd/hub",
	}

	wd, err := NewWebDriver(url, &cfg)

	return wd, err
}

func ServiceIsActive() bool {
	return serviceActive
}
