package main

import (
	"fmt"
	"strconv"
)

type complexFlag complex128

func (c *complexFlag) String() string {
	return fmt.Sprint(complex128(*c))
}

func (c *complexFlag) Set(s string) error {
	var v complex128

	switch s {
	case "a", "A":
		v = pointA
	case "b", "B":
		v = pointB
	default:
		var err error
		if v, err = strconv.ParseComplex(s, 128); err != nil {
			return err
		}
	}

	*c = complexFlag(v)
	return nil
}
