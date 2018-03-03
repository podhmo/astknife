package mirror

import (
	"go/ast"
	"go/token"

	"github.com/podhmo/astknife/patchwork4/lookup"
)

// State :
type State struct {
	Replacing   map[ast.Node]*lookup.Result
	Appending   []*lookup.Result
	RegionStack []*Region
	File        *ast.File // source file
	Base        int
	Option      Option
}

// Offset :
func (s *State) Offset() int {
	return int(s.RegionStack[len(s.RegionStack)-1].Offset)
}

// Option :
type Option struct {
}

// Region :
type Region struct {
	Offset int
	Pos    token.Pos
	End    token.Pos
}

// NewRegion :
func NewRegion(f *ast.File, base int) *Region {
	offset := int(-f.Pos()) + base
	return &Region{Offset: offset, Pos: token.Pos(base)}
}

// IsFixed :
func (r *Region) IsFixed() bool {
	return r.End != token.NoPos
}

// StartRegion :
func (s *State) StartRegion(src ast.Node, doc *ast.CommentGroup) {
	offset := int(-src.Pos()) + s.Base
	if doc != nil {
		offset += int(src.Pos() - doc.Pos())
	}
	r := &Region{Offset: offset, Pos: token.Pos(s.Base)}
	s.RegionStack = append(s.RegionStack, r)
}

// EndRegion :
func (s *State) EndRegion(dst ast.Node, comment *ast.CommentGroup) {
	end := dst.End()
	r := s.RegionStack[len(s.RegionStack)-1]
	if comment != nil {
		end += comment.End() - dst.End()
	}
	r.End = end
	s.Base = int(r.End)
	s.RegionStack = s.RegionStack[:len(s.RegionStack)-1]
}
