package dijkstra

import (
	"errors"
	"fmt"
)

var ErrWrongFormat = errors.New("wrong source format")
var ErrNoPath = errors.New("no path found")
var ErrLoopDetected = errors.New("infinite loop detected")
var ErrVertexNotFound = errors.New("vertex not found")
var ErrVertexAlreadyExists = errors.New("vertex already exists")
var ErrVertexNegative = errors.New("vertex is negative")
var ErrArcNotFound = errors.New("arc not found")
var ErrGraphNotValid = errors.New("graph is not valid")
var ErrMapNotFound = errors.New("mapping error, can not find mapped vertex")
var ErrArcHanging = errors.New("arc will be left hanging")

// not found/item validity
func newErrMapNotFound(a int) error {
	return fmt.Errorf("%w '%d'", ErrMapNotFound, a)
}
func newErrVertexNotFound(a int) error {
	return fmt.Errorf("%d %w", a, ErrVertexNotFound)
}
func newErrVertexNegative(a int) error {
	return fmt.Errorf("%d %w", a, ErrVertexNegative)
}
func newErrVertexAlreadyExists(a int) error {
	return fmt.Errorf("%d %w", a, ErrVertexAlreadyExists)
}
func newErrArcNotFound[T comparable](a, b T) error {
	return fmt.Errorf("%v->%v %w", a, b, ErrArcNotFound)
}
func newArcHanging[T comparable](a, b T) error {
	return fmt.Errorf("removing %v; %v->%v %w", b, a, b, ErrArcHanging)
}

// graph issues
func newErrLoop(a, b int) error {
	return fmt.Errorf("%w, from node '%d' to node '%d'", ErrLoopDetected, a, b)
}
func newErrGraphNotValid(a, b int) error {
	return fmt.Errorf("%w, arc %d->%d, %d not found", ErrGraphNotValid, a, b, b)
}
func newErrNoPath(a, b int) error {
	return fmt.Errorf("%d->%d %w", a, b, ErrNoPath)
}

// mappped
func newErrMappedVertexNotFound[T comparable](a T) error {
	return fmt.Errorf("%v %w", a, ErrVertexNotFound)
}
func newErrMappedVertexAlreadyExists[T comparable](a T) error {
	return fmt.Errorf("%v %w", a, ErrVertexAlreadyExists)
}
func newErrMappedArcNotFound[T comparable](a, b T) error {
	return fmt.Errorf("%v->%v %w", a, b, ErrArcNotFound)
}
