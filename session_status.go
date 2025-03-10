package delugeclient

import (
	"fmt"

	"github.com/gdm85/go-rencode"
)

// SessionStatus contains basic session status and statistics.
type SessionStatus struct {
	HasIncomingConnections bool
	UploadRate             float32
	DownloadRate           float32
	PayloadUploadRate      float32
	PayloadDownloadRate    float32
	TotalDownload          int64
	TotalUpload            int64
	NumPeers               int16
	DhtNodes               int16
}

// sessionStatusKeys is a slice with specific session status and statistics.
var sessionStatusKeys = rencode.NewList(
	"has_incoming_connections",
	"upload_rate",
	"download_rate",
	"payload_upload_rate",
	"payload_download_rate",
	"total_download",
	"total_upload",
	"num_peers",
	"dht_nodes",
)

// GetSessionStatus retrieves session status and statistics.
func (c *Client) GetSessionStatus() (*SessionStatus, error) {
	var args rencode.List
	args.Add(sessionStatusKeys)

	rd, err := c.rpcWithDictionaryResult("core.get_session_status", args, rencode.Dictionary{})
	if err != nil {
		return nil, err
	}

	if c.settings.Logger != nil {
		c.settings.Logger.Printf("session status: %+s", rd)
	}
	var data SessionStatus
	err = rd.ToStruct(&data, c.excludeV2tag)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// GetAASessionStatus retrieves every session status and statistics from libtorrent.
func (c *Client) GetAllSessionStatus() (map[string]int, error) {
	var args rencode.List
	args.Add("")

	rd, err := c.rpcWithDictionaryResult("core.get_session_status", args, rencode.Dictionary{})
	if err != nil {
		return nil, err
	}

	dic := make(map[string]int)
	keys := rd.Keys()
	values := rd.Values()
	for i := 0; i < len(keys); i++ {
		ks := fmt.Sprintf("%s", keys[i])
		var vs int
		switch values[i].(type) {
		case int:
			vs = values[i].(int)
		case int64:
			vs = int(values[i].(int64))
		case int32:
			vs = int(values[i].(int32))
		case int16:
			vs = int(values[i].(int16))
		case int8:
			vs = int(values[i].(int8))
		case float64:
			vs = int(values[i].(float64))
		case float32:
			vs = int(values[i].(float32))
		default:
			//fmt.Println("unknown type")
		}

		dic[ks] = vs
	}

	if c.settings.Logger != nil {
		c.settings.Logger.Printf("session status: %+s", rd)
	}

	return dic, nil
}
