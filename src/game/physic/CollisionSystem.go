package physic

import (
	"github.com/ByteArena/box2d"
)

type Box2dComponent struct {
	// Body is the box2d body
	Body *box2d.B2Body
}

type CollisionEntity struct {
	*Box2dComponent
}

type CollisionSystem struct {
	Entities []CollisionEntity
}

func (c *CollisionSystem) Add(box *Box2dComponent) {
	c.Entities = append(c.Entities, CollisionEntity{box})
}

func (c *CollisionSystem) NewListener(w *box2d.B2World) {
	w.SetContactListener(c)
}

func (c *CollisionSystem) Update(dt float32) {}

// BeginContact implements the B2ContactListener interface.
// when a BeginContact callback is made by box2d, it sends a message containing
// the information from the callback.
func (c *CollisionSystem) BeginContact(contact box2d.B2ContactInterface) {
	//fmt.Println(fmt.Sprintf("Begin contact is toching: %v", contact.IsTouching()))
}

// EndContact implements the B2ContactListener interface.
// when a EndContact callback is made by box2d, it sends a message containing
// the information from the callback.
func (c *CollisionSystem) EndContact(contact box2d.B2ContactInterface) {
	//fmt.Println(fmt.Sprintf("End contact is toching: %v", contact.IsTouching()))
}

// PreSolve implements the B2ContactListener interface.
// this is called after a contact is updated but before it goes to the solver.
// When it is called, a message is sent containing the information from the callback
func (c *CollisionSystem) PreSolve(contact box2d.B2ContactInterface, oldManifold box2d.B2Manifold) {
	//fmt.Println(fmt.Sprintf("PreSolve contact is toching: %v", contact.IsTouching()))
}

// PostSolve implements the B2ContactListener interface.
// this is called after the solver is finished.
// When it is called, a message is sent containing the information from the callback
func (c *CollisionSystem) PostSolve(contact box2d.B2ContactInterface, impulse *box2d.B2ContactImpulse) {
	//fmt.Println(fmt.Sprintf("PostSolve contact is toching: %v", contact.IsTouching()))
}
