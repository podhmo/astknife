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
	Pos token.Pos
	End token.Pos

	Ob     ast.Node
	Offset int // for another *ast.File
	Delta  int // for comment area
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
func (s *State) StartRegion(pat, rep ast.Node, doc *ast.CommentGroup) {
	base := s.Base
	delta := 0
	if len(s.RegionStack) > 1 {
		parentRegion := s.RegionStack[len(s.RegionStack)-1]
		if !parentRegion.IsFixed() {
			delta += parentRegion.Delta
		}
	}
	if doc != nil {
		delta += int(rep.Pos() - doc.Pos())
	}
	offset := int(-rep.Pos()) + base + delta
	// fmt.Printf("** start region (base=%d, offset=%d, delta=%d, comment=%v)\n", s.Base, offset, delta, doc != nil)
	r := &Region{Offset: offset, Pos: token.Pos(base + delta), Ob: pat, Delta: delta}
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
	// fmt.Printf("** end region (base=%d, offset=%d, comment=%v)\n", s.Base, r.Offset, comment != nil)
	s.RegionStack = s.RegionStack[:len(s.RegionStack)-1]
}
