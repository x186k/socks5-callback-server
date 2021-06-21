package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/txthinking/socks5"
)

var _ = getDefRouteIntfAddrIPv6

func getDefRouteIntfAddrIPv6() net.IP {
	const googleDNSIPv6 = "[2001:4860:4860::8888]:8080" // not important, does not hit the wire
	cc, err := net.Dial("udp6", googleDNSIPv6)          // doesnt send packets
	if err == nil {
		cc.Close()
		return cc.LocalAddr().(*net.UDPAddr).IP
	}
	return nil
}

//from http://github.com/x186k/sfu1
func getDefRouteIntfAddrIPv4() net.IP {
	const googleDNSIPv4 = "8.8.8.8:8080"       // not important, does not hit the wire
	cc, err := net.Dial("udp4", googleDNSIPv4) // doesnt send packets
	if err == nil {
		cc.Close()
		return cc.LocalAddr().(*net.UDPAddr).IP
	}
	return nil
}

type DefaultHandle struct {
}

func (h *DefaultHandle) UDPHandle(*socks5.Server, *net.UDPAddr, *socks5.Datagram) error {
	return fmt.Errorf("no udp support")
}

// TCPHandle auto handle request. You may prefer to do yourself.
func (h *DefaultHandle) TCPHandle(s *socks5.Server, c *net.TCPConn, r *socks5.Request) error {

	println(r.Cmd,c.RemoteAddr().String(),r.Address())
	if r.Cmd == socks5.CmdConnect {
		rc, err := r.Connect(c)
		if err != nil {
			return err
		}
		defer rc.Close()
		go func() {
			var bf [1024 * 2]byte
			for {
				if s.TCPTimeout != 0 {
					if err := rc.SetDeadline(time.Now().Add(time.Duration(s.TCPTimeout) * time.Second)); err != nil {
						return
					}
				}
				i, err := rc.Read(bf[:])
				if err != nil {
					return
				}
				if _, err := c.Write(bf[0:i]); err != nil {
					return
				}
			}
		}()
		var bf [1024 * 2]byte
		for {
			if s.TCPTimeout != 0 {
				if err := c.SetDeadline(time.Now().Add(time.Duration(s.TCPTimeout) * time.Second)); err != nil {
					return nil
				}
			}
			i, err := c.Read(bf[:])
			if err != nil {
				return nil
			}
			if _, err := rc.Write(bf[0:i]); err != nil {
				return nil
			}
		}
	}
	if r.Cmd == socks5.CmdUDP {
		caddr, err := r.UDP(c, s.ServerAddr)
		if err != nil {
			return err
		}
		ch := make(chan byte)
		defer close(ch)
		s.AssociatedUDP.Set(caddr.String(), ch, -1)
		defer s.AssociatedUDP.Delete(caddr.String())
		_, err = io.Copy(io.Discard, c)
		if err != nil {
			return err
		}
		if socks5.Debug {
			log.Printf("A tcp connection that udp %#v associated closed\n", caddr.String())
		}
		return nil
	}
	return socks5.ErrUnsupportCmd
}

func main() {

	myipv4 := getDefRouteIntfAddrIPv4()
	if myipv4 != nil {

		server := myipv4.String() + ":9999"

		println(server)

		s, err := socks5.NewClassicServer(server, myipv4.String(), "", "", 10, 10)
		if err != nil {
			panic(err)
		}
		var z *DefaultHandle
		err = s.ListenAndServe(z)
		if err != nil {
			panic(err)
		}
	}

}
