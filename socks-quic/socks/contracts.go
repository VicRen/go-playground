package socks

const (
	socks5Version = uint8(5)
)

const (
	ConnectCommand   = uint8(1)
	BindCommand      = uint8(2)
	AssociateCommand = uint8(3)
	ipv4Address      = uint8(1)
	fqdnAddress      = uint8(3)
	ipv6Address      = uint8(4)

	lenAddrType      = uint8(1)
	lenIPv4          = uint8(4)
	lenIPv6          = uint8(16)
	lenPort          = uint8(2)
	lenAssociateRSV  = uint8(2)
	lenAssociateFRAG = uint8(1)

	minDatagramLen = 4
)

const (
	successReply uint8 = iota
	serverFailure
	ruleFailure
	networkUnreachable
	hostUnreachable
	connectionRefused
	ttlExpired
	commandNotSupported
	addrTypeNotSupported
)
