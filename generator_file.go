package gengo

import (
	"fmt"
	"strings"

	descriptor "github.com/yuansudong/gengo/descriptor"
)

// GoPackage 用于描述一个golang的包
type GoPackage struct {
	Path string // Path 用于描述一个包的路径
	// Name 用于描述一个包名
	Name string
	// Alias 包的别名
	Alias string
}

// Standard 用于判断是不是golang的标准包
func (p GoPackage) Standard() bool {
	return !strings.Contains(p.Path, ".")
}

// String 返回一个包的完整路径
func (p GoPackage) String() string {
	if p.Alias == "" {
		return fmt.Sprintf("%q", p.Path)
	}
	return fmt.Sprintf("%s %q", p.Alias, p.Path)
}

// File 用于扩展文件
type File struct {
	*descriptor.FileDescriptorProto
	// GoPkg go包的集合
	GoPkg GoPackage
	// Messages 定义在这个文件中的Message
	Messages []*Message
	// Enums 定义在这个文件里的枚举
	Enums []*Enum
	// Services 定义在这个文件里的服务
	Services []*Service
	Imports  map[string]bool
}

// AddImportByPublic 用于增加
func (f *File) AddImportByPublic(imp string) {
	if f.Imports == nil {
		f.Imports = make(map[string]bool)
	}
	f.Imports[imp] = true
}

// AddImportByMessage 用于向该文件中添加要导入的路径
func (f *File) AddImportByMessage(m *Message) {
	if f.Imports == nil {
		f.Imports = make(map[string]bool)
	}
	if f.GoPkg.Path != m.File.GoPkg.Path {
		f.Imports[m.File.GoPkg.Path] = true
	}
}

// AddImportByEnum 用于增加枚举导入的包
func (f *File) AddImportByEnum(e *Enum) {
	if f.Imports == nil {
		f.Imports = make(map[string]bool)
	}
	if f.GoPkg.Path != e.File.GoPkg.Path {
		f.Imports[e.File.GoPkg.Path] = true
	}
}

// proto2 判断该协议是不是proto2
func (f *File) proto2() bool {
	return f.Syntax == nil || f.GetSyntax() == "proto2"
}
