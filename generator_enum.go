package gengo

import (
	"fmt"
	"strings"

	descriptor "github.com/yuansudong/gengo/descriptor"
)

// Enum 描述protobuf中的枚举类型
type Enum struct {
	// File 这个枚举类型在哪个文件中
	File *File
	// Outers is a list of outer messages if this enum is a nested type.
	Outers []string
	*descriptor.EnumDescriptorProto
	Index int
}

// FQEN 返回完整的枚举名称.
func (e *Enum) FQEN() string {
	components := []string{""}
	if e.File.Package != nil {
		components = append(components, e.File.GetPackage())
	}
	components = append(components, e.Outers...)
	components = append(components, e.GetName())
	return strings.Join(components, ".")
}

// GoType 返回一个与go相关的类型
func (e *Enum) GoType(currentPackage string) string {
	var components []string
	components = append(components, e.Outers...)
	components = append(components, e.GetName())

	name := strings.Join(components, "_")
	if e.File.GoPkg.Path == currentPackage {
		return name
	}
	pkg := e.File.GoPkg.Name
	if alias := e.File.GoPkg.Alias; alias != "" {
		pkg = alias
	}
	return fmt.Sprintf("%s.%s", pkg, name)
}
