package game

import (
	"fmt"
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
	"math/rand"
	"server/src/game/physic"
	"time"
)

const (
	screenWidth  = 640
	screenHeight = 480
	scale        = 4
)

type GameController struct {
	world           *box2d.B2World
	collisionSystem *physic.CollisionSystem
	characters      map[string]*Player
}

func NewGameController(world *box2d.B2World, system *physic.CollisionSystem) GameController {
	g := GameController{}
	g.world = world
	g.collisionSystem = system
	g.characters = make(map[string]*Player)
	return g
}

/*
func main() {

	gravity := box2d.MakeB2Vec2(0.0, -10.0)
	world := box2d.MakeB2World(gravity)
	sys := physic.CollisionSystem{}
	sys.NewListener(&world)

	game := newGameController(&world, &sys)

	{
		bd := box2d.MakeB2BodyDef()
		ground := world.CreateBody(&bd)

		shape := box2d.MakeB2EdgeShape()
		shape.Set(box2d.MakeB2Vec2(-20.0, -20), box2d.MakeB2Vec2(20.0, -25))
		ground.CreateFixture(&shape, 0.0)
		game.characters["plane"] = ground
		sys.Add(&physic.Box2dComponent{ground})
	}
	game.AddCircle("circle_1", 3, 5, 0.5)
	game.AddCircle("circle_2", 8, 5, 1)
	game.Run()

}
*/
func (g *GameController) Update() error {

	timeStep := 1.0 / 60.0
	velocityIterations := 8
	positionIterations := 3

	g.world.Step(timeStep, velocityIterations, positionIterations)
	g.collisionSystem.Update(float32(timeStep))
	// Now print the position and angle of the body.
	for name, character := range g.characters {
		position := character.Collider.GetPosition()
		angle := character.Collider.GetAngle()
		msg := fmt.Sprintf("(%s): %4.3f %4.3f %4.3f\n", name, position.X, position.Y, angle)
		fmt.Sprintf(msg)
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		rand.Seed(time.Now().UnixNano())
		//id := rand.Intn(100)
		//g.AddPlayerCollider("circle_"+string(id), 3, 3, 1)
	}

	return nil
}

func (g *GameController) Draw(screen *ebiten.Image) {
	/*
		vector.StrokeLine(screen, 100, 100, 300, 100, 1, color.RGBA{0xff, 0xff, 0xff, 0xff})
		vector.StrokeLine(screen, 50, 150, 50, 350, 1, color.RGBA{0xff, 0xff, 0x00, 0xff})
		vector.StrokeLine(screen, 50, 100+cf, 200+cf, 250, 4, color.RGBA{0x00, 0xff, 0xff, 0xff})

		vector.DrawFilledRect(screen, 50+cf, 50+cf, 100+cf, 100+cf, color.RGBA{0x80, 0x80, 0x80, 0xc0})
		vector.StrokeRect(screen, 300-cf, 50, 120, 120, 10+cf/4, color.RGBA{0x00, 0x80, 0x00, 0xff})

		vector.DrawFilledCircle(screen, 400, 400, 100, color.RGBA{0x80, 0x00, 0x80, 0x80})
	*/

	for _, character := range g.characters {
		position := character.Collider.GetPosition()
		vector.StrokeCircle(screen, float32(position.X*scale), -float32(position.Y*scale), float32(character.Collider.M_fixtureList.M_shape.GetRadius()*scale), 1, color.RGBA{0xff, 0x80, 0xff, 0xff})
	}
	/*
		for _, e := range g.collisionSystem.Entities {
			fmt.Println(e.Body.GetPosition().X)
			vector.StrokeCircle(screen, float32(e.Body.GetPosition().X*scale), float32(e.Body.GetPosition().X*scale), float32(e.Body.M_fixtureList.M_shape.GetRadius()*scale), 1, color.RGBA{0xff, 0x80, 0xff, 0xff})
		}
	*/
	//vector.StrokeCircle(screen, 400, 400, 10, 1, color.RGBA{0xff, 0x80, 0xff, 0xff})

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
}

func (g *GameController) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *GameController) Run() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Server Physic Window")
	if err := ebiten.RunGame(g); err != nil {
		log.Println(err)
	}
}

func (g *GameController) AddPlayerCollider(name string, x float64, y float64, r float64) {

	bd := box2d.MakeB2BodyDef()
	bd.Position.Set(x, y)
	bd.Type = box2d.B2BodyType.B2_dynamicBody
	bd.FixedRotation = false
	bd.AllowSleep = false

	body := g.world.CreateBody(&bd)

	shape := box2d.MakeB2CircleShape()
	shape.M_radius = r

	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &shape
	fd.Density = 20.0
	body.CreateFixtureFromDef(&fd)

	player := Player{Entity{Name: name, Collider: *body, Stats: EntityStat{}}}

	g.characters[name] = &player

	g.collisionSystem.Add(&physic.Box2dComponent{&player.Collider})
}
