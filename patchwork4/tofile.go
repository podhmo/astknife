package patchwork4

import (
	"fmt"
	"go/ast"
	"go/token"
	"sort"

	"github.com/podhmo/astknife/patchwork4/mirror"
)

// ToFile :
func ToFile(p *Patchwork, filename string) *ast.File {
	base := p.Fset.Base()
	s := &mirror.State{
		Replacing:   p.Replacing,
		Appending:   p.Appending,
		RegionStack: []*mirror.Region{mirror.NewRegion(p.File, base)},
		Base:        base,
	}
	f := &ast.File{
		Name:    mirror.Ident(p.File.Name, s), // xxx
		Scope:   ast.NewScope(nil),
		Package: token.Pos(base),
	}
	s.Base = int(f.Name.End()) + 3
	f.Imports = mirror.ImportSpecs(p.File.Imports, s)
	f.Decls = mirror.Decls(p.File.Decls, s)

	var comments []*ast.CommentGroup
	ast.Inspect(f, func(node ast.Node) bool {
		if node != nil {
			if c, ok := node.(*ast.CommentGroup); ok {
				comments = append(comments, c)
			}
		}
		return true
	})
	sort.Slice(comments, func(i, j int) bool { return comments[i].Pos() < comments[j].Pos() })
	f.Comments = comments // xxx

	// todo: new line

	fmt.Println("******* before", p.File.End(), p.File.Pos(), int(p.File.End()-p.File.Pos()))
	fmt.Println("******* after", f.End(), f.Pos(), int(f.End()-f.Pos()))
	p.Fset.AddFile(filename, -1, int(f.End()-f.Pos()))
	return f
}
