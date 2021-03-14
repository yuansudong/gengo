package gengo

import (
	"fmt"
	"strings"

	descriptor "github.com/yuansudong/gengo/descriptor"
)

// Message 用于描述一个消息结构体
type Message struct {
	// File 描述该消息属于哪个文件
	File *File
	// Outers 消息嵌套列表
	Outers []string
	*descriptor.DescriptorProto
	Fields []*Field
	Index  int
}

// _LookupField 根据名称查找字段
func (m *Message) LookupField(name string) *Field {
	for _, f := range m.Fields {
		if f.GetName() == name {
			return f
		}
	}
	return nil
}

// FQMN 返回一个完整的消息名称
func (m *Message) FQMN() string {
	components := []string{""}
	if m.File.Package != nil {
		components = append(components, m.File.GetPackage())
	}
	components = append(components, m.Outers...)
	components = append(components, m.GetName())
	return strings.Join(components, ".")
}

// GoType 用于返回一个Go类型
func (m *Message) GoType(currentPackage string) string {
	var components []string
	components = append(components, m.Outers...)
	components = append(components, m.GetName())

	name := strings.Join(components, "_")
	if m.File.GoPkg.Path == currentPackage {
		return name
	}
	pkg := m.File.GoPkg.Name
	if alias := m.File.GoPkg.Alias; alias != "" {
		pkg = alias
	}
	return fmt.Sprintf("%s.%s", pkg, name)
}
