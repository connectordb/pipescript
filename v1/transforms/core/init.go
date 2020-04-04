/*
Package core contains the core basic transforms which are available in PipeScript.
It should be imported by default by basically all users of PipeScript
*/
package core

import "github.com/heedy/pipescript"

// These scripts can be used internally
var iflast *pipescript.Script
var identity *pipescript.Script

func init() {
	// Manually generate the if last script
	last, err := Last.Script(nil)
	if err != nil {
		panic("Could not generate internal script for 'if last':" + err.Error())
	}
	iflast, err = If.Script([]*pipescript.Script{last})
	if err != nil {
		panic("Could not generate internal if script for 'if last':" + err.Error())
	}

	identity, err = IdentityTransform.Script([]*pipescript.Script{pipescript.ConstantScript(nil)})
	if err != nil {
		panic("Could not generate identity script: " + err.Error())
	}

	// Set up the default args for the transforms which use them
	Remember.Args[1].Default = identity
	Ifelse.Args[2].Default = identity
	Top.Args[1].Default = identity
	Bottom.Args[1].Default = identity
}

func Register() {
	Set.Register()
	Del.Register()

	If.Register()

	IdentityTransform.Register()
	D.Register()

	Ifelse.Register()

	Changed.Register()

	First.Register()
	Last.Register()

	//prev.Register()
	Next.Register()

	Count.Register()
	I.Register() // Same as count except 1 less
	Sum.Register()

	Length.Register()

	T.Register()
	Tshift.Register()
	Ttrue.Register()
	Tfalse.Register()
	Dt.Register()

	IMap.Register()
	Map.Register()
	Reduce.Register()
	Filter.Register()
	Apply.Register()

	Remember.Register()

	IWhile.Register()
	While.Register()

	Rand.Register()

	New.Register()

	AllTrue.Register()
	AnyTrue.Register()

	Bucket.Register()

	Top.Register()
	Bottom.Register()
}
