package comment

import (
	"go/ast"
	"go/token"
)

// Collect :
func Collect(comments []*ast.CommentGroup, pos, end token.Pos) []*ast.CommentGroup {
	r := make([]*ast.CommentGroup, 0, len(comments))

	for i, c := range comments {
		if c.Pos() >= pos {
			for ; i < len(comments); i++ {
				if comments[i].Pos() >= end {
					return r
				}
				r = append(r, comments[i])
			}
		}
	}
	return r
}

// CollectFromNode :
func CollectFromNode(comments []*ast.CommentGroup, node ast.Node) []*ast.CommentGroup {
	// todo: fix
	switch t := node.(type) {
	case *ast.GenDecl:
		if t.Doc != nil {
			return Collect(comments, t.Doc.Pos(), node.End())
		}
		return Collect(comments, node.Pos()-1, node.End())
	case *ast.FuncDecl:
		if t.Doc != nil {
			return Collect(comments, t.Doc.Pos(), node.End())
		}
		return Collect(comments, node.Pos()-1, node.End())
	default:
		return Collect(comments, node.Pos(), node.End())
	}
}
