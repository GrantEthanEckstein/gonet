package gonet

import (
	"crypto/tls"
	"log"
	"net"
)

// TLSConnection embeds tls.Conn, extending functionality.
// The main idea is to implement TLSConnection.Send and TLSConnection.Recv
// allowing for more simple websocket operation.
type TLSConnection struct {
	tls.Conn
	Settings         ConnectionConfiguration
	Listener         net.Conn
	Outgoing         *tls.Conn
	ListenerSettings ListenerConfiguration
}

// ConnectionConfiguration is included for
// the express purpose of extensibility.
type ConnectionConfiguration struct {
	TLSConnectionConfiguration tls.Config
	Protocol                   string
	Addr                       string
}

// ListenerConfiguration is also included
// for the express purpose of extensibility.
type ListenerConfiguration struct {
	Cert tls.Certificate
	Addr string
}

// Check is a basic errer parsing solution.
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// NewTLS is the constructor for TLSConnection.
// The primary parameter is the protocol used for the connection.
// The insecure parameter will set the configuration for using a
// self-signed certificate. This configuration can be replaced by using
// TLSConnection.SetConfig.
func NewTLS(network string, insecure bool) TLSConnection {
	c := TLSConnection{
		Conn: tls.Conn{},
		Settings: ConnectionConfiguration{
			Protocol: network,
			Addr:     "",
		},
		ListenerSettings: ListenerConfiguration{},
	}
	if insecure {
		c.Settings.TLSConnectionConfiguration = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: true,
		}
	}
	return c
}

// SetConfig is to more easily access c.Settings.TLSConnectionConfiguration
func (c *TLSConnection) SetConfig(config tls.Config) {
	c.Settings.TLSConnectionConfiguration = &config
}

// Recv listens for data of length `len int` using tls.Listen.Accept
func (c *TLSConnection) Recv(len int, laddr, certFile, keyFile string) []byte {
	c.ListenerSettings.Addr = laddr

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	check(err)

	ln, err := tls.Listen(c.Settings.Protocol, c.ListenerSettings.Addr, c.Settings.TLSConnectionConfiguration)
	check(err)

	conn, err := ln.Accept()
	check(err)

	data := make([]byte, len)
	_, err = conn.Read(data)
	check(err)

	c.Listener = conn

	return data
}

// Send is used to send data using tls.Dial.
func (c *TLSConnection) Send(data []byte, addr string) {
	var err error
	c.Outgoing, err = tls.Dial(c.Settings.Protocol, addr, c.Settings.TLSConnectionConfiguration)
	check(err)

	c.Outgoing.Write(data)
}
