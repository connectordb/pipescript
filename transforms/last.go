package transforms

import "github.com/connectordb/pipescript"

type lastTransform struct{}

// Copy creates a copy of the last transform
func (lt *lastTransform) Copy() pipescript.TransformInstance {
	return &lastTransform{}
}

// Next returns the next element of the transform
func (lt *lastTransform) Next(ti *pipescript.TransformIterator) (*pipescript.Datapoint, error) {
	te := ti.Next()
	if te.IsFinished() {
		return te.Get()
	}
	// Peek at the next datapoint, to find out if it is nil (ie, the current datapoint is the last one)
	te2 := ti.Peek(0)

	return te.Set(te2.IsFinished())
}

var last = pipescript.Transform{
	Name:         "last",
	Description:  "Returns true if last datapoint of a sequence, and false otherwise",
	OutputSchema: `{"type": "boolean"}`,
	OneToOne:     true,
	Generator: func(name string, args []*pipescript.Script) ([]*pipescript.Script, pipescript.TransformInstance, bool, error) {
		return nil, &lastTransform{}, false, nil
	},
}