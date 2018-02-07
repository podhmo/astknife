package astknife

import (
	"go/ast"
)

// IsMethod :
func IsMethod(fn *ast.FuncDecl) bool {
	return fn.Recv != nil
}

// IsSameTypeOrPointer :
func IsSameTypeOrPointer(ob *ast.Object, fn ast.Node) bool {
	switch t := fn.(type) {
	case *ast.StarExpr:
		return IsSameTypeOrPointer(ob, t.X)
	case *ast.Ident:
		return t.Obj.Name == ob.Name
	default:
		return false
	}
}
