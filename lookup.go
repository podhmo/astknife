package astknife

import (
	"go/ast"
)

// LookupAllMethods :
func LookupAllMethods(f *ast.File, obname string) []*ast.FuncDecl {
	ob := f.Scope.Lookup(obname)
	if ob == nil {
		return nil
	}

	var r []*ast.FuncDecl
	for _, decl := range f.Decls {
		if decl, ok := decl.(*ast.FuncDecl); ok {
			if IsMethod(decl) && IsSameTypeOrPointer(ob, decl.Recv.List[0].Type) {
				r = append(r, decl)
			}
		}
	}
	return r
}

// LookupMethod :
func LookupMethod(f *ast.File, obname string, name string) *ast.FuncDecl {
	ob := f.Scope.Lookup(obname)
	if ob == nil {
		return nil
	}

	for _, decl := range f.Decls {
		if decl, ok := decl.(*ast.FuncDecl); ok {
			if IsMethod(decl) && IsSameTypeOrPointer(ob, decl.Recv.List[0].Type) {
				if decl.Name.Name == name {
					return decl
				}
			}
		}
	}
	return nil
}
