package common

import "fmt"

type CoordXYZ struct {
	X int32
	Y int32
	Z int32
}

func (c CoordXYZ) String() string {
	return fmt.Sprintf("%d,%d,%d", c.X, c.Y, c.Z)
}
