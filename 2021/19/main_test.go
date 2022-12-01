package main

import (
	"fmt"
	"github.com/jasontconnell/advent/common"
	"testing"
)

func TestRotate(t *testing.T) {
	p := xyz{1, 2, 3}

	pts := map[xyz]matrix{}
	mx := []matrix{
		{x0, y0, z0}, {x90, y0, z0}, {x180, y0, z0}, {x270, y0, z0},
		{x0, y90, z0}, {x90, y90, z0}, {x180, y90, z0}, {x270, y90, z0},
		{x0, y180, z0}, {x90, y180, z0}, {x180, y180, z0}, {x270, y180, z0},
		{x0, y270, z0}, {x90, y270, z0}, {x180, y270, z0}, {x270, y270, z0},

		{x0, y0, z90}, {x90, y0, z90}, {x180, y0, z90}, {x270, y0, z90},
		{x0, y90, z90}, {x90, y90, z90}, {x180, y90, z90}, {x270, y90, z90},
		{x0, y180, z90}, {x90, y180, z90}, {x180, y180, z90}, {x270, y180, z90},
		{x0, y270, z90}, {x90, y270, z90}, {x180, y270, z90}, {x270, y270, z90},
		{x0, y0, z90}, {x90, y0, z90}, {x180, y0, z90}, {x270, y0, z90},
		{x0, y90, z180}, {x90, y90, z180}, {x180, y90, z180}, {x270, y90, z180},
		{x0, y90, z180}, {x90, y90, z180}, {x180, y90, z180}, {x270, y180, z180},
		{x0, y90, z180}, {x90, y90, z180}, {x180, y90, z180}, {x270, y180, z180},
		{x0, y180, z270}, {x90, y180, z270}, {x180, y180, z270}, {x270, y180, z270},
		{x0, y270, z270}, {x90, y270, z270}, {x180, y270, z270}, {x270, y270, z270},

		{x0, y0, z90}, {x90, y0, z90}, {x180, y0, z90}, {x270, y0, z90},
		{x0, y90, z180}, {x90, y90, z180}, {x180, y90, z180}, {x270, y90, z180},
		{x0, y180, z270}, {x90, y180, z270}, {x180, y180, z270}, {x270, y180, z270},
		{x0, y270, z270}, {x90, y270, z270}, {x180, y270, z270}, {x270, y270, z270},
	}

	for _, m := range mx {
		r := p
		r = rotateX(r, m.x)
		r = rotateY(r, m.y)
		r = rotateZ(r, m.z)

		if m2, ok := pts[r]; ok {
			fmt.Println("pt already added", r, "by", m2, "tried", m)
		}
		pts[r] = m
	}

	fmt.Println(len(pts))
	for _, v := range pts {
		fmt.Print("{" + v.String() + "},")
	}
}

func TestRotateMatrix(t *testing.T) {
	p := xyz{3, 4, 5}

	rs := []xyz{
		{p.x, p.y, p.z}, {p.x, p.y, -p.z}, {p.x, -p.y, -p.z}, {p.x, -p.y, p.z},
		{-p.x, p.y, p.z}, {-p.x, p.y, -p.z}, {-p.x, -p.y, -p.z}, {-p.x, -p.y, p.z},

		{p.y, p.x, p.z}, {p.y, p.x, -p.z}, {-p.y, p.x, -p.z}, {-p.y, p.x, p.z},
		{p.y, -p.x, p.z}, {p.y, -p.x, -p.z}, {-p.y, -p.x, -p.z}, {-p.y, -p.x, p.z},

		{p.z, p.y, p.x}, {-p.z, p.y, p.x}, {-p.z, -p.y, p.x}, {p.z, -p.y, p.x},
		{p.z, p.y, -p.x}, {-p.z, p.y, -p.x}, {-p.z, -p.y, -p.x}, {p.z, -p.y, -p.x},

		{p.z, p.x, p.y}, {-p.z, p.x, p.y}, {-p.z, p.x, -p.y}, {p.z, p.x, -p.y},
		{p.z, -p.x, p.y}, {-p.z, -p.x, p.y}, {-p.z, -p.x, -p.y}, {p.z, -p.x, -p.y},
	}

	fmt.Println(rs)
}

func TestScratch(t *testing.T) {
	x, y, z := 1, 2, 3

	pts := []xyz{
		{x, y, z}, {x, z, y}, {y, x, z}, {y, z, x}, {z, x, y}, {z, y, x},
		{-x, y, z}, {-x, z, y}, {y, -x, z}, {y, z, -x}, {z, -x, y}, {z, y, -x},
		{x, -y, z}, {x, z, -y}, {-y, x, z}, {-y, z, x}, {z, x, -y}, {z, -y, x},
		{x, y, -z}, {x, -z, y}, {y, x, -z}, {y, -z, x}, {-z, x, y}, {-z, y, x},
	}

	t.Log(pts)
}

func TestExample(t *testing.T) {
	in, err := common.ReadStrings("example.txt")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	scanners := parseInput(in)

	rd := getRotationDelta(scanners[4], scanners[1], 12)
	fmt.Println(rd)
	rd = getRotationDelta(scanners[3], scanners[1], 12)
	fmt.Println(rd)

	// getRotationDelta(scanners[0], scanners[1], 12)
	// getRotationDelta(scanners[0], scanners[2], 12)
	// getRotationDelta(scanners[0], scanners[3], 12)
	// getRotationDelta(scanners[0], scanners[4], 12)

	// getRotationDelta(scanners[1], scanners[2], 12)
	// getRotationDelta(scanners[1], scanners[3], 12)
	// getRotationDelta(scanners[1], scanners[4], 12)

	// getRotationDelta(scanners[2], scanners[3], 12)
	// getRotationDelta(scanners[2], scanners[4], 12)

	// getRotationDelta(scanners[3], scanners[4], 12)

	for _, s := range scanners {
		fmt.Println(s.id, len(s.points))
	}

	// if !locateAllScanners(scanners, 12) {
	// 	t.Fail()
	// }
}
