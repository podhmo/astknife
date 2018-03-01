package patchwork2

import (
	"fmt"
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
	tokenf   *token.File
	base     int
	comments []*ast.CommentGroup
}

func (a *aggregator) setBase(pos token.Pos) {
	fmt.Println("!! update base", a.base, "->", int(pos))
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
	tokenf := fset.AddFile(filename, -1, size)
	base := tokenf.Base()
	f := &ast.File{
		Name:    mirror.Ident(fref.File.Name, base, false), // xxx
		Scope:   ast.NewScope(nil),
		Package: token.Pos(base),
	}

	a := &aggregator{base: int(f.Name.End()), f: f, fref: fref, tokenf: tokenf}
	for _, dref := range fref.Decls {
		if len(dref.Specs) == 0 {
			decl := a.aggregateDeclRef(dref)
			f.Decls = append(f.Decls, decl)
		} else {
			decl := a.aggregateGencDeclRef(dref)
			f.Decls = append(f.Decls, decl)
		}
	}
	f.Comments = a.comments
	return f
}

// aggregateGencDeclRef
func (a *aggregator) aggregateGencDeclRef(dref *DeclRef) ast.Decl {
	if dref.Replacement != nil {
		offset := int(-dref.Replacement.Pos()) + a.base
		if len(dref.Comments) > 0 {
			offset += int(dref.Replacement.Pos() - dref.Comments[0].Pos())
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
		// xxxx :
		// if len(dref.Comments) > 0 {
		// 	base += (dref.Comments[len(dref.Comments)-1].End() - dref.Replacement.End())
		// }
		new.Specs = specs
		if new.Lparen != token.NoPos {
			new.Lparen = token.Pos(int(new.Lparen) + offset)
		}
		if new.Rparen != token.NoPos {
			new.Rparen = token.Pos(int(new.Rparen) + offset)
		}
		new.Doc = mirror.CommentGroup(new.Doc, offset, false)
		return &new
	}

	if dref.Original != nil {
		offset := int(-dref.Original.Pos()) + a.base
		fmt.Println("**GenDecl***", "@offset", offset, "@base", a.base)
		if len(dref.Comments) > 0 {
			offset += int(dref.Original.Pos() - dref.Comments[0].Pos())
			a.comments = append(a.comments, moveComments(dref.Comments, offset)...)
		}

		decl := dref.Original.(*ast.GenDecl)
		new := *decl
		new.TokPos = token.Pos(int(new.TokPos) + offset)
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
		if new.Lparen != token.NoPos {
			new.Lparen = token.Pos(int(new.Lparen) + offset)
		}
		if new.Rparen != token.NoPos {
			new.Rparen = token.Pos(int(new.Rparen) + offset)
		}
		new.Doc = mirror.CommentGroup(new.Doc, offset, false)
		return &new
	}

	panic("something wrong")
}

// aggregateDeclRef
func (a *aggregator) aggregateDeclRef(dref *DeclRef) ast.Decl {
	if dref.Replacement != nil {
		offset := int(-dref.Replacement.Pos()) + a.base
		if len(dref.Comments) > 0 {
			offset += int(dref.Replacement.Pos() - dref.Comments[0].Pos())
			a.comments = append(a.comments, moveComments(dref.Comments, offset)...)
		}
		decl := mirror.Decl(dref.Replacement, offset, false)
		base := decl.End()
		if len(dref.Comments) > 0 {
			base += (dref.Comments[len(dref.Comments)-1].End() - dref.Replacement.End())
		}
		a.setBase(base)
		return decl
	}

	if dref.Original != nil {
		offset := int(-dref.Original.Pos()) + a.base
		fmt.Println("**FuncDecl***", "@offset", offset, "@base", a.base)
		if len(dref.Comments) > 0 {
			offset += int(dref.Original.Pos() - dref.Comments[0].Pos())
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
			offset += int(sref.Replacement.Pos() - sref.Comments[0].Pos())
			// a.comments = append(a.comments, moveComments(sref.Comments, offset)...)
		}

		spec := mirror.Spec(sref.Replacement, offset, false)
		base := spec.End()
		// if len(sref.Comments) > 0 {
		// 	base += (sref.Comments[len(sref.Comments)-1].End() - sref.Replacement.End())
		// }
		a.setBase(base)
		return spec
	}

	if sref.Original != nil {
		offset := int(-sref.Original.Pos()) + a.base
		fmt.Println("**SpeC***", "@offset", offset, "@base", a.base)
		// if len(sref.Comments) > 0 {
		// 	a.comments = append(a.comments, moveComments(sref.Comments, offset)...)
		// }
		spec := mirror.Spec(sref.Original, offset, false)
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
	size := int(fref.File.End() - fref.File.Pos()) // xxx;
	cmapMap := map[*ast.File]ast.CommentMap{}
	cmap := ast.NewCommentMap(fset, fref.File, fref.File.Comments)

	t := &trimmer{fset: fset, fref: fref, cmapMap: cmapMap, cmap: cmap}
	for _, dref := range fref.Decls {
		fmt.Println("@SIZE@", size)
		size += t.trimDeclRef(dref)
	}
	fref.Comments = cmap.Comments()
	if len(fref.Comments) > 0 {
		cend := fref.Comments[len(fref.Comments)-1].End()
		if cend > fref.File.End() {
			size += int(cend - fref.File.End())
		}
	}
	fmt.Println("@SIZE@@@", size)
	return size + 2
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
		delta += int(dref.Replacement.End() - dref.Replacement.Pos())
		added := t.getCmap(dref.File).Filter(dref.Replacement)
		dref.Comments = added.Comments()

		if len(dref.Comments) > 0 {
			delta += int(dref.Replacement.Pos() - dref.Comments[0].Pos())
			delta += int(dref.Comments[len(dref.Comments)-1].End() - dref.Replacement.End())
		}

		if dref.Original != nil {
			delta -= int(dref.Original.End() - dref.Original.Pos())
			removed := t.cmap.Filter(dref.Original)
			for k := range removed {
				delete(t.cmap, k)
			}
			removedComments := removed.Comments()
			if len(removedComments) > 0 {
				delta -= int(dref.Replacement.Pos() - removedComments[0].Pos())
				delta -= int(removedComments[len(removedComments)-1].End() - dref.Replacement.End())
			}
		}
	} else if dref.Original != nil {
		added := t.getCmap(dref.File).Filter(dref.Original)
		dref.Comments = added.Comments()
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
		delta += int(sref.Replacement.End() - sref.Replacement.Pos())
		added := t.getCmap(sref.File).Filter(sref.Replacement)
		sref.Comments = added.Comments()
		if len(sref.Comments) > 0 {
			delta += int(sref.Replacement.Pos() - sref.Comments[0].Pos())
			delta += int(sref.Comments[len(sref.Comments)-1].End() - sref.Replacement.End())
		}

		if sref.Original != nil {
			delta -= int(sref.Original.End() - sref.Original.Pos())
			removed := t.cmap.Filter(sref.Original)
			for k := range removed {
				delete(t.cmap, k)
			}
			removedComments := removed.Comments()
			if len(removedComments) > 0 {
				delta -= int(sref.Replacement.Pos() - removedComments[0].Pos())
				delta -= int(removedComments[len(removedComments)-1].End() - sref.Replacement.End())
			}
		}
	} else if sref.Original != nil {
		added := t.getCmap(sref.File).Filter(sref.Original)
		sref.Comments = added.Comments()
	}
	return delta
}
