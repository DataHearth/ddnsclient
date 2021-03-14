package http

import (
	"io/ioutil"
	"net"
	h "net/http"

	"github.com/datahearth/ddnsclient/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// HTTP is the base interface to interact with websites
type HTTP interface {
	CheckIPAddr(addr string) (bool, error)
	GetServerIP() (string, error)
	GetSubdomainIP(addr string) (string, error)
	RetrieveSubdomainIP(addr string) error
	RetrieveServerIP() error
}

type http struct {
	logger      logrus.FieldLogger
	serverIP    string
	subdomainIP string
	webIP       string
}

// NewHTTP instanciate a new http implementation
func NewHTTP(logger logrus.FieldLogger) (HTTP, error) {
	if logger == nil {
		return nil, utils.ErrNilLogger
	}

	logger = logger.WithField("package", "http")
	webIP := viper.GetString("web-ip")

	if webIP == "" {
		logger.Infoln("web-ip field not set in config. Using default")
		webIP = "http://dynamicdns.park-your-domain.com/getip"
	}

	return &http{
		logger: logger,
		webIP:  webIP,
	}, nil
}

// RetrieveServerIP will use the defined web-ip service to get the server public address and save it to the struct
func (http *http) RetrieveServerIP() error {
	resp, err := h.Get(http.webIP)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return utils.ErrWrongStatusCode
	}
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	http.serverIP = string(d)

	return nil
}

// GetServerIP will return the IP addr of the server.
// If you haven't triggered the function before manually or with CheckIPAddr,
// it'll trigger the RetrieveServerIP function and save it. Then, it'll return the IP
func (http *http) GetServerIP() (string, error) {
	if http.serverIP == "" {
		if err := http.RetrieveServerIP(); err != nil {
			return "", err
		}
	}

	return http.serverIP, nil
}

// RetrieveSubdomainIP will retrieve the subdomain IP with a HEAD request then save it to the struct
func (http *http) RetrieveSubdomainIP(addr string) error {
	resp, err := h.Head(addr)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return utils.ErrWrongStatusCode
	}

	h, _, err := net.SplitHostPort(resp.Request.RemoteAddr)
	if err != nil {
		return err
	}

	http.subdomainIP = h

	return nil
}

// GetSubdomainIP will return the IP addr of the subdomain.
// If you haven't triggered the function before manually or with CheckIPAddr,
// it'll trigger the RetrieveSubdomainIP function and save it. Then, it'll return the IP
func (http *http) GetSubdomainIP(addr string) (string, error) {
	if http.subdomainIP == "" {
		if err := http.RetrieveSubdomainIP(addr); err != nil {
			return "", err
		}
	}

	return http.subdomainIP, nil
}

// CheckIPAddr will get both ip addr (subdomain from addr and server) and compare it.
// If one of them is missing, it'll retrieve it and then save it
func (http *http) CheckIPAddr(addr string) (bool, error) {
	srvIP, err := http.GetServerIP()
	if err != nil {
		return false, err
	}

	subIP, err := http.GetSubdomainIP(addr)
	if err != nil {
		return false, err
	}

	if srvIP != subIP {
		return false, nil
	}

	return true, nil
}
