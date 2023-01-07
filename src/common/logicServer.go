package common

import (
	"bufio"
	"encoding/json"
	"github.com/ByteArena/box2d"
	"io"
	"log"
	"net"
	"os"
	"server/src/controller"
	"server/src/game"
	"server/src/game/physic"
	"strconv"
	"strings"
)

const (
	lMenu  int32 = 0
	lWorld       = 1
)

const (
	UUIDPackage                  int32 = 101
	pNewPlayer                         = 102
	pAccessAllowed                     = 103
	pMoveBroadcast                     = 104
	pDisconnect                        = 105
	pEnterTheWorld                     = 106
	pCreateCharacter                   = 107
	pSuccessfulCreationCharacter       = 108
	pGetCharacterList                  = 109
)

type LogicServer struct {
	Addr               string
	Clients            map[int]*Client
	CurrentID          int
	DataBaseController *controller.DataBaseController
	GameController     game.GameController
}

func (s *LogicServer) Init(addr string, maxPlayers int) {
	s.Addr = addr
	s.Clients = make(map[int]*Client, maxPlayers)
	s.DataBaseController = controller.DataBaseConnect(os.Getenv("DATABASE_STRING"))

	gravity := box2d.MakeB2Vec2(0.0, 0.0)
	world := box2d.MakeB2World(gravity)
	collisionSystem := physic.CollisionSystem{}
	collisionSystem.NewListener(&world)

	s.GameController = game.NewGameController(&world, &collisionSystem)

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
		s.Clients[s.CurrentID] = &Client{"", conn, 0, 0, game.Player{}}
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
		access, uid := s.DataBaseController.CheckAccessToken(pack.UUID)
		if access {
			s.Clients[id].AccountID = uid
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

		} else {
			log.Println("Токен не валиден, игрок отключается")
			s.Clients[id].Conn.Close()
		}

	case pEnterTheWorld:
		log.Println("Enter")
		charactedID, _ := strconv.Atoi(pack.Data)
		err := s.NewPlayer(id, charactedID)
		s.Clients[id].Location = lWorld
		if err != nil {
			log.Println("Error:", err)
			return err
		}

	case pMoveBroadcast:
		if s.Clients[id].inWorld() {
			s.MoveBroadCastWithout(pack, id)
		}
	case pCreateCharacter:
		data := strings.Split(pack.Data, "|")
		characterName := data[0]
		appearanceData := data[1]
		if s.DataBaseController.CreateCharacter(characterName, appearanceData, s.Clients[id].AccountID) {
			s.SingleCast(Package{pSuccessfulCreationCharacter, "", ""}, id)
		}
	case pGetCharacterList:
		var charactersData []string
		charactersData = s.DataBaseController.GetCharacters(s.Clients[id].AccountID)
		for _, character := range charactersData {
			characterDataPackage, err := Encode(Package{pGetCharacterList, character, pack.UUID})
			if err != nil {
				log.Println("Error:", err)
				return err
			}
			_, err = s.Clients[id].Conn.Write(characterDataPackage)
		}
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

func (s LogicServer) MoveBroadCastWithout(pack Package, id int) error {
	var transformData physic.TransformData
	json.Unmarshal([]byte(pack.Data), &transformData)

	s.GameController.SetPosition(pack.UUID, transformData)

	data, err := Encode(pack)
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	for _, client := range s.Clients {
		if client.Conn != s.Clients[id].Conn && client.Conn != nil && s.Clients[id].inWorld() {
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
		if client.Conn != s.Clients[id].Conn && client.Conn != nil && s.Clients[id].inWorld() {
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
	s.GameController.RemovePlayerCollider(s.Clients[id].Uuid)
	s.DataBaseController.RemoveAccessToken(s.Clients[id].Uuid)
	delete(s.Clients, id)
}

func (s LogicServer) NewPlayer(id int, characterID int) error {

	s.GameController.AddPlayerCollider(s.Clients[id].Uuid, 0, 0, 0.5)

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
