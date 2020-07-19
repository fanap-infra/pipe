package rtp

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"log"
)

//Rtp handle rtp connection to camera.
type Rtp struct {
	subscriberList map[int]chan []byte
	Reader         *bufio.Reader
	ConnectionType string
}

//NewRTPConnection ...
func NewRTPConnection(reader *bufio.Reader, connectionType string) *Rtp {
	subscriberList := make(map[int]chan []byte)
	rtp := Rtp{Reader: reader, ConnectionType: connectionType, subscriberList: subscriberList}
	rtp.strtReading()
	return &rtp
}

func (conn *Rtp) parseTCPRaw() (rawPkt []byte, channelNum int, err error) {
	if _, err := conn.Reader.ReadBytes('$'); err != nil {
		return nil, -1, err
	}
	header := make([]byte, 3)
	if _, err := io.ReadFull(conn.Reader, header); err != nil {
		return nil, -1, err
	}
	channelNum = int(header[0])
	length := binary.BigEndian.Uint16(header[1:3])
	rawPkt = make([]byte, length)
	n, err := io.ReadFull(conn.Reader, rawPkt)
	if err != nil {
		return nil, -1, err
	}
	if uint16(n) != length {
		return nil, -1, errors.New("Lenght of packet and read doesn't match")
	}
	return rawPkt, channelNum, nil
}

//SubscribeChannel read rtp packet from tcp connection.
func (conn *Rtp) SubscribeChannel(channelNum int) chan []byte {
	ch := make(chan []byte, 500)
	conn.subscriberList[channelNum] = ch
	return ch
}

func (conn *Rtp) read() (rawPkt []byte, channelNum int, err error) {
	if conn.ConnectionType == "TCP" {
		return conn.parseTCPRaw()
	}
	return nil, -1, errors.New("Select ConnectionType")
}

func (conn *Rtp) strtReading() {
	go func() {
		for {
			pkt, channel, err := conn.read()
			if err != nil {
				log.Println(err)
			}
			ch, exist := conn.subscriberList[channel]
			if exist == false {
				continue
			}
			ch <- pkt
		}
	}()
}
