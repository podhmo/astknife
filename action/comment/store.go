package comment

import (
	"fmt"
	"go/ast"
	"go/token"
)

// Store :
type Store struct {
	Fset *token.FileSet
	F    *ast.File
	Refs []*Ref
}

// New :
func New(fset *token.FileSet, f *ast.File) *Store {
	refs := make([]*Ref, len(f.Comments))
	for i, c := range f.Comments {
		refs[i] = &Ref{FromOther: false, Comment: c}
	}
	return &Store{
		Fset: fset,
		F:    f,
		Refs: refs,
	}
}

// Ref :
type Ref struct {
	Comment   *ast.CommentGroup
	FromOther bool // todo: not boolean
}

// ReplaceFromOther :
func (s *Store) ReplaceFromOther(node ast.Node, replacement ast.Node, comments []*ast.CommentGroup) {
	pos := node.Pos()
	end := node.End()
	new := make([]*Ref, 0, len(s.Refs)+len(comments))
	fmt.Println(len(comments))
	for i, ref := range s.Refs {
		fmt.Printf("%q %d %d %v\n", ref.Comment.Text(), ref.Comment.Pos(), pos, ref.Comment.Pos() >= pos)
		if ref.FromOther {
			new = append(new, ref)
			continue
		}
		if ref.Comment.Pos() >= pos {
			if len(new) > 0 && !new[len(new)-1].FromOther {
				// pop latest
				// latest := new[len(new)-1]
				new = new[:len(new)-1]

				hackstartPos := token.Pos(int(node.Pos() - 1))
				// if latest.Comment.Pos() > hackstartPos {
				// 	hackstartPos = token.Pos(int(latest.Comment.Pos() - 1))
				// }

				new = append(new, &Ref{
					Comment: &ast.CommentGroup{
						List: []*ast.Comment{{Slash: hackstartPos, Text: "// *hack* start"}},
					},
					FromOther: true,
				})
			}

			for _, c := range comments {
				new = append(new, &Ref{Comment: c, FromOther: true})
			}
			new = append(new, &Ref{
				Comment: &ast.CommentGroup{
					List: []*ast.Comment{{Slash: token.Pos(int(replacement.End() - 1)), Text: "// *hack* end"}},
				},
				FromOther: true,
			})
			for ; i < len(s.Refs); i++ {
				ref := s.Refs[i]
				if ref.FromOther || end < ref.Comment.Pos() {
					break
				}
			}
			new = append(new, s.Refs[i:]...)
			break
		}
		new = append(new, ref)
	}
	s.Refs = new
}

// Comments :
func (s *Store) Comments() []*ast.CommentGroup {
	r := make([]*ast.CommentGroup, len(s.Refs))
	for i := range s.Refs {
		r[i] = s.Refs[i].Comment
	}
	return r
}
