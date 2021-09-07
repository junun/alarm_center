package inject

import (
	"log"

	"go.uber.org/dig"
)
// DI 服务di容器
type DI struct {
	container *dig.Container
}

// New create di entry
func New() *DI {
	return &DI{
		container: dig.New(),
	}
}

// BuildProvides build container provides.
func (d *DI) BuildProvides(provides ...interface{}) {
	var err error
	for k := range provides {
		err = d.container.Provide(provides[k])
		if err != nil {
			log.Fatalln("provide error: ", err)
		}
	}
}

// Invoke invoke provides.
func (d *DI) Invoke(fn interface{}) error {
	return d.container.Invoke(fn)
}
