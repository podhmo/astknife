package lookup

import (
	"go/ast"
)

// isMethod :
func isMethod(fn *ast.FuncDecl) bool {
	return fn.Recv != nil
}

// isSameTypeOrPointer :
func isSameTypeOrPointer(ob *ast.Object, fn ast.Node) bool {
	switch t := fn.(type) {
	case *ast.StarExpr:
		return isSameTypeOrPointer(ob, t.X)
	case *ast.Ident:
		return t.Obj.Name == ob.Name
	default:
		return false
	}
}
