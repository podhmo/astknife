package patchwork5

import (
	"fmt"
	"go/ast"
	"go/token"
)

func parseASTFile(f *ast.File) *File {
	p := &parserFromAST{f: f}
	p.parseDecls(f.Decls)
	end := f.End()
	if len(f.Comments) > 0 {
		cend := f.Comments[len(f.Comments)-1].End()
		if end < cend {
			end = cend
		}
	}
	p.parseComments(end)
	return &File{
		Regions: p.regions,
	}
}

type parserFromAST struct {
	f       *ast.File
	regions []*Region

	currentPos int
	commentIdx int
}

func (p *parserFromAST) parseComments(pos token.Pos) {
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
		p.regions = append(p.regions, &Region{
			f:      p.f,
			Origin: int(comments[0].Pos()),
			Ref:    &CommentRef{Comments: comments},
		})
	}
}

func (p *parserFromAST) dropComments(pos token.Pos) {
	for i := p.commentIdx; i < len(p.f.Comments); i++ {
		c := p.f.Comments[i]
		if pos <= c.Pos() {
			if pos == c.Pos() {
				p.commentIdx++
			}
			break
		}
		p.commentIdx = i
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
		p.parseComments(origin)
		p.regions = append(p.regions, &Region{
			Origin: int(origin),
			Ref:    &DeclHeadRef{decl: x},
		})
		p.parseSpecs(x.Specs)
		p.regions = append(p.regions, &Region{Ref: &DeclTailRef{decl: x}})
		p.dropComments(x.End())

	case *ast.FuncDecl:
		if x.Doc != nil {
			origin = x.Doc.Pos()
		}
		p.parseComments(origin)
		p.regions = append(p.regions, &Region{
			Origin: int(origin),
			Ref:    &DeclRef{Decl: x, name: x.Name.Name},
		})
		p.dropComments(x.End())

	case *ast.BadDecl:
		p.parseComments(origin)
		p.regions = append(p.regions, &Region{
			Origin: int(origin),
			Ref:    &DeclRef{Decl: x},
		})
		p.dropComments(x.End())

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
		p.parseComments(origin)

		p.regions = append(p.regions, &Region{
			Origin: int(origin),
			Ref:    &SpecRef{Spec: x},
		})

		end := x.End()
		if x.Comment != nil {
			cend := x.Comment.End()
			if end < cend {
				end = cend
			}
		}
		p.dropComments(end)

	case *ast.ValueSpec:
		if x.Doc != nil {
			origin = x.Doc.Pos()
		}
		p.parseComments(origin)

		p.regions = append(p.regions, &Region{
			Origin: int(origin),
			Ref:    &SpecRef{Spec: x},
		})

		end := x.End()
		if x.Comment != nil {
			cend := x.Comment.End()
			if end < cend {
				end = cend
			}
		}
		p.dropComments(end)

	case *ast.TypeSpec:
		if x.Doc != nil {
			origin = x.Doc.Pos()
		}
		p.parseComments(origin)

		p.regions = append(p.regions, &Region{
			Origin: int(origin),
			Ref:    &SpecRef{Spec: x, name: x.Name.Name},
		})

		end := x.End()
		if x.Comment != nil {
			cend := x.Comment.End()
			if end < cend {
				end = cend
			}
		}
		p.dropComments(end)

	default:
		panic(fmt.Sprintf("invalid spec %q", x))
	}
}

func dumpRegions(regions []*Region) {
	for i, r := range regions {
		fmt.Println(i, r, r.Origin)
	}
}
