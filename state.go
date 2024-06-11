package main

// Type for representing current state of user
type State int

// State constatnts
const (
	Welcome State = iota
	HasImage
	HasAudio
	HasBoth
	Generating
)
