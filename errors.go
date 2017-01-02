package dijkstra

import "errors"

//ErrWrongFormat is thrown when the source file input is in an incorrect format
var ErrWrongFormat = errors.New("Wrong source format")

//ErrNoPath is thrown when there is no path from src to dest
var ErrNoPath = errors.New("Wrong source format")

//ErrMixMapping is thrown when there is a mixture of integers and strings in the input file
var ErrMixMapping = errors.New("Potential mixing of integer and string node ID's :" + ErrWrongFormat.Error())

//ErrLoopDetected is thrown when a loop is detected, causing the distance to go to infinity
var ErrLoopDetected = errors.New("Infinite loop detected")
