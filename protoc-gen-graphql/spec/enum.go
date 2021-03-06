package spec

import (
	"strings"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// Enum spec wraps EnumDescriptorProto with keeping file definition.
type Enum struct {
	descriptor *descriptor.EnumDescriptorProto
	*File

	paths  []int
	values []*EnumValue

	*Dependencies
}

func NewEnum(
	d *descriptor.EnumDescriptorProto,
	f *File,
	paths ...int,
) *Enum {

	e := &Enum{
		descriptor: d,
		File:       f,
		paths:      paths,
		values:     make([]*EnumValue, 0),

		Dependencies: NewDependencies(),
	}
	for i, v := range d.GetValue() {
		ps := make([]int, len(paths))
		copy(ps, paths)
		e.values = append(e.values, NewEnumValue(v, f, append(ps, 2, i)...))
	}
	return e
}

func (e *Enum) Comment() string {
	return e.File.getComment(e.paths)
}

func (e *Enum) Name() string {
	return e.descriptor.GetName()
}

func (e *Enum) SingleName() string {
	spl := strings.Split(e.Name(), ".")
	return spl[len(spl)-1]
}

func (e *Enum) Values() []*EnumValue {
	return e.values
}

func (e *Enum) FullPath() string {
	return e.File.Package() + "." + e.Name()
}
