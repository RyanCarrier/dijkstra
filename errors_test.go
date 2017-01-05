package dijkstra

import (
	"fmt"
	"testing"
)

func TestErrLoop(t *testing.T) {
	if newErrLoop(0, 1).Error() != fmt.Sprint(ErrLoopDetected.Error(), "From node '", 0, "' to node '", 1, "'") {
		t.Error("ErrLoop doesn't match")
	}
}
