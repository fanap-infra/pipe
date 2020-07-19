package rtsp

import (
	"errors"
	"fmt"
	"strings"
)

//Play RTSP command.
func (conn *Rtsp) Play() error {
	if err := conn.prepareStage(stageSetupDone); err != nil {
		return err
	}
	req := request{
		method: "PLAY",
		uri:    trimURLUser(conn.url),
	}
	req.header = append(req.header, "Session: "+conn.sessionKey)
	if err := conn.writeRequest(req); err != nil {
		return err
	}
	headers, err := conn.readRtspResponse()
	if err != nil {
		return err
	}
	if string(headers[0]) != "RTSP/1.0 200 OK" {
		return errors.New("Describe failed")
	}
	conn.stage = stagePlayDone
	return nil
}

//Describe RTSP Describe command
func (conn *Rtsp) Describe() error {
	req := request{
		method: "DESCRIBE",
		uri:    trimURLUser(conn.url),
		header: []string{"Accept: application/sdp"},
	}
	if err := conn.writeRequest(req); err != nil {
		return err
	}
	headers, err := conn.readRtspResponse()
	if err != nil {
		return err
	}
	if string(headers[0]) != "RTSP/1.0 200 OK" {
		return errors.New("Describe failed")
	}
	if len(conn.SDP.MediaDescriptions) == 0 {
		return errors.New("Describe failed, Cannot read SDP")
	}
	conn.stage = stageDescribeDone
	return nil
}

//Setup RTSP command
func (conn *Rtsp) Setup() error {
	if err := conn.prepareStage(stageDescribeDone); err != nil {
		return err
	}
	for index, mediaDesc := range conn.SDP.MediaDescriptions {
		uri := ""
		control, _ := mediaDesc.Attribute("control")
		if strings.HasPrefix(control, "rtsp://") {
			uri = control
		} else {
			uri = trimURLUser(conn.url) + "/" + control
		}
		req := request{method: "SETUP", uri: uri}
		req.header = append(req.header,
			fmt.Sprintf("Transport: RTP/AVP/TCP;unicast;interleaved=%d-%d",
				index*2, index*2+1))
		if conn.sessionKey != "" {
			req.header = append(req.header, "Session: "+conn.sessionKey)
		}
		if err := conn.writeRequest(req); err != nil {
			return err
		}
		if _, err := conn.readRtspResponse(); err != nil {
			return err
		}
	}
	conn.stage = stageSetupDone
	return nil
}

func (conn *Rtsp) prepareStage(stage int) error {
	for conn.stage < stage {
		switch conn.stage {
		case stageIdle:
			if err := conn.Describe(); err != nil {
				return err
			}
		case stageDescribeDone:
			if err := conn.Setup(); err != nil {
				return err
			}
		case stageSetupDone:
			if err := conn.Play(); err != nil {
				return err
			}
		}
	}
	return nil
}
