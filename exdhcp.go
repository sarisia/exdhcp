package exdhcp

import (
	"context"
	"encoding/csv"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/nclient4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ExdhcpClient struct {
	InterfaceName string
	ServerAddress string

	dhcpcli *nclient4.Client
	log     *zap.SugaredLogger

	leases []*nclient4.Lease
}

func New(ifname, serverAddr string, verbose bool) (*ExdhcpClient, error) {
	opts := make([]nclient4.ClientOpt, 0)

	loggerConfig := zap.NewDevelopmentConfig()
	if !verbose {
		loggerConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}
	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, err
	}
	log := logger.Sugar()

	if serverAddr != "" {
		log.Debugf("using fixed server address: %s", serverAddr)
		addr := net.UDPAddr{
			IP:   net.IP(serverAddr),
			Port: 67,
		}
		opts = append(opts, nclient4.WithServerAddr(&addr))
	}

	dhcpcli, err := nclient4.New(ifname, opts...)
	if err != nil {
		log.Errorf("initialization failed: %s", err)
		return nil, err
	}

	log.Sync()

	return &ExdhcpClient{
		InterfaceName: ifname,
		ServerAddress: serverAddr,
		dhcpcli:       dhcpcli,
		log:           log,
	}, nil
}

func (e *ExdhcpClient) Start(ctx context.Context, numTries int, timeout int, release bool, exportCSV bool) {
	e.log.Info("starting...")

	e.leases = make([]*nclient4.Lease, 0)

	if numTries > 0 {
		for i := 0; i < numTries; i++ {
			lease, err := e.doRequest(ctx, timeout, release)
			if lease != nil {
				e.leases = append(e.leases, lease)
			}
			if err != nil {
				e.log.Errorf("dhcp failed: %s", err)
				break
			}
		}
	} else {
		for {
			lease, err := e.doRequest(ctx, timeout, release)
			if lease != nil {
				e.leases = append(e.leases, lease)
			}
			if err != nil {
				e.log.Errorf("dhcp failed: %s", err)
				break
			}
		}
	}

	if exportCSV {
		e.exportCSV()
	}

	e.log.Infof("Summary: leased %d addresses!", len(e.leases))
	e.log.Sync()
}

func (e *ExdhcpClient) doRequest(parent context.Context, timeout int, release bool) (*nclient4.Lease, error) {
	ctx, cancel := context.WithTimeout(parent, time.Duration(timeout*int(time.Second)))
	defer cancel()

	hwaddr := randomMAC()
	hwaddrMod := dhcpv4.WithHwAddr(hwaddr)
	broadcastMod := dhcpv4.WithBroadcast(true)

	e.log.Infof("try lease (ClientHWAddr=%s)", hwaddr.String())
	lease, err := e.dhcpcli.Request(ctx, hwaddrMod, broadcastMod)
	if err != nil {
		return nil, err
	}
	e.log.Infof("leased! (IP=%s, ClientHWAddr=%s)", lease.ACK.YourIPAddr.String(), lease.ACK.ClientHWAddr.String())
	e.log.Debugf("Offer summary: \n%s", lease.Offer.Summary())
	e.log.Debugf("ACK summary: \n%s", lease.ACK.Summary())

	if release {
		err = e.dhcpcli.Release(lease)
		if err != nil {
			e.log.Errorf("dhcp release failed: %s", err)
			return lease, err
		}
	}

	return lease, err
}

func (e *ExdhcpClient) exportCSV() {
	filename := fmt.Sprintf("exdhcp-%s.csv", time.Now().Format("20060102150405"))
	f, err := os.Create(filename)
	if err != nil {
		e.log.Errorf("failed to open file (filename=%s): %s", filename, err)
		return
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	// write header
	writer.Write([]string{"IP", "clientHWAddr"})

	for _, l := range e.leases {
		writer.Write([]string{l.ACK.YourIPAddr.String(), l.ACK.ClientHWAddr.String()})
	}

	writer.Flush()
	e.log.Infof("CSV saved as %s", filename)
}

func randomMAC() net.HardwareAddr {
	hwaddrBytes := [6]byte{
		0xDE,
		0xAD,
		byte(rand.Intn(0x29)),
		byte(rand.Intn(0x7f)),
		byte(rand.Intn(0xff)),
		byte(rand.Intn(0xff)),
	}
	return net.HardwareAddr(hwaddrBytes[:])
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
