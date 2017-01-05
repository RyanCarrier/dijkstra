package dijkstra

import (
	"errors"
	"fmt"
)

//ErrWrongFormat is thrown when the source file input is in an incorrect format
var ErrWrongFormat = errors.New("Wrong source format")

//ErrNoPath is thrown when there is no path from src to dest
var ErrNoPath = errors.New("No path found")

//ErrMixMapping is thrown when there is a mixture of integers and strings in the input file
var ErrMixMapping = errors.New("Potential mixing of integer and string node ID's :" + ErrWrongFormat.Error())

//ErrLoopDetected is thrown when a loop is detected, causing the distance to go
// to inf (or -inf), or just generally loop forever
var ErrLoopDetected = errors.New("Infinite loop detected")

//NewErrLoop generates a new error with details for loop error
func newErrLoop(a, b int) error {
	return errors.New(fmt.Sprint(ErrLoopDetected.Error(), "From node '", a, "' to node '", b, "'"))
}
