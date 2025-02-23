package rdp

import (
	"context"
	"fmt"
	"time"

	"github.com/praetorian-inc/fingerprintx/pkg/plugins"
	"github.com/praetorian-inc/fingerprintx/pkg/plugins/services/rdp"
	"github.com/khulnasoft-lab/vulmap/pkg/protocols/common/protocolstate"
)

// RDPClient is a client for rdp servers
type RDPClient struct{}

type IsRDPResponse struct {
	IsRDP bool
	OS    string
}

// IsRDP checks if the given host and port are running rdp server.
//
// If connection is successful, it returns true.
// If connection is unsuccessful, it returns false and error.
//
// The Name of the OS is also returned if the connection is successful.
func (c *RDPClient) IsRDP(host string, port int) (IsRDPResponse, error) {
	resp := IsRDPResponse{}

	timeout := 5 * time.Second
	conn, err := protocolstate.Dialer.Dial(context.TODO(), "tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return resp, err
	}
	defer conn.Close()

	server, isRDP, err := rdp.DetectRDP(conn, timeout)
	if err != nil {
		return resp, err
	}
	if !isRDP {
		return resp, nil
	}
	resp.IsRDP = true
	resp.OS = server
	return resp, nil
}

type CheckRDPAuthResponse struct {
	PluginInfo *plugins.ServiceRDP
	Auth       bool
}

// CheckRDPAuth checks if the given host and port are running rdp server
// with authentication and returns their metadata.
func (c *RDPClient) CheckRDPAuth(host string, port int) (CheckRDPAuthResponse, error) {
	resp := CheckRDPAuthResponse{}

	timeout := 5 * time.Second
	conn, err := protocolstate.Dialer.Dial(context.TODO(), "tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return resp, err
	}
	defer conn.Close()

	pluginInfo, auth, err := rdp.DetectRDPAuth(conn, timeout)
	if err != nil {
		return resp, err
	}
	if !auth {
		return resp, nil
	}
	resp.Auth = true
	resp.PluginInfo = pluginInfo
	return resp, nil
}
