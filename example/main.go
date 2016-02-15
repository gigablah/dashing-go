package main

import (
  "gopkg.in/gigablah/dashing-go.v1"
	_  "gopkg.in/gigablah/dashing-go.v1/example/jobs"
)

func main() {
	d := dashing.NewDashing()
	d.Start()
}
