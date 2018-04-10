package patchwork5

import (
	"fmt"
	"go/ast"
)

func parseASTFile(f *ast.File) *File {
	p := &parserFromAST{f: f}
	p.parseDecls(f.Decls)
	return &File{
		Regions: p.regions,
	}
}

type parserFromAST struct {
	f       *ast.File
	regions []*Region
	current int
}

func (p *parserFromAST) parseDecls(decls []ast.Decl) {
	for i := range decls {
		p.parseDecl(decls[i])
	}
}

func (p *parserFromAST) parseDecl(decl ast.Decl) {
	if decl == nil {
		return
	}
	switch x := decl.(type) {
	case *ast.GenDecl:
		p.regions = append(p.regions, &Region{Ref: &DeclHeadRef{decl: x}})
		p.parseSpecs(x.Specs)
		p.regions = append(p.regions, &Region{Ref: &DeclTailRef{decl: x}})
	case *ast.FuncDecl:
		p.regions = append(p.regions, &Region{Ref: &DeclRef{Decl: x, name: x.Name.Name}})
	case *ast.BadDecl:
		p.regions = append(p.regions, &Region{Ref: &DeclRef{Decl: x}})
	default:
		panic(fmt.Sprintf("invalid decl %q", x))
	}
}

func (p *parserFromAST) parseSpecs(specs []ast.Spec) {
	for i := range specs {
		p.parseSpec(specs[i])
	}
}

func (p *parserFromAST) parseSpec(spec ast.Spec) {
	if spec == nil {
		return
	}
	switch x := spec.(type) {
	case *ast.ImportSpec:
		p.regions = append(p.regions, &Region{Ref: &SpecRef{Spec: x}})
	case *ast.ValueSpec:
		p.regions = append(p.regions, &Region{Ref: &SpecRef{Spec: x}})
	case *ast.TypeSpec:
		p.regions = append(p.regions, &Region{Ref: &SpecRef{Spec: x, name: x.Name.Name}})
	default:
		panic(fmt.Sprintf("invalid spec %q", x))
	}
}

func dumpRegions(regions []*Region) {
	for i, r := range regions {
		fmt.Println(i, r)
	}
}
