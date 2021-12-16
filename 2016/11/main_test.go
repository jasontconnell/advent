package main

import (
	"testing"
)

func TestGetAllMoves(t *testing.T) {
	floors := make([]map[string]bool, 4)
	floors[0] = map[string]bool{"HH-G": true, "HH-M": true, "LL-M": true}
	floors[1] = map[string]bool{"LL-G": true}
	floors[2] = map[string]bool{}
	floors[3] = map[string]bool{}

	elevator := map[string]bool{}

	level := 0

	mvs := getAllMoves(floors, elevator, level)

	t.Log(mvs)
}

func TestValidMoves(t *testing.T) {
	floors := make([]map[string]bool, 4)
	floors[0] = map[string]bool{"LL-M": true}
	floors[1] = map[string]bool{}
	floors[2] = map[string]bool{"HH-G": true, "HH-M": true, "LL-G": true}
	floors[3] = map[string]bool{}

	elevator := map[string]bool{}

	level := 2

	mvs := getValidMoves(floors, elevator, level)

	t.Log(mvs)
}

func TestIsComplete(t *testing.T) {
	floors := make([]map[string]bool, 4)
	floors[0] = map[string]bool{}
	floors[1] = map[string]bool{}
	floors[2] = map[string]bool{"HH-G": true}
	floors[3] = map[string]bool{"HH-M": true}

	elevator := map[string]bool{"LL-G": true, "LL-M": true}

	level := 3

	c := isComplete(floors, elevator, level)
	t.Log(c)
}

func TestStates(t *testing.T) {
	// floors := make([]map[string]bool, 4)
	// floors[0] = map[string]bool{"LL-M": true, "HH-M": true}
	// floors[1] = map[string]bool{"HH-G": true}
	// floors[2] = map[string]bool{"LL-G": true}
	// floors[3] = map[string]bool{}

	// elevator := map[string]bool{}

	// mvs := []struct {
	// 	floor      []string
	// 	elevator   []string
	// 	start, end int
	// }{
	// 	{
	// 		elevator: []string{"HH-M"},
	// 		start:    0,
	// 		end:      1,
	// 	},
	// 	{
	// 		elevator: []string{"HH-G"},
	// 		start:    1,
	// 		end:      2,
	// 	},
	// 	{
	// 		floor: []string{"HH-G"},
	// 		start: 2,
	// 		end:   1,
	// 	},
	// 	{
	// 		start: 1,
	// 		end:   0,
	// 	},
	// 	{
	// 		elevator: []string{"LL-M"},
	// 		start:    0,
	// 		end:      1,
	// 	},
	// 	{
	// 		start: 1,
	// 		end:   2,
	// 	},
	// 	{
	// 		start: 2,
	// 		end:   3,
	// 	},
	// 	{
	// 		floor: []string{"LL-M"},
	// 		start: 3,
	// 		end:   2,
	// 	},
	// 	{
	// 		floor:    []string{"HH-M"},
	// 		elevator: []string{"HH-G", "LL-G"},
	// 		start:    2,
	// 		end:      3,
	// 	},
	// 	{
	// 		floor:    []string{"HH-G", "LL-G"},
	// 		elevator: []string{"LL-M"},
	// 		start:    3,
	// 		end:      2,
	// 	},
	// 	{
	// 		elevator: []string{"HH-M"},
	// 		start:    2,
	// 		end:      3,
	// 	},
	// }

	// nst := newState(floors, elevator, 0)
	// keys := []string{}

	// keys = append(keys, getKey(nst.floors, nst.elevator, nst.currentLevel))

	// for _, smv := range mvs {

	// 	fakemvs := getValidMoves(nst.floors, nst.elevator, nst.currentLevel)

	// 	contains := false
	// 	for _, mv := range fakemvs {
	// 		ff := sortedKeys(mv.floor)
	// 		fe := sortedKeys(mv.elevator)

	// 		ffs := strings.Join(ff, ",")
	// 		sfs := strings.Join(smv.floor, ",")
	// 		localf := ffs == sfs

	// 		fes := strings.Join(fe, ",")
	// 		ses := strings.Join(smv.elevator, ",")
	// 		locale := fes == ses

	// 		if localf && locale {
	// 			t.Log(ffs, sfs, "---", fes, ses)
	// 			contains = true
	// 		}

	// 	}

	// 	if !contains {
	// 		t.Log(smv, "fails")
	// 		t.Fail()
	// 	}

	// 	mv := move{
	// 		floor:    maplist(smv.floor),
	// 		elevator: maplist(smv.elevator),
	// 	}
	// 	keys = append(keys, getKey(nst.floors, nst.elevator, nst.currentLevel))

	// 	nst = transitState(nst, mv, smv.start, smv.end)
	// 	sm := stateMove{mv: mv, from: smv.start, to: smv.end}
	// 	if !validState(nst, sm) {
	// 		t.Log("not valid state", sm)
	// 		t.Fail()
	// 	} else {
	// 		t.Log("valid state", sm)
	// 	}
	// 	recordMove(nst, sm)

	// 	keys = append(keys, getKey(nst.floors, nst.elevator, nst.currentLevel))
	// }
	// t.Log("\n", strings.Join(keys, "\n"))
	// t.Log(len(nst.moves), nst.moves[:len(nst.moves)])
}
