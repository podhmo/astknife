package patchwork5

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
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

func (p *parserFromAST) parsePaddingComments(pos token.Pos) {
	var comments []*ast.CommentGroup
	for i := p.commentIdx; i < len(p.f.Comments); i++ {
		c := p.f.Comments[i]
		if pos <= c.Pos() {
			if pos == c.Pos() {
				p.commentIdx++
			}
			break
		}
		comments = append(comments, c)
		p.commentIdx = i
	}
	if len(comments) > 0 {
		pos := comments[0].Pos()
		region := &Region{
			f:        p.f,
			Origin:   int(pos),
			Ref:      &CommentRef{},
			Comments: comments,
			Lines:    p.calcLines(pos, comments[len(comments)-1].End()),
		}
		p.regions = append(p.regions, region)
	}
}

func (p *parserFromAST) collectComments(pos token.Pos) []*ast.CommentGroup {
	var comments []*ast.CommentGroup
	for i := p.commentIdx; i < len(p.f.Comments); i++ {
		c := p.f.Comments[i]
		if pos < c.Pos() {
			break
		}
		comments = append(comments, c)
		p.commentIdx = i
	}
	return comments
}

// calcLines returns offset list of each region's origin
func (p *parserFromAST) calcLines(pos, pend token.Pos) []int {
	start := int(pos) - p.base
	end := int(pend) - p.base
	var lines []int
	for i := p.linesIdx; i < len(p.lines); i++ {
		if p.lines[i] < start {
			continue
		}
		lines = append(lines, (p.lines[i] - start))
		if p.lines[i] > end {
			break
		}
	}
	return lines
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
		p.regions = append(p.regions, &Region{
			f:        p.f,
			Origin:   int(origin),
			Ref:      &DeclHeadRef{decl: x},
			Lines:    p.calcLines(origin, headEnd),
			Comments: p.collectComments(headEnd),
		})
		p.parseSpecs(x.Specs)

		tailOrigin := x.End()
		if x.Rparen != token.NoPos {
			tailOrigin = x.Rparen
		}
		end := x.End()
		p.regions = append(p.regions, &Region{
			f:        p.f,
			Origin:   int(tailOrigin),
			Ref:      &DeclTailRef{decl: x},
			Lines:    p.calcLines(tailOrigin, end),
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
			Origin:   int(origin),
			Ref:      &DeclRef{Decl: x, name: x.Name.Name},
			Lines:    p.calcLines(origin, end),
			Comments: p.collectComments(end),
		})

	case *ast.BadDecl:
		p.parsePaddingComments(origin)
		end := x.End()
		p.regions = append(p.regions, &Region{
			f:        p.f,
			Origin:   int(origin),
			Ref:      &DeclRef{Decl: x},
			Lines:    p.calcLines(origin, end),
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
			Origin:   int(origin),
			Ref:      &SpecRef{Spec: x},
			Lines:    p.calcLines(origin, end),
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
			Origin:   int(origin),
			Ref:      &SpecRef{Spec: x},
			Lines:    p.calcLines(origin, end),
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
			Origin:   int(origin),
			Ref:      &SpecRef{Spec: x, name: x.Name.Name},
			Lines:    p.calcLines(origin, end),
			Comments: p.collectComments(end),
		})

	default:
		panic(fmt.Sprintf("invalid spec %q", x))
	}
}

func dumpRegions(regions []*Region) {
	for i, r := range regions {
		fmt.Println(i, r, "origin", r.Origin, "lines", r.Lines)
	}
}
