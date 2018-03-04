package buildtree

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"

	"github.com/podhmo/astknife/patchwork4/debug"
	"github.com/podhmo/astknife/patchwork4/lookup"
)

// State :
type State struct {
	Fset  *token.FileSet
	File  *ast.File // source file
	Debug *debug.Debug

	Replacing map[ast.Node]*lookup.Result
	Appending []*lookup.Result

	RegionStack []*Region
	Lines       map[int]int

	FileBase int
	Base     int
}

// Offset :
func (s *State) Offset() int {
	return int(s.RegionStack[len(s.RegionStack)-1].Offset)
}

// Region :
type Region struct {
	Pos token.Pos
	End token.Pos

	pat       ast.Node
	rep       ast.Node
	sourcePos token.Pos
	sourceEnd token.Pos
	Offset    int // for another *ast.File
	Delta     int // for comment area
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
	var sourcePos token.Pos
	if doc != nil {
		sourcePos = doc.Pos()
		delta += int(rep.Pos() - doc.Pos())
	} else {
		sourcePos = rep.Pos()
	}

	offset := int(-rep.Pos()) + base + delta
	// fmt.Printf("** start region (base=%d, offset=%d, delta=%d, comment=%v)\n", s.Base, offset, delta, doc != nil)
	r := &Region{Offset: offset, Pos: token.Pos(base + delta), pat: pat, rep: rep, sourcePos: sourcePos, Delta: delta}
	s.RegionStack = append(s.RegionStack, r)
}

// EndRegion :
func (s *State) EndRegion(dst ast.Node, comment *ast.CommentGroup) {
	end := dst.End()
	r := s.RegionStack[len(s.RegionStack)-1]
	if comment != nil {
		r.sourceEnd = r.rep.End() + (comment.End() - dst.End())
		end += comment.End() - dst.End()
	} else {
		r.sourceEnd = r.rep.End()
	}

	r.End = end
	fmt.Printf("region fixed (%d, %d). (base=%d, original=(%s, %s))\n", r.Pos, r.End, s.Base, s.Fset.Position(r.sourcePos), s.Fset.Position(r.sourceEnd))

	if s.Debug != nil {
		fmt.Println("----------------------------------------")
		f := s.Fset.File(r.sourcePos)
		pos := s.Fset.Position(r.sourcePos)
		source := s.Debug.SourceMap[pos.Filename]
		fmt.Println(f.Name(), source[f.Offset(r.sourcePos):f.Offset(r.sourceEnd)])
		fmt.Println("----------------------------------------")
		{
			// todo: cache
			// f := s.Fset.File(r.sourcePos)
			lines := reflect.ValueOf(f).Elem().FieldByName("lines")
			for i := 0; i < lines.Len(); i++ {
				p := f.Base() + int(lines.Index(i).Int())
				if int(r.sourceEnd) < p {
					break
				}
				if int(r.sourcePos) <= p {
					s.Lines[p] = r.Offset + p - s.FileBase
					fmt.Println("	", p, "->", r.Offset+p, "(", f.Name(), s.Lines[p], ")")
				}
			}
		}
	}

	s.Base = int(r.End)
	// fmt.Printf("** end region (base=%d, offset=%d, comment=%v)\n", s.Base, r.Offset, comment != nil)
	s.RegionStack = s.RegionStack[:len(s.RegionStack)-1]
}
