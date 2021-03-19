package subdomain

import (
	"errors"
	"net"
	"strings"
	"time"

	"github.com/datahearth/ddnsclient/pkg/utils"
	"github.com/sirupsen/logrus"
)

// ErrIpLength is thrown when subdomain no or multiples remote IP address
var ErrIpLenght = errors.New("zero or more than 1 ips have been found")

type (
	PendingSubdomains map[time.Time]Subdomain
	subdomain         struct {
		logger        logrus.FieldLogger
		subdomainAddr string
		ip            string
	}
	Subdomain interface {
		CheckIPAddr(srvIP string) (bool, error)
		GetSubdomainIP() string
		retrieveSubdomainIP() error
		GetSubdomainAddr() string
		SubIsPending(sbs PendingSubdomains) bool
		FindSubdomain(sbs PendingSubdomains) Subdomain
	}
)

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
	ips, err := net.LookupIP(sd.subdomainAddr)
	if err != nil {
		return err
	}

	if len(ips) != 1 {
		return ErrIpLenght
	}

	ip := ips[0].String()
	if strings.Contains(ip, ":") {
		ip, _, err = net.SplitHostPort(ip)
		if err != nil {
			return utils.ErrSplitAddr
		}
	}

	sd.ip = ip

	return nil
}

// CheckIPAddr will compare the server IP and the subdomain IP
func (sd *subdomain) CheckIPAddr(srvIP string) (bool, error) {
	if err := sd.retrieveSubdomainIP(); err != nil {
		return false, err
	}

	if srvIP != sd.ip {
		return false, nil
	}

	return true, nil
}

// GetSubdomainIP returns the subdomain IP
func (sd *subdomain) GetSubdomainIP() string {
	return sd.ip
}

// GetSubdomainAddr returns the subdomain address
func (sd *subdomain) GetSubdomainAddr() string {
	return sd.subdomainAddr
}
