package rtp

import (
	"behnama/pipe/conn"
	"behnama/pipe/rtsp"
	"bufio"
	"os"
	"testing"
	"time"

	pionRtp "github.com/pion/rtp"

	"log"

	"github.com/stretchr/testify/suite"
)

type rtpSuite struct {
	suite.Suite
	rtsp *rtsp.Rtsp
	rtp  *Rtp
}

func (suite *rtpSuite) SetupTest() {
	conn, url, err := conn.NewRtspTCPConn("rtsp://192.168.14.23:554/rtsp_tunnel?h26x=4&line=1&inst=1vcd=2",
		time.Duration(15*time.Second),
		time.Duration(5*time.Second),
		time.Duration(5*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	readWriter := bufio.NewReadWriter(reader, writer)
	rtsp, err := rtsp.NewRtsp(readWriter, *url)
	if err != nil {
		log.Fatal(err)
	}
	rtsp.Play()
	rtp := NewRTPConnection(reader, "TCP")
	suite.rtp = rtp
	suite.rtsp = rtsp
}

// func saveTestVideo() {
// 	f, _ := os.OpenFile("test.mp4", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	f.Write((packet.Payload))
// }

func (suite *rtpSuite) TestReadRTP() {
	packetChan := suite.rtp.SubscribeChannel(0)
	counter := 0
	f, _ := os.Open("test.mp4")
	for rawPacket := range packetChan {
		if counter == 200 {
			break
		}
		packet := pionRtp.Packet{}
		err := packet.Unmarshal(rawPacket)
		if err != nil {
			suite.FailNow(err.Error())
		}
		buffer := make([]byte, len(packet.Payload))
		f.Read(buffer)
		suite.Equal(packet.Payload, buffer)
		counter++
	}
}

func TestRtpSuite(t *testing.T) {
	suite.Run(t, new(rtpSuite))
}
