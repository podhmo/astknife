package patchwork5

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"

	"github.com/k0kubun/pp"
)

func parseASTFile(fset *token.FileSet, f *ast.File) *File {
	tokenf := fset.File(f.Pos())
	reflines := reflect.ValueOf(tokenf).Elem().FieldByName("lines")
	lines := make([]int, reflines.Len())
	for i := range lines {
		lines[i] = int(reflines.Index(i).Int())
	}

	p := &parserFromAST{
		f:     f,
		base:  tokenf.Base(),
		lines: lines,
	}
	pp.ColoringEnabled = false
	pp.Println(f.Comments)
	p.parseDecls(f.Decls)
	end := f.End()
	if len(f.Comments) > 0 {
		cend := f.Comments[len(f.Comments)-1].End()
		if end < cend {
			end = cend
		}
	}
	p.parsePaddingComments(end)
	return &File{
		Regions: p.regions,
	}
}

type parserFromAST struct {
	f     *ast.File
	lines []int
	base  int

	regions    []*Region
	currentPos int
	commentIdx int
	linesIdx   int
}

func (p *parserFromAST) collectComments(end token.Pos) []*ast.CommentGroup {
	var comments []*ast.CommentGroup
	for i := p.commentIdx; i < len(p.f.Comments); {
		c := p.f.Comments[i]
		if end <= c.End() {
			break
		}
		comments = append(comments, c)
		i++
		p.commentIdx = i
	}
	return comments
}

func (p *parserFromAST) collectLines(pos, pend token.Pos) []int {
	// start := int(pos) - p.base
	end := int(pend) - p.base
	lines := []int{}
	for i := p.linesIdx; i < len(p.lines); {
		// if p.lines[i] < start {
		// 	continue
		// }
		if p.lines[i] > end {
			break
		}
		lines = append(lines, int(p.lines[i]+p.base))
		i++
		p.linesIdx = i
	}
	return lines
}

func (p *parserFromAST) parsePaddingComments(pos token.Pos) {
	var comments []*ast.CommentGroup
	for i := p.commentIdx; i < len(p.f.Comments); {
		c := p.f.Comments[i]
		if pos <= c.Pos() {
			break
		}
		comments = append(comments, c)
		i++
		p.commentIdx = i
	}
	if len(comments) > 0 {
		pos := comments[0].Pos()
		end := comments[len(comments)-1].End()
		region := &Region{
			f:        p.f,
			Pos:      pos,
			End:      end,
			Ref:      &CommentRef{},
			Comments: comments,
			Lines:    p.collectLines(pos, end+1),
		}
		p.regions = append(p.regions, region)
	}
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

	origin := decl.Pos()

	switch x := decl.(type) {
	case *ast.GenDecl:
		if x.Doc != nil {
			origin = x.Doc.Pos()
		}
		p.parsePaddingComments(origin)

		headEnd := x.Specs[0].Pos() // panic?
		// xxx :
		switch x := x.Specs[0].(type) {
		case *ast.ImportSpec:
			if x.Doc != nil {
				headEnd = x.Doc.Pos()
			}
		case *ast.ValueSpec:
			if x.Doc != nil {
				headEnd = x.Doc.Pos()
			}
		case *ast.TypeSpec:
			if x.Doc != nil {
				headEnd = x.Doc.Pos()
			}
		default:
		}
		p.regions = append(p.regions, &Region{
			f:        p.f,
			Pos:      origin,
			End:      headEnd,
			Ref:      &DeclHeadRef{decl: x},
			Lines:    p.collectLines(origin, headEnd),
			Comments: p.collectComments(headEnd),
		})
		p.parseSpecs(x.Specs)

		tailPos := x.End()
		if x.Rparen != token.NoPos {
			tailPos = x.Rparen
		}
		end := x.End()
		p.regions = append(p.regions, &Region{
			f:        p.f,
			Pos:      tailPos,
			End:      end,
			Ref:      &DeclTailRef{decl: x},
			Lines:    p.collectLines(tailPos, end),
			Comments: p.collectComments(end),
		})

	case *ast.FuncDecl:
		if x.Doc != nil {
			origin = x.Doc.Pos()
		}
		p.parsePaddingComments(origin)
		end := x.End()
		p.regions = append(p.regions, &Region{
			f:        p.f,
			Pos:      origin,
			End:      end,
			Ref:      &DeclRef{Decl: x, name: x.Name.Name},
			Lines:    p.collectLines(origin, end),
			Comments: p.collectComments(end),
		})

	case *ast.BadDecl:
		p.parsePaddingComments(origin)
		end := x.End()
		p.regions = append(p.regions, &Region{
			f:        p.f,
			Pos:      origin,
			End:      end,
			Ref:      &DeclRef{Decl: x},
			Lines:    p.collectLines(origin, end),
			Comments: p.collectComments(end),
		})

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
	origin := spec.Pos()

	switch x := spec.(type) {
	case *ast.ImportSpec:
		if x.Doc != nil {
			origin = x.Doc.Pos()
		}
		p.parsePaddingComments(origin)

		end := x.End()
		if x.Comment != nil {
			cend := x.Comment.End()
			if end < cend {
				end = cend
			}
		}
		p.regions = append(p.regions, &Region{
			f:        p.f,
			Pos:      origin,
			End:      end,
			Ref:      &SpecRef{Spec: x},
			Lines:    p.collectLines(origin, end),
			Comments: p.collectComments(end),
		})

	case *ast.ValueSpec:
		if x.Doc != nil {
			origin = x.Doc.Pos()
		}
		p.parsePaddingComments(origin)

		end := x.End()
		if x.Comment != nil {
			cend := x.Comment.End()
			if end < cend {
				end = cend
			}
		}
		p.regions = append(p.regions, &Region{
			f:        p.f,
			Pos:      origin,
			End:      end,
			Ref:      &SpecRef{Spec: x},
			Lines:    p.collectLines(origin, end),
			Comments: p.collectComments(end),
		})

	case *ast.TypeSpec:
		if x.Doc != nil {
			origin = x.Doc.Pos()
		}
		p.parsePaddingComments(origin)

		end := x.End()
		if x.Comment != nil {
			cend := x.Comment.End()
			if end < cend {
				end = cend
			}
		}
		p.regions = append(p.regions, &Region{
			f:        p.f,
			Pos:      origin,
			End:      end,
			Ref:      &SpecRef{Spec: x, name: x.Name.Name},
			Lines:    p.collectLines(origin, end),
			Comments: p.collectComments(end),
		})

	default:
		panic(fmt.Sprintf("invalid spec %q", x))
	}
}
