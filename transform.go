package pipescript

import "fmt"

// TransformInstance is the interface which underlies each transform. Copy() should return an exact copy
// of this transform instance as it is right now.
type TransformInstance interface {
	Next(*TransformIterator) (*Datapoint, error)
	Copy() TransformInstance
}

// TransformGenerator creates a new TransformInstance from the generator. The args passed in are the scripts
// that pass in. The generator returns the script arguments that are to be used during runtime, the TransformInstance
// and an error. The Script array is used to allow argument magic in transforms, such as using constants on init
// and removing the args during computation. It is also used to implement more advanced things, such as the select statement,
// which actually uses the Script as its arg, and not the transformed arguments. constant notifies if the transform will return
// a constant value. If this is True, then the constant can be extracted during optimization by calling the transform with a dummy
// datapoint.
type TransformGenerator func(name string, args []*Script) (script []*Script, ti TransformInstance, constant bool, err error)

// TransformArg represents an argument passed into the transform function
type TransformArg struct {
	Description string      `json:"description"`       // A description of what the arg represents
	Optional    bool        `json:"optional"`          // Whether the arg is optional
	Default     interface{} `json:"default,omitempty"` // If the arg is optional, what is its default value
	Constant    bool        `json:"constant"`          // If the argument must be a constant (ie, not part of a transform)
}

// Transform is the struct which holds the name, docstring, and generator for a transform function
type Transform struct {
	Name         string         `json:"name"`              // The name of the transform
	Description  string         `json:"description"`       // The description of the transform - a docstring
	InputSchema  string         `json:"ischema,omitempty"` // The schema of the input datapoint that the given transform expects (optional)
	OutputSchema string         `json:"oschema,omitempty"` // The schema of the output data that this transform gives (optional).
	Args         []TransformArg `json:"args"`              // The arguments that the transform accepts
	OneToOne     bool           `json:"one_to_one"`        //Whether or not the transform is one to one

	Generator TransformGenerator `json:"-"` // The generator function of the transform
}

var (
	// TransformRegistry is the map of all the transforms that are currently registered.
	// Do not manually add/remove elements from this map.
	// Use Transform.Register to insert new transforms.
	TransformRegistry = make(map[string]Transform)
)

// Register registers the transform with the system. Note that it is not threadsafe. Register is
// assumed to be run once at the startup of the query system. Adding functions during runtime is
// not supported.
func (t Transform) Register() error {
	if t.Name == "" || t.Generator == nil {
		err := fmt.Errorf("Attempted to register invalid transform: '%s'", t.Name)
		return err
	}
	_, ok := TransformRegistry[t.Name]
	if ok {
		err := fmt.Errorf("A transform with the name '%s' already exists.", t.Name)
		return err
	}

	TransformRegistry[t.Name] = t

	return nil
}

// Script generates a Script from the given Transform, when given the arguments passed into the Transform
// as scripts. It validates the arguments and information given transform and argument metadata
func (t *Transform) Script(args []*Script) (*Script, error) {
	// First, check the argument count
	if len(args) > len(t.Args) {
		return nil, fmt.Errorf("%d arguments were passed to '%s', which only accepts %d", len(args), t.Name, len(t.Args))
	}

	// Now check the ordering/Constantness of arguments
	for i := range t.Args {
		// Check if the argument is given
		if len(args) >= i {
			// The argument is given
			if !args[i].IsOneToOne {
				return nil, fmt.Errorf("Argument %d of transform '%s' must be OneToOne", i+1, t.Name)
			}
			if t.Args[i].Constant && !args[i].Constant {
				return nil, fmt.Errorf("Argument %d of transform '%s' must be a constant.", i+1, t.Name)
			}
		} else {
			// The argument was not given
			if !t.Args[i].Optional {
				return nil, fmt.Errorf("Argument %d of transform '%s' is required.", i+1, t.Name)
			}

			// The argument was optional and was NOT passed in. We set it up using a ConstantScript using the default value
			// given in the transform config. Note that we assume that all previous arguments MUST exist. This means that there can't
			// be optional arguments in between required arguments.
			args = append(args, ConstantScript(t.Args[i].Default))
		}
	}

	// All of the arguments were checked. Send the args to the Generator
	args, ti, c, err := t.Generator(t.Name, args)
	if err != nil {
		return nil, err
	}

	// Now we have everything necessary to generate a Script out of the Transform
	pe, err := NewPipelineElement(args, ti)

	return &Script{
		input:      pe,
		output:     pe,
		IsOneToOne: t.OneToOne,
		Constant:   c,
	}, err
}