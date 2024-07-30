package subdomain

import (
	"net"
	"net/http"
	h "net/http"
	"net/http/httptrace"
	"strings"
	"time"

	"github.com/datahearth/ddnsclient/pkg/utils"
	"github.com/sirupsen/logrus"
)

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
	var remoteAddr string
	logger := sd.logger.WithField("component", "retrieve-subdomain-ip")

	// * create HEAD request
	req, err := http.NewRequest("HEAD", "https://"+sd.subdomainAddr, nil)
	if err != nil {
		return err
	}
	// * create a trace to get server remote address
	trace := &httptrace.ClientTrace{
		GotConn: func(gci httptrace.GotConnInfo) {
			remoteAddr = gci.Conn.RemoteAddr().String()
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	// * create a client to perform the request
	client := new(h.Client)
	client.Timeout = 5 * time.Second

	// * perform the request
	resp, err := client.Do(req)
	if err != nil {
		// todo: ignoring errors is bad. Implement a solution to scrape 100% of the time remote addr
		logger.WithError(err).WithFields(logrus.Fields{
			"subdomain": sd.subdomainAddr,
		}).Errorln(utils.ErrHeadRemoteIP.Error())
		sd.ip = ""
		return nil
	}
	if resp.StatusCode != 200 && remoteAddr == "" {
		logger.WithFields(logrus.Fields{
			"status-code": resp.StatusCode,
			"subdomain":   sd.subdomainAddr,
		}).Errorln(utils.ErrWrongStatusCode.Error())
		return utils.ErrWrongStatusCode
	}

	// * check if remote address contains a port
	if strings.Contains(remoteAddr, ":") {
		remoteAddr, _, err = net.SplitHostPort(remoteAddr)
		if err != nil {
			logger.WithError(err).WithField("remote-address", remoteAddr).Errorln(utils.ErrSplitAddr.Error())
			return utils.ErrSplitAddr
		}
	}

	sd.ip = remoteAddr

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
