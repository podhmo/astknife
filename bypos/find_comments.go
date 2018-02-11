package bypos

import (
	"go/ast"
	"go/token"
	"log"

	"golang.org/x/tools/go/ast/astutil"
)

// FindComments :
func FindComments(files []*ast.File, pos token.Pos) *ast.CommentGroup {
	file := FindFile(files, pos)
	if file == nil {
		return nil
	}

	nodes, _ := astutil.PathEnclosingInterval(file, pos, pos)
	if len(nodes) <= 0 {
		return nil
	}

	switch t := nodes[0].(type) {
	case *ast.GenDecl:
		return t.Doc
	case *ast.FuncDecl:
		return t.Doc
	case *ast.Ident:
		if t.Obj == nil {
			return nil
		}
		switch x := t.Obj.Decl.(type) {
		case *ast.Field:
			if x.Doc != nil {
				return x.Doc
			}
			return x.Comment
		case *ast.ImportSpec:
			if x.Doc != nil {
				return x.Doc
			}
			return x.Comment
		case *ast.ValueSpec:
			if x.Doc != nil {
				return x.Doc
			}
			return x.Comment
		case *ast.TypeSpec:
			if x.Doc != nil {
				return x.Doc
			}
			return x.Comment
		default:
			log.Printf("default2: %#v\n", x)
			return nil
		}
	default:
		log.Printf("default: %#v\n", t)
		return nil
	}
}
