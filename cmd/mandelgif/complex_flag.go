package main

import (
	"fmt"
	"strconv"
	"strings"
)

type complexFlag complex128

func (c *complexFlag) String() string {
	return fmt.Sprint(complex128(*c))
}

func (c *complexFlag) Set(s string) error {
	var v complex128

	switch strings.ToUpper(s) {
	case "A":
		v = pointA
	case "B":
		v = pointB
	case "C":
		v = pointC
	case "D":
		v = pointD
	case "E":
		v = pointE
	case "F":
		v = pointF
	default:
		var err error
		if v, err = strconv.ParseComplex(s, 128); err != nil {
			return err
		}
	}

	*c = complexFlag(v)
	return nil
}
