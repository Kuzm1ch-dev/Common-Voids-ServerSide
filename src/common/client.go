package common

import (
	"net"
	"server/src/game"
)

type Client struct {
	Uuid      string
	Conn      net.Conn
	AccountID int32
	Location  int32
	Player    game.Player
}

func (client *Client) inWorld() bool {
	if client.Location == 1 {
		return true
	} else {
		return false
	}
}

func (client *Client) SetUUID(uuid string) {
	client.Uuid = uuid
}
