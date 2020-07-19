package conn

import (
	"net"
	"net/url"
	"time"
)

//TCPWithDeadline ...
type TCPWithDeadline struct {
	conn          net.Conn
	readDeadLine  time.Duration
	writeDeadLine time.Duration
}

//Read ...
func (c *TCPWithDeadline) Read(b []byte) (
	n int, err error) {
	c.conn.SetReadDeadline(time.Now().Add(c.readDeadLine))
	return c.conn.Read(b)
}

//Write ...
func (c *TCPWithDeadline) Write(b []byte) (
	n int, err error) {
	c.conn.SetWriteDeadline(time.Now().Add(c.writeDeadLine))
	return c.conn.Write(b)
}

//Close ...
func (c *TCPWithDeadline) Close() error {
	return c.Close()
}

//NewRtspTCPConn create new rtsp tcp connection
func NewRtspTCPConn(
	urlString string,
	dialTimeout, readDeadLine, writeDeadLine time.Duration) (
	*TCPWithDeadline,
	*url.URL,
	error) {
	url, err := url.Parse(urlString)
	if err != nil {
		return nil, nil, err
	}
	if _, _, err := net.SplitHostPort(url.Host); err != nil {
		url.Host = url.Host + ":554"
	}
	dailer := net.Dialer{Timeout: dialTimeout}
	conn, err := dailer.Dial("tcp", url.Host)
	if err != nil {
		return nil, nil, err
	}
	tcpConn := &TCPWithDeadline{conn, readDeadLine, writeDeadLine}
	return tcpConn, url, nil
}
