package quic_chan

type RecvData struct {
	SessionID uint64
	StreamID  uint64
	Data      []byte
}
