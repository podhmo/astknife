package patchwork4

import (
	"fmt"
	"go/ast"
	"go/token"
	"sort"

	"github.com/podhmo/astknife/patchwork4/buildtree"
)

// ToFile :
func ToFile(p *Patchwork, filename string) *ast.File {
	base := p.Fset.Base()
	fmt.Println("root.start", base)
	s := &buildtree.State{
		Fset:        p.Fset,
		Debug:       p.Debug,
		Replacing:   p.Replacing,
		Appending:   p.Appending,
		Lines:       map[int]int{},
		RegionStack: []*buildtree.Region{buildtree.NewRegion(p.File, base)},
		FileBase:    base,
		Base:        base,
	}
	f := &ast.File{
		Name:    buildtree.Ident(p.File.Name, s), // xxx
		Scope:   ast.NewScope(nil),
		Package: token.Pos(base),
	}
	s.Base = int(f.Name.End())
	fmt.Println("root.start2", s.Base)
	f.Imports = buildtree.ImportSpecs(p.File.Imports, s)
	f.Decls = buildtree.Decls(p.File.Decls, s)

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

	tokenf := p.Fset.AddFile(filename, -1, int(f.End()-f.Pos()))
	// todo: new line
	lines := make([]int, len(s.Lines))
	for i, pos := range lines {
		lines[i] = pos
	}
	tokenf.SetLines(lines)
	tokenf.SetLines([]int{0, 10})
	fmt.Println("root.end", f.End()+1)
	return f
}
