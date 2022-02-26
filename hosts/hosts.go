package hosts

import (
	"crypto/tls"
	"crypto/x509"
	log "github.com/sirupsen/logrus"
	"net"
	"strings"
	"sync"
	"time"
)

var (
	info   []*Info
	lock   sync.Mutex
	done   = make(chan bool)
	logger = log.WithFields(log.Fields{
		"package": "hosts",
	})
)

// GetHosts returns all configured hosts
func GetHosts() []*Info {
	lock.Lock()
	defer lock.Unlock()
	return info
}

// UpdateEvery runs an update on all hosts every interval
func UpdateEvery(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				logger.Info("updating hosts")
				updateAllHosts()
			}
		}
	}()
}

func updateAllHosts() {
	lock.Lock()
	defer lock.Unlock()
	for _, i := range info {
		logger.WithFields(log.Fields{
			"host": i.HostString,
		}).Info("retrieving certs")
		i.GetCerts(5 * time.Second)
	}
}

// Info provides information for a host
type Info struct {
	HostString string
	Certs      []*x509.Certificate
}

// AddHost adds a new host to watch
func AddHost(hostString string) {
	if !strings.Contains(hostString, ":") {
		hostString = hostString + ":443"
	}
	info = append(info, &Info{
		HostString: hostString,
	})
}

// AddHosts adds a set of hosts to watch
func AddHosts(hostStrings []string) {
	for _, h := range hostStrings {
		AddHost(h)
	}
}

// GetCerts retrieves certs for the configure Host
func (i *Info) GetCerts(timeout time.Duration) error {
	dialer := &net.Dialer{Timeout: timeout}
	conn, err := tls.DialWithDialer(dialer, "tcp", i.HostString,
		&tls.Config{
			InsecureSkipVerify: true,
		})
	if err != nil {
		return err
	}

	defer conn.Close()

	if err := conn.Handshake(); err != nil {
		return err
	}

	peerCerts := conn.ConnectionState().PeerCertificates
	i.Certs = make([]*x509.Certificate, 0, len(peerCerts))
	logger.WithFields(log.Fields{
		"num_certs": len(peerCerts),
		"host":      i.HostString,
	}).Info("found certs")

	for _, cert := range peerCerts {
		if cert.IsCA {
			continue
		}
		i.Certs = append(i.Certs, cert)
	}

	return nil
}
