package dataStructures

import "gitlab.com/elixxir/primitives/ndf"

type Ndf struct{
	f ndf.NetworkDefinition
	hash []byte
}


func(ndf *Ndf)Update(
