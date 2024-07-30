package subdomain

import (
	"net"
	h "net/http"

	"github.com/datahearth/ddnsclient/pkg/utils"
	"github.com/sirupsen/logrus"
)

// HTTP is the base interface to interact with websites
type Subdomain interface {
	CheckIPAddr(srvIP string) (bool, error)
	GetSubdomainIP() string
	retrieveSubdomainIP() error
}

type subdomain struct {
	logger        logrus.FieldLogger
	subdomainAddr string
	ip            string
}

// NewSubdomain instanciate a new http implementation
func NewSubdomain(logger logrus.FieldLogger, subdomainAddr string) (Subdomain, error) {
	if logger == nil {
		return nil, utils.ErrNilLogger
	}

	logger = logger.WithField("package", "http")

	return &subdomain{
		logger:        logger,
		subdomainAddr: subdomainAddr,
	}, nil
}

// RetrieveSubdomainIP will retrieve the subdomain IP with a HEAD request
func (sd *subdomain) retrieveSubdomainIP() error {
	resp, err := h.Head(sd.subdomainAddr)
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

	sd.ip = h

	return nil
}

// CheckIPAddr will compare the srvIP passed in parameter and the subIP retrieved from the head request
func (sd *subdomain) CheckIPAddr(srvIP string) (bool, error) {
	if err := sd.retrieveSubdomainIP(); err != nil {
		return false, err
	}

	if srvIP != sd.ip {
		return false, nil
	}

	return true, nil
}

func (sd *subdomain) GetSubdomainIP() string {
	return sd.ip
}
