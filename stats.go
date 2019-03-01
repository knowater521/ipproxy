package ipproxy

import (
	"sync/atomic"
	"time"
)

func (p *proxy) trackStats() {
	ticker := time.NewTicker(p.opts.StatsInterval)
	for {
		select {
		case <-p.closeCh:
			return
		case <-ticker.C:
			numTCPOrigins, numTCPClients, numUDPClients := p.ConnCounts()
			log.Debugf("TCP Origins: %v   TCP Clients: %v    UDP Clients: %v", numTCPOrigins, numTCPClients, numUDPClients)
			log.Debugf("Accepted Packets: %d    Rejected Packets: %d", p.AcceptedPackets(), p.RejectedPackets())
		}
	}
}

func (p *proxy) ConnCounts() (numTCPOrigins int, numTCPClients int, numUDPClients int) {
	p.tcpOriginsMx.Lock()
	origins := make([]*origin, 0, len(p.tcpOrigins))
	for _, o := range p.tcpOrigins {
		origins = append(origins, o)
	}
	p.tcpOriginsMx.Unlock()
	numTCPOrigins = len(origins)
	for _, o := range origins {
		numTCPClients += o.numClients()
	}

	p.udpConnTrackMx.Lock()
	numUDPClients = len(p.udpConnTrack)
	p.udpConnTrackMx.Unlock()

	return
}

func (p *proxy) acceptedPacket() {
	atomic.AddInt64(&p.acceptedPackets, 1)
}

func (p *proxy) AcceptedPackets() int {
	return int(atomic.LoadInt64(&p.acceptedPackets))
}

func (p *proxy) rejectedPacket() {
	atomic.AddInt64(&p.rejectedPackets, 1)
}

func (p *proxy) RejectedPackets() int {
	return int(atomic.LoadInt64(&p.rejectedPackets))
}
