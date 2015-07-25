package ipmsg

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"time"
)

type IPMSG struct {
	ClientData ClientData
	Conn       *net.UDPConn
	Conf       *IPMSGConfig
	Broadcast  []net.IP
	PacketNum  int
}

type IPMSGConfig struct {
	NickName  string
	GroupName string
	UserName  string
	HostName  string
	Port      int
}

const (
	DefaultPort int = 2425
	Buflen      int = 65535
)

func NewIPMSGConf() *IPMSGConfig {
	conf := &IPMSGConfig{
		Port: DefaultPort,
	}
	return conf
}

func NewIPMSG(conf *IPMSGConfig) (*IPMSG, error) {
	ipmsg := &IPMSG{
		PacketNum: 0,
	}
	ipmsg.Conf = conf
	// UDP server
	service := fmt.Sprintf(":%v", conf.Port)
	//fmt.Println("service = ", service)
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	if err != nil {
		return ipmsg, err
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return ipmsg, err
	}
	ipmsg.Conn = conn
	return ipmsg, err
}

func (ipmsg *IPMSG) Close() error {
	conn := ipmsg.Conn
	if conn == nil {
		err := errors.New("Conn is not defined")
		return err
	}
	err := conn.Close()
	return err
}

func (ipmsg *IPMSG) SendMSG(msg string, addr *net.UDPAddr) error {
	clientdata := NewClientData("", addr)
	clientdata.Version = 1
	clientdata.PacketNum = ipmsg.GetNewPacketNum()
	clientdata.User = "user"
	clientdata.Host = "host"
	clientdata.Command = BR_ENTRY
	clientdata.Option = msg
	//pp.Println("clientdata.String=", clientdata.String())

	conn := ipmsg.Conn
	_, err := conn.WriteToUDP([]byte(clientdata.String()), addr)
	if err != nil {
		return err
	}
	return nil
}

func (ipmsg *IPMSG) RecvMSG() (*ClientData, error) {
	var buf [Buflen]byte
	conn := ipmsg.Conn
	_, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		return nil, err
	}
	trimmed := bytes.Trim(buf[:], "\x00")
	clientdata := NewClientData(string(trimmed[:]), addr)
	return clientdata, nil
	//return string(trimmed[:]), addr, nil
}

// convert net.Addr to net.UDPAddr
func (ipmsg *IPMSG) UDPAddr() (*net.UDPAddr, error) {
	conn := ipmsg.Conn
	if conn == nil {
		err := errors.New("Conn is not defined")
		return nil, err
	}
	addr := conn.LocalAddr()
	network := addr.Network()
	str := addr.String()
	udpAddr, err := net.ResolveUDPAddr(network, str)
	return udpAddr, err
}

func (ipmsg *IPMSG) AddBroadCast(ip net.IP) {
	bc := ipmsg.Broadcast
	bc = append(bc, ip)
	ipmsg.Broadcast = bc
}

func (ipmsg *IPMSG) GetNewPacketNum() int {
	ipmsg.PacketNum++
	return int(time.Now().Unix()) + ipmsg.PacketNum
}
