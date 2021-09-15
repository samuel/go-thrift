// This file is automatically generated. Do not modify.

package gentest

import (
	"fmt"
)

var _ = fmt.Sprintf

type NestedColor struct {
	Rgb *Rgb `thrift:"1,required" json:"rgb"`
}

func (n *NestedColor) GetRgb() (val Rgb) {
	if n != nil && n.Rgb != nil {
		return *n.Rgb
	}

	return
}

type Rgb struct {
	Red   *int32 `thrift:"1,required" json:"red"`
	Green *int32 `thrift:"2,required" json:"green"`
	Blue  *int32 `thrift:"3,required" json:"blue"`
}

func (r *Rgb) GetRed() (val int32) {
	if r != nil && r.Red != nil {
		return *r.Red
	}

	return
}

func (r *Rgb) GetGreen() (val int32) {
	if r != nil && r.Green != nil {
		return *r.Green
	}

	return
}

func (r *Rgb) GetBlue() (val int32) {
	if r != nil && r.Blue != nil {
		return *r.Blue
	}

	return
}
