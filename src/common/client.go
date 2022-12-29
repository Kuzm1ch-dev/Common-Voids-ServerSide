package common

import (
	"net"
	"server/src/game"
)

type Client struct {
	Uuid   string
	Conn   net.Conn
	Player game.Player
}

func (client *Client) SetUUID(uuid string) {
	client.Uuid = uuid
}
