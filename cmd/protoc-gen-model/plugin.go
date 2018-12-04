package main

import (
	extrapb "github.com/fzerorubigd/protobuf/extra"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/gogo/protobuf/vanity"
	"github.com/kr/pretty"
)

type plugin struct {
	*generator.Generator
	generator.PluginImports
	useGogoImport bool
}

func newPlugin(useGogoImport bool) generator.Plugin {
	return &plugin{
		useGogoImport: useGogoImport,
	}
}

func (p *plugin) Name() string {
	return "model"
}

func (p *plugin) Init(g *generator.Generator) {
	p.Generator = g
}

func (p *plugin) Generate(file *generator.FileDescriptor) {
	if !p.useGogoImport {
		vanity.TurnOffGogoImport(file.FileDescriptorProto)
	}
	p.PluginImports = generator.NewPluginImports(p.Generator)

	p.P("/*")
	for _, msg := range file.Messages() {
		//if msg.Options.

		//msg.DescriptorProto.GetOptions().GetMapEntry()
		//I := proto.ExtensionDesc(*extrapb.E_IsModel)
		ex, err := proto.GetExtension(msg.Options, extrapb.E_SchemaName)
		if err != nil {
			p.P(err.Error())
		} else {
			p.P(pretty.Sprint(ex))
		}

		//for _, opt := range msg. {
		//p.P(pretty.Sprint(msg.Options.XXX_InternalExtensions))
	}
	//}
	p.P("*/")

}
