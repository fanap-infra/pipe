package rtsp

import (
	"behnama/pipe/conn"
	"bufio"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type RtspSuite struct {
	suite.Suite
	conn *Rtsp
}

func (suite *RtspSuite) SetupTest() {
	conn, url, err := conn.NewRtspTCPConn("rtsp://127.0.0.1:8554/test",
		time.Duration(15*time.Second),
		time.Duration(5*time.Second),
		time.Duration(5*time.Second),
	)
	if err != nil {
		suite.FailNow("RTSP Conn failed!")
	}
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	readWriter := bufio.NewReadWriter(reader, writer)
	suite.conn, _ = NewRtsp(readWriter, *url)
}

func (suite *RtspSuite) TestDescribe() {
	suite.conn.Describe()
	suite.Equal(suite.conn.stage, stageDescribeDone)
}

func (suite *RtspSuite) TestSetup() {
	suite.conn.Setup()
	suite.Equal(suite.conn.stage, stageSetupDone)
}

func (suite *RtspSuite) TestPlay() {
	suite.conn.Play()
	suite.Equal(suite.conn.stage, stagePlayDone)
}

func TestRtspSuite(t *testing.T) {
	suite.Run(t, new(RtspSuite))
}
