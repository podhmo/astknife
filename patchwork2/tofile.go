package patchwork2

import (
	"go/ast"
	"go/token"

	"github.com/podhmo/astknife/patchwork2/mirror"
)

// ToFile :
func (fref *FileRef) ToFile(fset *token.FileSet, filename string) *ast.File {
	// calculate size and trim comments
	size := trim(fset, fref)
	// aggregate file (create new AST and fix positions)
	return aggregate(fset, fref, filename, size)
}

type aggregator struct {
	f        *ast.File
	fref     *FileRef
	base     int
	comments []*ast.CommentGroup
}

func (a *aggregator) setBase(pos token.Pos) {
	a.base = int(pos)
}
func moveComments(xs []*ast.CommentGroup, offset int) []*ast.CommentGroup {
	ys := make([]*ast.CommentGroup, len(xs))
	for i, x := range xs {
		ys[i] = mirror.CommentGroup(x, offset, false)
	}
	return ys
}

func aggregate(fset *token.FileSet, fref *FileRef, filename string, size int) *ast.File {
	base := fset.AddFile(filename, -1, size).Base()
	f := &ast.File{
		Name:    mirror.Ident(fref.File.Name, base, false),
		Scope:   ast.NewScope(nil),
		Package: token.Pos(base),
	}

	a := &aggregator{base: int(f.Name.End() + 1), f: f, fref: fref}
	for _, dref := range fref.Decls {
		if len(dref.Specs) == 0 {
			decl := a.aggregateDeclRef(dref)
			f.Decls = append(f.Decls, decl)
		} else {
			decl := a.aggregateGencDeclRef(dref)
			f.Decls = append(f.Decls, decl)
		}
	}
	f.Comments = append(a.comments, moveComments(fref.Comments, a.base)...)
	return f
}

// aggregateGencDeclRef
func (a *aggregator) aggregateGencDeclRef(dref *DeclRef) ast.Decl {
	if dref.Replacement != nil {
		rpos := dref.Replacement.Pos()
		offset := int(-rpos) + a.base
		if len(dref.Comments) > 0 {
			for _, c := range dref.Comments {
				pos := c.Pos()
				if pos < rpos {
					offset += int(rpos - pos)
				}
			}
			a.comments = append(a.comments, moveComments(dref.Comments, offset)...)
		}

		decl := dref.Replacement.(*ast.GenDecl)
		new := *decl

		// xxx:
		if decl.Lparen.IsValid() {
			new.Lparen = token.Pos(int(new.Lparen) + offset)
			// <token> ( spec, ... )
			a.setBase(token.Pos(a.base + int(decl.Lparen-decl.Pos()) + len(decl.Tok.String())))
			defer a.setBase(token.Pos(a.base + int(decl.End()-decl.Rparen)))
		}
		specs := make([]ast.Spec, len(decl.Specs))
		for i, sref := range dref.Specs {
			specs[i] = a.aggregateSpecRef(sref)
		}
		new.Specs = specs
		return &new
	}

	if dref.Original != nil {
		offset := int(-dref.File.Pos()) + a.base
		if len(dref.Comments) > 0 {
			a.comments = append(a.comments, moveComments(dref.Comments, offset)...)
		}

		decl := dref.Original.(*ast.GenDecl)
		new := *decl

		// xxx:
		if decl.Lparen.IsValid() {
			new.Lparen = token.Pos(int(new.Lparen) + offset)
			// <token> ( spec, ... )
			a.setBase(token.Pos(a.base + int(decl.Lparen-decl.Pos()) + len(decl.Tok.String())))
			defer a.setBase(token.Pos(a.base + int(decl.End()-decl.Rparen)))
		}
		specs := make([]ast.Spec, len(decl.Specs))
		for i, sref := range dref.Specs {
			specs[i] = a.aggregateSpecRef(sref)
		}
		new.Specs = specs
		return &new
	}

	panic("something wrong")
}

// aggregateDeclRef
func (a *aggregator) aggregateDeclRef(dref *DeclRef) ast.Decl {
	if dref.Replacement != nil {
		rpos := dref.Replacement.Pos()
		offset := int(-rpos) + a.base
		if len(dref.Comments) > 0 {
			for _, c := range dref.Comments {
				pos := c.Pos()
				if pos < rpos {
					offset += int(rpos - pos)
				}
			}
			a.comments = append(a.comments, moveComments(dref.Comments, offset)...)
		}

		decl := mirror.Decl(dref.Replacement, offset, false)
		a.setBase(decl.End())
		return decl
	}

	if dref.Original != nil {
		offset := int(-dref.File.Pos()) + a.base
		if len(dref.Comments) > 0 {
			a.comments = append(a.comments, moveComments(dref.Comments, offset)...)
		}
		decl := mirror.Decl(dref.Original, offset, false)
		a.setBase(decl.End())
		return decl
	}

	panic("something wrong")
}

// aggregateSpecRef
func (a *aggregator) aggregateSpecRef(sref *SpecRef) ast.Spec {
	if sref.Replacement != nil {
		rpos := sref.Replacement.Pos()
		offset := int(-rpos) + a.base
		if len(sref.Comments) > 0 {
			for _, c := range sref.Comments {
				pos := c.Pos()
				if pos < rpos {
					offset += int(rpos - pos)
				}
			}
			a.comments = append(a.comments, moveComments(sref.Comments, offset)...)
		}

		spec := mirror.Spec(sref.Replacement, offset, false)
		a.setBase(spec.End())
		return spec
	}

	if sref.Original != nil {
		offset := int(-sref.File.Pos()) + a.base
		if len(sref.Comments) > 0 {
			a.comments = append(a.comments, moveComments(sref.Comments, offset)...)
		}
		spec := mirror.Spec(sref.Replacement, offset, false)
		a.setBase(spec.End())
		return spec
	}

	panic("something wrong")
}

type trimmer struct {
	fset    *token.FileSet
	fref    *FileRef
	cmap    ast.CommentMap
	cmapMap map[*ast.File]ast.CommentMap
}

// trim :
func trim(fset *token.FileSet, fref *FileRef) int {
	size := int(fref.File.End() - fref.File.Pos())
	cmapMap := map[*ast.File]ast.CommentMap{}
	cmap := ast.NewCommentMap(fset, fref.File, fref.File.Comments)

	t := &trimmer{fset: fset, fref: fref, cmapMap: cmapMap, cmap: cmap}
	for _, dref := range fref.Decls {
		size += t.trimDeclRef(dref)
	}
	fref.Comments = cmap.Comments()
	return size
}

func (t *trimmer) getCmap(f *ast.File) ast.CommentMap {
	cmap, ok := t.cmapMap[f]
	if !ok {
		cmap = ast.NewCommentMap(t.fset, f, f.Comments)
		t.cmapMap[f] = cmap
	}
	return cmap
}

// trimDeclRef  :
func (t *trimmer) trimDeclRef(dref *DeclRef) int {
	delta := 0
	if dref.Replacement != nil {
		added := t.getCmap(dref.File).Filter(dref.Replacement)
		for _, cs := range added {
			for _, c := range cs {
				delta += int(c.End() - c.Pos())
			}
		}
		dref.Comments = added.Comments()

		if dref.Original != nil {
			removed := t.cmap.Filter(dref.Original)
			for k := range removed {
				delete(t.cmap, k)
			}
			for _, cs := range removed {
				for _, c := range cs {
					delta -= int(c.End() - c.Pos())
				}
			}
		}
	}
	for _, sref := range dref.Specs {
		delta += t.trimSpecRef(sref)
	}
	return delta
}

// trimSpecRef  :
func (t *trimmer) trimSpecRef(sref *SpecRef) int {
	delta := 0
	if sref.Replacement != nil {
		added := t.getCmap(sref.File).Filter(sref.Replacement)
		for _, cs := range added {
			for _, c := range cs {
				delta += int(c.End() - c.Pos())
			}
		}
		sref.Comments = added.Comments()

		if sref.Original != nil {
			removed := t.cmap.Filter(sref.Original)
			for k := range removed {
				delete(t.cmap, k)
			}
			for _, cs := range removed {
				for _, c := range cs {
					delta -= int(c.End() - c.Pos())
				}
			}
		}
	}
	return delta
}
