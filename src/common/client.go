package common

import (
	"github.com/google/uuid"
	"net"
)

type Client struct {
	Uuid uuid.UUID
	Conn net.Conn
}
