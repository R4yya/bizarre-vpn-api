package coreErrors

import "errors"

type ErrorVpnServer = error

var (
	ErrorVpnServerInvalidIp                   ErrorVpnServer = errors.New("invalid IP")
	ErrorVpnServerInvalidPort                 ErrorVpnServer = errors.New("invalid port")
	ErrorVpnServerNameIsTooSmall              ErrorVpnServer = errors.New("name too small")
	ErrorVpnServerUniqueConstraint            ErrorVpnServer = errors.New("UNIQUE constraint failed")
	ErrorVpnServerUniqueConstraintHostAndPort ErrorVpnServer = errors.New("UNIQUE constraint failed: vpn_servers.adapter_host, vpn_servers.adapter_port")
	ErrorVpnServerUniqueConstraintName        ErrorVpnServer = errors.New("UNIQUE constraint failed: vpn_servers.name")
)
