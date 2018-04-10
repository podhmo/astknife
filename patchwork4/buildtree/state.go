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

	FileBase          int
	Base              int
	CollectedComments map[token.Pos]*ast.CommentGroup
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
	Comments   []*ast.CommentGroup
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
	var comments []*ast.CommentGroup
	{
		if s.LatestRegion.sourceFile == file {
			for _, cg := range comment.Collect(file.Comments, s.LatestRegion.sourceEnd, pat.Pos()) {
				// fmt.Printf("%s	!! %q(%d)\n", strings.Repeat("  ", len(s.RegionStack)), cg.Text(), delta)
				if _, ok := s.CollectedComments[cg.Pos()]; !ok {
					comments = append(comments, cg)
					// todo: + delta?
				}
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
	fmt.Printf("%sstart region (base=%d, offset=%d, delta=%d, comment=%v, file=%s)\n", strings.Repeat("  ", len(s.RegionStack)), s.Base, offset, delta, doc != nil, s.Fset.File(file.Pos()).Name())
	r := &Region{
		Offset:     offset,
		Pos:        token.Pos(base + delta),
		pat:        pat,
		rep:        rep,
		sourceFile: file,
		sourcePos:  sourcePos,
		Delta:      delta,
		Comments:   comments,
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
	fmt.Printf("%sregion fixed %T (%d, %d). (base=%d, original=(%s, %s))\n", strings.Repeat("  ", len(s.RegionStack)), r.pat, r.Pos, r.End, s.Base, s.Fset.Position(r.sourcePos), s.Fset.Position(r.sourceEnd))

	for _, cg := range comment.Collect(r.sourceFile.Comments, r.sourcePos, r.pat.Pos()) {
		r.Comments = append(r.Comments, cg)
	}
	for _, cg := range comment.Collect(r.sourceFile.Comments, r.pat.End(), r.sourceEnd) {
		r.Comments = append(r.Comments, cg)
	}
	for _, cg := range r.Comments {
		if _, ok := s.CollectedComments[cg.Pos()]; !ok {
			fmt.Printf("%s	(%d) ** %q\n", strings.Repeat("  ", len(s.RegionStack)), cg.Pos(), cg.Text())
			s.CollectedComments[cg.Pos()] = CommentGroup(cg, s)
		}
	}

	s.Base = int(r.End)
	s.RegionStack = s.RegionStack[:len(s.RegionStack)-1]
	s.LatestRegion = r
	fmt.Printf("%send region (base=%d, offset=%d, lastComment=%v, file=%ss)\n", strings.Repeat("  ", len(s.RegionStack)), s.Base, r.Offset, lastComment != nil, s.Fset.File(r.sourceFile.Pos()).Name())
}
