package ipmsg

import (
	"errors"
	"net"
	"testing"
)

func TestGetNewPackNum(t *testing.T) {
	conf := NewIPMSGConf()
	ipmsg, err := NewIPMSG(conf)
	if err != nil {
		t.Errorf("ipmsg error is not nil '%v'", err)
	}
	defer ipmsg.Close()

	if 0 != ipmsg.PacketNum {
		t.Errorf("ipmsg.PacketNum should be 0 but '%v'", ipmsg.PacketNum)
	}
	num := ipmsg.GetNewPacketNum()
	if num == 0 {
		t.Errorf("ipmsg.GetNewPacketNum returns 0")
	}
	if 1 != ipmsg.PacketNum {
		t.Errorf("ipmsg.PacketNum should be 1 but '%v'", ipmsg.PacketNum)
	}
}

func getv4loopback() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagLoopback != 0 {
			addrs, err := iface.Addrs()
			if err != nil {
				return "", err
			}
			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}
				ip = ip.To4()
				if ip == nil {
					continue
				}
				if ip.IsLoopback() {
					return ip.String(), nil
				}
			}
		}
	}
	return "", errors.New("you do not have loopback?")
}

func TestNewIPMSG(t *testing.T) {
	conf := NewIPMSGConf()
	conf.NickName = "testuser"
	conf.GroupName = "testgroup"
	conf.UserName = "testuser"
	conf.HostName = "testhost"
	laddr, err := getv4loopback()
	if err != nil {
		t.Errorf("loopback not found '%v'", err)
	}
	conf.Local = laddr
	client, err := NewIPMSG(conf)
	if err != nil {
		t.Errorf("client error is not nil '%v'", err)
	}
	defer client.Close()

	serverConf := NewIPMSGConf()
	serverConf.Port = 12425
	serverConf.Local = laddr
	server, err := NewIPMSG(serverConf)
	if err != nil {
		t.Errorf("server error is not nil '%v'", err)
	}
	defer server.Close()

	saddr, err := server.UDPAddr()
	if err != nil {
		t.Errorf("failed to resolve to UDP '%v'", err)
	}

	// client sends message to server
	testmsg := "hogehoge"
	err = client.SendMSG(saddr, testmsg, BR_ENTRY)
	if err != nil {
		t.Errorf("client.SendMSG return error '%v'", err)
	}
	// server receives message from client
	received, err := server.RecvMSG()
	//pp.Println("received = ", received)
	if err != nil {
		t.Errorf("server.RecvMSG return error '%v'", err)
	}

	if testmsg != received.Option {
		//if testmsg != received {
		t.Errorf("received is not much to sent msg")
	}
	//pp.Println("received  = ", received)
}
