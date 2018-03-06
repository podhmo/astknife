package buildtree

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/podhmo/astknife/action/comment"
	"github.com/podhmo/astknife/patchwork4/debug"
	"github.com/podhmo/astknife/patchwork4/lookup"
)

// State :
type State struct {
	Fset  *token.FileSet
	Debug *debug.Debug

	Replacing map[ast.Node]*lookup.Result
	Appending []*lookup.Result

	RegionStack  []*Region
	LatestRegion *Region
	Lines        map[int]int

	FileBase int
	Base     int
}

// Offset :
func (s *State) Offset() int {
	return int(s.RegionStack[len(s.RegionStack)-1].Offset)
}

// File :
func (s *State) File() *ast.File {
	return s.RegionStack[len(s.RegionStack)-1].sourceFile
}

// Region :
type Region struct {
	Pos token.Pos
	End token.Pos

	pat        ast.Node
	rep        ast.Node
	sourcePos  token.Pos
	sourceEnd  token.Pos
	sourceFile *ast.File
	Offset     int // for another *ast.File
	Delta      int // for comment area
}

// NewRegion :
func NewRegion(f *ast.File, base int) *Region {
	offset := int(-f.Pos()) + base
	return &Region{
		Offset:     offset,
		Pos:        token.Pos(base),
		pat:        f,
		sourceFile: f,
	}
}

// IsFixed :
func (r *Region) IsFixed() bool {
	return r.End != token.NoPos
}

// StartRegion :
func (s *State) StartRegion(pat, rep ast.Node, doc *ast.CommentGroup, file *ast.File) {
	base := s.Base
	delta := 0
	if len(s.RegionStack) > 1 {
		parentRegion := s.RegionStack[len(s.RegionStack)-1]
		if !parentRegion.IsFixed() {
			delta += parentRegion.Delta
		}
	}
	var sourcePos token.Pos
	{
		if s.LatestRegion.sourceFile == s.RegionStack[0].sourceFile {
			for _, cg := range comment.Collect(file.Comments, s.LatestRegion.pat.End(), pat.Pos()) {
				fmt.Printf("%s	!! %q\n", strings.Repeat("  ", len(s.RegionStack)), cg.Text())
			}
		}
	}
	if doc != nil {
		sourcePos = doc.Pos()
		delta += int(rep.Pos() - doc.Pos())
	} else {
		sourcePos = rep.Pos()
	}

	offset := int(-rep.Pos()) + base + delta
	fmt.Printf("%sstart region (base=%d, offset=%d, delta=%d, comment=%v)\n", strings.Repeat("  ", len(s.RegionStack)), s.Base, offset, delta, doc != nil)
	r := &Region{
		Offset:     offset,
		Pos:        token.Pos(base + delta),
		pat:        pat,
		rep:        rep,
		sourceFile: file,
		sourcePos:  sourcePos,
		Delta:      delta,
	}
	s.RegionStack = append(s.RegionStack, r)
}

// EndRegion :
func (s *State) EndRegion(dst ast.Node, lastComment *ast.CommentGroup) {
	end := dst.End()
	r := s.RegionStack[len(s.RegionStack)-1]
	if lastComment != nil {
		r.sourceEnd = r.rep.End() + (lastComment.End() - dst.End())
		end += lastComment.End() - dst.End()
	} else {
		r.sourceEnd = r.rep.End()
	}

	r.End = end
	fmt.Printf("%sregion fixed (%d, %d). (base=%d, original=(%s, %s))\n", strings.Repeat("  ", len(s.RegionStack)), r.Pos, r.End, s.Base, s.Fset.Position(r.sourcePos), s.Fset.Position(r.sourceEnd))

	if s.Debug != nil {
		for _, cg := range comment.Collect(r.sourceFile.Comments, r.sourcePos, r.sourceEnd) {
			fmt.Printf("%s	** %q\n", strings.Repeat("  ", len(s.RegionStack)), cg.Text())
		}
	}

	s.Base = int(r.End)
	s.RegionStack = s.RegionStack[:len(s.RegionStack)-1]
	s.LatestRegion = r
	fmt.Printf("%send region (base=%d, offset=%d, lastComment=%v)\n", strings.Repeat("  ", len(s.RegionStack)), s.Base, r.Offset, lastComment != nil)
}
