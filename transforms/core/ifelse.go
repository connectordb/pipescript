package core

import (
	"errors"

	"github.com/connectordb/pipescript"
)

type ifelseTransformStruct struct {
	ifScript   *pipescript.Script
	elseScript *pipescript.Script
	iter       *pipescript.SingleDatapointIterator
}

func (t ifelseTransformStruct) Copy() (pipescript.TransformInstance, error) {
	newif, err := t.ifScript.Copy()
	if err != nil {
		return nil, err
	}
	newelse, err := t.elseScript.Copy()
	if err != nil {
		return nil, err
	}
	iter := &pipescript.SingleDatapointIterator{}
	newif.SetInput(iter)
	newelse.SetInput(iter)
	return ifelseTransformStruct{newif, newelse, iter}, nil
}

func (t ifelseTransformStruct) Next(ti *pipescript.TransformIterator) (dp *pipescript.Datapoint, err error) {
	te := ti.Next()
	if te.IsFinished() {
		// Clear the internal scripts
		t.iter.Set(nil, nil)
		t.ifScript.Next()
		t.iter.Set(nil, nil)
		t.elseScript.Next()
		return te.Get()
	}

	v, err := te.Args[0].Bool()
	if err != nil {
		return nil, err
	}
	t.iter.Set(te.Datapoint, nil)
	if v {
		return t.ifScript.Next()
	}
	return t.elseScript.Next()

}

var Ifelse = pipescript.Transform{
	Name:        "ifelse",
	Description: "A conditional. This is what an if statement would be in other languages.",
	Args: []pipescript.TransformArg{
		{
			Description: "The statement to check for truth",
		},
		{
			Description: "pipe to run if conditional is true",
			Hijacked:    true,
		},
		{
			Description: "Pipe to run if conditional is false",
			Hijacked:    true,
			Optional:    true,
		},
	},
	OneToOne: true,
	Generator: func(name string, args []*pipescript.Script) (*pipescript.TransformInitializer, error) {
		if args[1].Peek || args[2].Peek {
			return nil, errors.New("ifelse cannot be used with transforms that peek.")
		}
		if args[2].Constant {
			dp, err := args[2].GetConstant()
			if err != nil {
				return nil, err
			}
			if dp.Data == nil {
				args[2], err = IdentityTransform.Script([]*pipescript.Script{args[2]}) // Get an identity if no else
			}
		}

		iter := &pipescript.SingleDatapointIterator{}
		args[1].SetInput(iter)
		args[2].SetInput(iter)
		return &pipescript.TransformInitializer{Args: args[0:1], Transform: ifelseTransformStruct{args[1], args[2], iter}}, nil
	},
}
