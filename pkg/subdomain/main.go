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

	logger = logger.WithField("pkg", "subdomain")

	return &subdomain{
		logger:        logger,
		subdomainAddr: subdomainAddr,
	}, nil
}

// RetrieveSubdomainIP will retrieve the subdomain IP with a HEAD request
func (sd *subdomain) retrieveSubdomainIP() error {
	logger := sd.logger.WithField("component", "retrieve-subdomain-ip")

	resp, err := h.Head(sd.subdomainAddr)
	if err != nil {
		logger.WithError(err).WithField("subdomain", sd.subdomainAddr).Errorln(utils.ErrHeadRemoteIP.Error())
		return utils.ErrHeadRemoteIP
	}
	if resp.StatusCode != 200 {
		logger.WithField("status-code", resp.StatusCode).Errorln(utils.ErrWrongStatusCode.Error())
		return utils.ErrWrongStatusCode
	}

	host, _, err := net.SplitHostPort(resp.Request.RemoteAddr)
	if err != nil {
		logger.WithError(err).WithField("remote-address", resp.Request.RemoteAddr).Errorln()
		return utils.ErrSplitAddr
	}

	sd.ip = host

	return nil
}

// CheckIPAddr will compare the srvIP passed in parameter and the subIP retrieved from the head request
func (sd *subdomain) CheckIPAddr(srvIP string) (bool, error) {
	if err := sd.retrieveSubdomainIP(); err != nil {
		sd.logger.WithError(err).WithField("component", "check-ip-address").Errorln("failed to retrieve subdomain ip address")
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
