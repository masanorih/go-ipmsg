package ipmsg

import (
	"net"
	"testing"
)

func TestParse(t *testing.T) {
	addr := new(net.Addr)
	client := NewClientData("", addr)
	//client.Parse("1:2:user:host:3:4")
	client.Parse("1:2:user:host:3:nick\x00group\x00")
	if client.Version != 1 {
		t.Errorf("client.Version should be 1 but '%v'", client.Version)
	}
	if client.PacketNum != 2 {
		t.Errorf("client.PacketNum should be 2 but '%v'", client.PacketNum)
	}
	if client.Nick != "nick" {
		t.Errorf("client.Nick should be 'nick' but '%v'", client.Nick)
	}
	if client.Group != "group" {
		t.Errorf("client.Group should be 'group' but '%v'", client.Group)
	}
	if client.NickName() != "nick@group" {
		t.Errorf("client.NickName() should be 'nick@group' but '%v'",
			client.NickName)
	}
}
