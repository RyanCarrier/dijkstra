package dijkstra

import (
	"slices"
	"testing"
)

func TestLists(t *testing.T) {
	inputs := []currentDistance{
		{0, 1},
		{1, 2},
		{2, 3},
		{3, 4},
		{4, 5},
		{5, 6},
		{6, 7},
		{7, 8},
		{8, 9},
		{9, 10},
		{10, 1 * 2},
		{11, 2 * 2},
		{12, 3 * 2},
		{13, 4 * 2},
		{14, 5 * 2},
		{15, 6 * 2},
		{16, 7 * 2},
		{17, 8 * 2},
		{18, 9 * 2},
		{19, 10 * 2},
		{99, 10},
	}
	shortResult := []currentDistance{
		{0, 1},
		{1, 2},
		{10, 1 * 2},
		{2, 3},
		{3, 4},
		{11, 2 * 2},
		{4, 5},
		{5, 6},
		{12, 3 * 2},
		{6, 7},
		{7, 8},
		{13, 4 * 2},
		{8, 9},
		{9, 10},
		{99, 10},
		{14, 5 * 2},
		{15, 6 * 2},
		{16, 7 * 2},
		{17, 8 * 2},
		{18, 9 * 2},
		{19, 10 * 2},
	}
	longResult := make([]currentDistance, len(shortResult))
	copy(longResult, shortResult)
	slices.Reverse(longResult)
	lists := []struct {
		listEnum int
		listName string
	}{
		{listShortPQ, "listShortPQ"},
		{listLongPQ, "listLongPQ"},
		{listShortLL, "listShortLL"},
		{listLongLL, "listLongLL"},
	}
	for _, list := range lists {
		t.Run(list.listName, func(t *testing.T) {
			dl := (&Graph{}).getList(list.listEnum)
			var result []currentDistance
			if list.listEnum == listShortPQ || list.listEnum == listShortLL {
				result = shortResult
			} else {
				result = longResult
			}
			for _, input := range inputs {
				dl.PushOrdered(input)
			}
			for i, want := range result {
				if got := dl.PopOrdered(); got.distance != want.distance {
					t.Errorf("Incorrect order, index %d\nwant:%v\n got:%v", i, want, got)
				}
			}

		})
	}

}
