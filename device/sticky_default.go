//go:build !linux

package device

import (
	"github.com/tailscale/wireguard-go-awg/conn"
	"github.com/tailscale/wireguard-go-awg/rwcancel"
)

func (device *Device) startRouteListener(bind conn.Bind) (*rwcancel.RWCancel, error) {
	return nil, nil
}
