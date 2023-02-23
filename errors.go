package dijkstra

import (
	"errors"
	"fmt"
)

//ErrWrongFormat is thrown when the source file input is in an incorrect format
var ErrWrongFormat = errors.New("wrong source format")

//ErrNoPath is thrown when there is no path from src to dest
var ErrNoPath = errors.New("no path found")

//ErrMixMapping is thrown when there is a mixture of integers and strings in the input file
var ErrMixMapping = errors.New("potential mixing of integer and string node ID's :" + ErrWrongFormat.Error())

//ErrLoopDetected is thrown when a loop is detected, causing the distance to go
// to inf (or -inf), or just generally loop forever
var ErrLoopDetected = errors.New("infinite loop detected")

//ErrNodeNotFound is thrown when a node is not found in the graph when it is being requested/used
var ErrNodeNotFound = errors.New("node not found")

//ErrNoMap is thrown when the map is not being used but is being requested/accessed
var ErrNoMap = errors.New("map is not being used/initialised")

//ErrAlreadyCalculating is thrown when the algorithm is already running
var ErrAlreadyCalculating = errors.New("already calculating")

//NewErrLoop generates a new error with details for loop error
func newErrLoop(a, b int) error {
	return errors.New(fmt.Sprint(ErrLoopDetected.Error(), " from node '", a, "' to node '", b, "'"))
}
