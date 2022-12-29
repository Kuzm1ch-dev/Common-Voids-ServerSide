package common

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"server/src/controller"
	"server/src/game"
)

const (
	UUIDPackage    int32 = 101
	pNewPlayer           = 102
	pAccessAllowed       = 103
	pBroadcast           = 104
	pDisconnect          = 105
)

type LogicServer struct {
	Addr               string
	Clients            map[int]*Client
	CurrentID          int
	DataBaseController *controller.DataBaseController
}

func (s *LogicServer) Init(addr string, maxPlayers int) {
	s.Addr = addr
	s.Clients = make(map[int]*Client, maxPlayers)
	s.DataBaseController = controller.DataBaseConnect(os.Getenv("DATABASE_STRING"))
}

func (s *LogicServer) ListenAndServe() error {
	Listener, err := net.Listen("tcp", s.Addr+":30000")
	if err != nil {
		log.Println("Error:", err)
		return err
	}

	log.Println("Logic server started ", s.Addr)

	for {
		conn, err := Listener.Accept()
		if err != nil {
			log.Println("accept failed, err:", err)
			continue
		}
		s.Clients[s.CurrentID] = &Client{"", conn, game.Player{}}
		go s.ConnectionHandler(s.CurrentID)
		s.CurrentID++
	}
}

func (s LogicServer) ConnectionHandler(id int) error {

	defer s.CloseConnection(id)

	//Вроверяем валидность Access Token'а
	checkAccessToken, err := Encode(Package{UUIDPackage, "", ""})
	if err != nil {
		log.Println("Encode Package Error:", err)
		return err
	}

	log.Println("Check Access Token ", id)
	_, err = s.Clients[id].Conn.Write(checkAccessToken)
	if err != nil {
		log.Println("Write Package Error:", err)
		return err
	}

	for {
		reader := bufio.NewReader(s.Clients[id].Conn)
		pack, err := Decode(reader)
		if err == io.EOF {
			log.Println("Error:", err)
			return err
		}
		if err != nil {
			log.Println("decode msg failed, err:", err)
			return err
		}

		err = s.handleReceivedPacket(pack, id)
		if err != nil {
			log.Println("Error:", err)
			return err
		}
	}
}

func (s LogicServer) handleReceivedPacket(pack Package, id int) error {
	switch pack.Code {
	case UUIDPackage:
		if s.DataBaseController.CheckAccessToken(pack.UUID) {
			accessAllowed, err := Encode(Package{pAccessAllowed, "", pack.UUID})
			if err != nil {
				log.Println("Error:", err)
				return err
			}

			s.Clients[id].SetUUID(pack.UUID)
			log.Println("Токен валиден, игрок авторизируется")
			_, err = s.Clients[id].Conn.Write(accessAllowed)

			if err != nil {
				log.Println("Error:", err)
				return err
			}

			err = s.NewPlayer(id)
			if err != nil {
				log.Println("Error:", err)
				return err
			}
		} else {
			log.Println("Токен не валиден, игрок отключается")
			s.Clients[id].Conn.Close()
		}
	case pBroadcast:
		s.BroadCastWithout(pack, id)
	default:
		break
	}
	return nil
}

func (s LogicServer) SingleCast(pack Package, id int) error {
	data, err := Encode(pack)
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	_, err = s.Clients[id].Conn.Write(data)
	if err != nil {
		log.Println("Error:", err.Error())
		return err
	}
	return nil
}

func (s LogicServer) BroadCast(pack Package) error {
	data, err := Encode(pack)
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	for _, client := range s.Clients {
		if client.Conn != nil {
			_, err := client.Conn.Write(data)
			if err != nil {
				log.Println("Error:", err.Error())
				return err
			}
		}
	}
	return nil
}

func (s LogicServer) BroadCastWithout(pack Package, id int) error {
	data, err := Encode(pack)
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	for _, client := range s.Clients {
		if client.Conn != s.Clients[id].Conn && client.Conn != nil {
			_, err := client.Conn.Write(data)
			if err != nil {
				log.Println("Error:", err.Error())
				return err
			}
		}
	}
	return nil
}

func (s LogicServer) CloseConnection(id int) {
	err := s.Clients[id].Conn.Close()
	if err != nil {
		log.Println("Error:", err)
	}
	s.BroadCastWithout(Package{pDisconnect, "", s.Clients[id].Uuid}, id)
	log.Println(s.Clients[id].Uuid + " Отключился.")
	delete(s.Clients, id)
}

func (s LogicServer) NewPlayer(id int) error {
	for _, v := range s.Clients {
		if v != s.Clients[id] {
			data, err := Encode(Package{pNewPlayer, "", s.Clients[id].Uuid})
			log.Println("Отсылаем ", v.Uuid, " игроку пакет с новым игроком: ", s.Clients[id].Uuid)
			if err != nil {
				log.Println("Error:", err.Error())
			}

			//Говорим всем, кто на сервере, что мы зашли
			_, err = v.Conn.Write(data)
			if err != nil {
				log.Println("Error:", err.Error())
			}

			//Подгружаем всех игроков, которые уже на сервере
			aboutPlayer, err := Encode(Package{pNewPlayer, "", v.Uuid})
			log.Println("Отсылаем ", s.Clients[id].Uuid, " игроку пакет с другим игроком: ", v.Uuid)
			_, err = s.Clients[id].Conn.Write(aboutPlayer)
			if err != nil {
				log.Println("Error:", err.Error())
			}
		}
	}
	return nil
}
