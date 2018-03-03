package mirror

import (
	"fmt"
	"go/ast"
	"go/token"
)

// File :
func File(f *ast.File, s *State) *ast.File {
	return &ast.File{
		Doc:        CommentGroup(f.Doc, s),
		Package:    token.Pos(int(f.Package) + s.Offset()),
		Name:       Ident(f.Name, s),
		Decls:      Decls(f.Decls, s),
		Scope:      ast.NewScope(nil), // todo:
		Imports:    ImportSpecs(f.Imports, s),
		Unresolved: Idents(f.Unresolved, s),
		Comments:   CommentGroups(f.Comments, s),
	}
}

// Decls :
func Decls(xs []ast.Decl, s *State) []ast.Decl {
	ys := make([]ast.Decl, len(xs))
	for i := range xs {
		ys[i] = Decl(xs[i], s)
	}
	return ys
}

// Decl :
func Decl(decl ast.Decl, s *State) ast.Decl {
	if decl == nil {
		return nil
	}
	original := decl
	if rep, ok := s.Replacing[decl]; ok {
		decl = rep.Object.Decl.(ast.Decl)
	}
	switch x := decl.(type) {
	case *ast.GenDecl:
		s.StartRegion(original, x, x.Doc)
		new := &ast.GenDecl{
			Doc:    CommentGroup(x.Doc, s),
			TokPos: token.Pos(int(x.TokPos) + s.Offset()),
			Tok:    x.Tok,
			Lparen: x.Lparen,
			Specs:  Specs(x.Specs, s),
			Rparen: x.Rparen,
		}
		if new.Lparen != token.NoPos {
			new.Lparen = token.Pos(int(x.Lparen) + s.Offset())
		}
		if new.Rparen != token.NoPos {
			new.Rparen = token.Pos(int(x.Rparen) + s.Offset())
		}
		s.EndRegion(new, nil)
		return new
	case *ast.FuncDecl:
		s.StartRegion(original, x, x.Doc)
		new := &ast.FuncDecl{
			Doc:  CommentGroup(x.Doc, s),
			Recv: FieldList(x.Recv, s),
			Name: Ident(x.Name, s),
			Type: FuncType(x.Type, s),
			Body: BlockStmt(x.Body, s),
		}
		s.EndRegion(new, nil)
		return new
	case *ast.BadDecl:
		s.StartRegion(original, x, nil)
		new := &ast.BadDecl{
			From: token.Pos(int(x.From) + s.Offset()),
			To:   token.Pos(int(x.To) + s.Offset()),
		}
		s.EndRegion(new, nil)
		return new
	default:
		panic(fmt.Sprintf("invalid decl %q", x))
	}
}

// Specs :
func Specs(xs []ast.Spec, s *State) []ast.Spec {
	ys := make([]ast.Spec, len(xs))
	for i := range xs {
		ys[i] = Spec(xs[i], s)
	}
	return ys
}

// ImportSpecs :
func ImportSpecs(xs []*ast.ImportSpec, s *State) []*ast.ImportSpec {
	ys := make([]*ast.ImportSpec, len(xs))
	for i := range xs {
		ys[i] = Spec(xs[i], s).(*ast.ImportSpec)
	}
	return ys
}

// Spec :
func Spec(spec ast.Spec, s *State) ast.Spec {
	if spec == nil {
		return nil
	}
	original := spec
	if rep, ok := s.Replacing[spec]; ok {
		spec = rep.Object.Decl.(ast.Spec)
	}

	switch x := spec.(type) {
	case *ast.ImportSpec:
		s.StartRegion(original, x, x.Doc)
		endpos := x.EndPos
		if endpos != token.NoPos {
			endpos = token.Pos(int(x.EndPos) + s.Offset())
		}
		new := &ast.ImportSpec{
			Doc:     CommentGroup(x.Doc, s),
			Name:    Ident(x.Name, s),
			Path:    BasicLit(x.Path, s),
			Comment: CommentGroup(x.Comment, s),
			EndPos:  endpos,
		}
		s.EndRegion(new, new.Comment)
		return new

	case *ast.ValueSpec:
		s.StartRegion(original, x, x.Doc)
		new := &ast.ValueSpec{
			Doc:     CommentGroup(x.Doc, s),
			Names:   Idents(x.Names, s),
			Type:    Expr(x.Type, s),
			Values:  Exprs(x.Values, s),
			Comment: CommentGroup(x.Comment, s),
		}
		s.EndRegion(new, new.Comment)
		return new
	case *ast.TypeSpec:
		s.StartRegion(original, x, x.Doc)
		assign := x.Assign
		if assign != token.NoPos {
			assign = token.Pos(int(x.Assign) + s.Offset())
		}
		new := &ast.TypeSpec{
			Doc:     CommentGroup(x.Doc, s),
			Name:    Ident(x.Name, s),
			Assign:  assign,
			Type:    Expr(x.Type, s),
			Comment: CommentGroup(x.Comment, s),
		}
		s.EndRegion(new, new.Comment)
		return new
	default:
		panic(fmt.Sprintf("invalid spec %q", x))
	}
}

// FieldList :
func FieldList(x *ast.FieldList, s *State) *ast.FieldList {
	if x == nil {
		return nil
	}

	opening := x.Opening
	if opening != token.NoPos {
		opening = token.Pos(int(x.Opening) + s.Offset())
	}
	closing := x.Closing
	if closing != token.NoPos {
		closing = token.Pos(int(x.Closing) + s.Offset())
	}
	return &ast.FieldList{
		Opening: opening,
		List:    Fields(x.List, s),
		Closing: closing,
	}
}

// Fields :
func Fields(xs []*ast.Field, s *State) []*ast.Field {
	ys := make([]*ast.Field, len(xs))
	for i := range xs {
		ys[i] = Field(xs[i], s)
	}
	return ys
}

// Field :
func Field(x *ast.Field, s *State) *ast.Field {
	if x == nil {
		return nil
	}
	s.StartRegion(x, x, x.Doc)
	new := &ast.Field{
		Doc:     CommentGroup(x.Doc, s),
		Names:   Idents(x.Names, s),
		Type:    Expr(x.Type, s),
		Tag:     BasicLit(x.Tag, s),
		Comment: CommentGroup(x.Comment, s),
	}
	s.EndRegion(new, new.Comment)
	return new
}

// Idents :
func Idents(xs []*ast.Ident, s *State) []*ast.Ident {
	ys := make([]*ast.Ident, len(xs))
	for i := range xs {
		ys[i] = Ident(xs[i], s)
	}
	return ys
}

// Ident :
func Ident(x *ast.Ident, s *State) *ast.Ident {
	if x == nil {
		return nil
	}
	// fmt.Println("Idnt", x.Name, "NamePos", x.NamePos, "->", token.Pos(int(x.NamePos)+s.Offset()), "offset=", s.Offset())
	return &ast.Ident{
		NamePos: token.Pos(int(x.NamePos) + s.Offset()),
		Name:    x.Name,
		// Object
	}
}

// FuncType :
func FuncType(x *ast.FuncType, s *State) *ast.FuncType {
	if x == nil {
		return nil
	}
	f := x.Func
	if f != token.NoPos {
		f = token.Pos(int(x.Func) + s.Offset())
	}
	new := &ast.FuncType{
		Func:    f,
		Params:  FieldList(x.Params, s),
		Results: FieldList(x.Results, s),
	}
	return new
}

// Stmts :
func Stmts(xs []ast.Stmt, s *State) []ast.Stmt {
	ys := make([]ast.Stmt, len(xs))
	for i := range xs {
		ys[i] = Stmt(xs[i], s)
	}
	return ys
}

// Stmt :
func Stmt(stmt ast.Stmt, s *State) ast.Stmt {
	if stmt == nil {
		return nil
	}
	switch x := stmt.(type) {
	case *ast.BadStmt:
		return &ast.BadStmt{
			From: token.Pos(int(x.From) + s.Offset()),
			To:   token.Pos(int(x.To) + s.Offset()),
		}
	case *ast.DeclStmt:
		return &ast.DeclStmt{
			Decl: Decl(x.Decl, s),
		}
	case *ast.EmptyStmt:
		return &ast.EmptyStmt{
			Semicolon: token.Pos(int(x.Semicolon) + s.Offset()),
			Implicit:  x.Implicit,
		}
	case *ast.LabeledStmt:
		return &ast.LabeledStmt{
			Label: Ident(x.Label, s),
			Colon: token.Pos(int(x.Colon) + s.Offset()),
			Stmt:  Stmt(x.Stmt, s),
		}
	case *ast.ExprStmt:
		return &ast.ExprStmt{
			X: Expr(x.X, s),
		}
	case *ast.SendStmt:
		return &ast.SendStmt{
			Chan:  Expr(x.Chan, s),
			Arrow: token.Pos(int(x.Arrow) + s.Offset()),
			Value: Expr(x.Value, s),
		}
	case *ast.IncDecStmt:
		return &ast.IncDecStmt{
			X:      Expr(x.X, s),
			TokPos: token.Pos(int(x.TokPos) + s.Offset()),
			Tok:    x.Tok,
		}
	case *ast.AssignStmt:
		return &ast.AssignStmt{
			Lhs:    Exprs(x.Lhs, s),
			TokPos: token.Pos(int(x.TokPos) + s.Offset()),
			Tok:    x.Tok,
			Rhs:    Exprs(x.Rhs, s),
		}
	case *ast.GoStmt:
		return &ast.GoStmt{
			Go:   token.Pos(int(x.Go) + s.Offset()),
			Call: CallExpr(x.Call, s),
		}
	case *ast.DeferStmt:
		return &ast.DeferStmt{
			Defer: token.Pos(int(x.Defer) + s.Offset()),
			Call:  CallExpr(x.Call, s),
		}
	case *ast.ReturnStmt:
		return &ast.ReturnStmt{
			Return:  token.Pos(int(x.Return) + s.Offset()),
			Results: Exprs(x.Results, s),
		}
	case *ast.BranchStmt:
		return &ast.BranchStmt{
			TokPos: token.Pos(int(x.TokPos) + s.Offset()),
			Tok:    x.Tok,
			Label:  Ident(x.Label, s),
		}
	case *ast.BlockStmt:
		return BlockStmt(x, s)

	case *ast.IfStmt:
		return &ast.IfStmt{
			If:   token.Pos(int(x.If) + s.Offset()),
			Init: Stmt(x.Init, s),
			Cond: Expr(x.Cond, s),
			Body: BlockStmt(x.Body, s),
			Else: Stmt(x.Else, s),
		}
	case *ast.CaseClause:
		return &ast.CaseClause{
			Case:  token.Pos(int(x.Case) + s.Offset()),
			List:  Exprs(x.List, s),
			Colon: token.Pos(int(x.Colon) + s.Offset()),
			Body:  Stmts(x.Body, s),
		}
	case *ast.SwitchStmt:
		return &ast.SwitchStmt{
			Switch: token.Pos(int(x.Switch) + s.Offset()),
			Init:   Stmt(x.Init, s),
			Tag:    Expr(x.Tag, s),
			Body:   BlockStmt(x.Body, s),
		}
	case *ast.TypeSwitchStmt:
		return &ast.TypeSwitchStmt{
			Switch: token.Pos(int(x.Switch) + s.Offset()),
			Init:   Stmt(x.Init, s),
			Assign: Stmt(x.Assign, s),
			Body:   BlockStmt(x.Body, s),
		}
	case *ast.CommClause:
		return &ast.CommClause{
			Case:  token.Pos(int(x.Case) + s.Offset()),
			Comm:  Stmt(x.Comm, s),
			Colon: token.Pos(int(x.Colon) + s.Offset()),
			Body:  Stmts(x.Body, s),
		}
	case *ast.SelectStmt:
		return &ast.SelectStmt{
			Select: token.Pos(int(x.Select) + s.Offset()),
			Body:   BlockStmt(x.Body, s),
		}
	case *ast.ForStmt:
		return &ast.ForStmt{
			For:  token.Pos(int(x.For) + s.Offset()),
			Init: Stmt(x.Init, s),
			Cond: Expr(x.Cond, s),
			Post: Stmt(x.Post, s),
			Body: BlockStmt(x.Body, s),
		}
	case *ast.RangeStmt:
		return &ast.RangeStmt{
			For:    token.Pos(int(x.For) + s.Offset()),
			Key:    Expr(x.Key, s),
			Value:  Expr(x.Value, s),
			TokPos: token.Pos(int(x.TokPos) + s.Offset()),
			Tok:    x.Tok,
			X:      Expr(x.X, s),
			Body:   BlockStmt(x.Body, s),
		}
	default:
		panic(fmt.Sprintf("invalid stmt %q", x))
	}
}

// CallExpr :
func CallExpr(x *ast.CallExpr, s *State) *ast.CallExpr {
	if x == nil {
		return nil
	}
	return &ast.CallExpr{
		Fun:      Expr(x.Fun, s),
		Lparen:   token.Pos(int(x.Lparen) + s.Offset()),
		Args:     Exprs(x.Args, s),
		Ellipsis: token.Pos(int(x.Ellipsis) + s.Offset()),
		Rparen:   token.Pos(int(x.Rparen) + s.Offset()),
	}
}

// Exprs :
func Exprs(xs []ast.Expr, s *State) []ast.Expr {
	ys := make([]ast.Expr, len(xs))
	for i := range xs {
		ys[i] = Expr(xs[i], s)
	}
	return ys
}

// Expr :
func Expr(expr ast.Expr, s *State) ast.Expr {
	if expr == nil {
		return nil
	}
	switch x := expr.(type) {
	case *ast.BadExpr:
		return &ast.BadExpr{
			From: token.Pos(int(x.From) + s.Offset()),
			To:   token.Pos(int(x.To) + s.Offset()),
		}
	case *ast.Ident:
		return Ident(x, s)
	case *ast.Ellipsis:
		return &ast.Ellipsis{
			Ellipsis: token.Pos(int(x.Ellipsis) + s.Offset()),
			Elt:      Expr(x.Elt, s),
		}
	case *ast.BasicLit:
		return BasicLit(x, s)
	case *ast.FuncLit:
		return &ast.FuncLit{
			Type: FuncType(x.Type, s),
			Body: BlockStmt(x.Body, s),
		}
	case *ast.CompositeLit:
		return &ast.CompositeLit{
			Type:   Expr(x.Type, s),
			Lbrace: token.Pos(int(x.Lbrace) + s.Offset()),
			Elts:   Exprs(x.Elts, s),
			Rbrace: token.Pos(int(x.Rbrace) + s.Offset()),
		}
	case *ast.ParenExpr:
		return &ast.ParenExpr{
			Lparen: token.Pos(int(x.Lparen) + s.Offset()),
			X:      Expr(x.X, s),
			Rparen: token.Pos(int(x.Rparen) + s.Offset()),
		}
	case *ast.SelectorExpr:
		return &ast.SelectorExpr{
			X:   Expr(x.X, s),
			Sel: Ident(x.Sel, s),
		}
	case *ast.IndexExpr:
		return &ast.IndexExpr{
			X:      Expr(x.X, s),
			Lbrack: token.Pos(int(x.Lbrack) + s.Offset()),
			Index:  Expr(x.Index, s),
			Rbrack: token.Pos(int(x.Rbrack) + s.Offset()),
		}
	case *ast.SliceExpr:
		return &ast.SliceExpr{
			X:      Expr(x.X, s),
			Lbrack: token.Pos(int(x.Lbrack) + s.Offset()),
			Low:    Expr(x.Low, s),
			High:   Expr(x.High, s),
			Max:    Expr(x.Max, s),
			Slice3: x.Slice3,
			Rbrack: token.Pos(int(x.Rbrack) + s.Offset()),
		}
	case *ast.TypeAssertExpr:
		return &ast.TypeAssertExpr{
			X:      Expr(x.X, s),
			Lparen: token.Pos(int(x.Lparen) + s.Offset()),
			Type:   Expr(x.Type, s),
			Rparen: token.Pos(int(x.Rparen) + s.Offset()),
		}
	case *ast.CallExpr:
		return CallExpr(x, s)
	case *ast.StarExpr:
		return &ast.StarExpr{
			Star: token.Pos(int(x.Star) + s.Offset()),
			X:    Expr(x.X, s),
		}
	case *ast.UnaryExpr:
		return &ast.UnaryExpr{
			OpPos: token.Pos(int(x.OpPos) + s.Offset()),
			Op:    x.Op,
			X:     Expr(x.X, s),
		}
	case *ast.BinaryExpr:
		return &ast.BinaryExpr{
			X:     Expr(x.X, s),
			OpPos: token.Pos(int(x.OpPos) + s.Offset()),
			Op:    x.Op,
			Y:     Expr(x.Y, s),
		}
	case *ast.KeyValueExpr:
		return &ast.KeyValueExpr{
			Key:   Expr(x.Key, s),
			Colon: token.Pos(int(x.Colon) + s.Offset()),
			Value: Expr(x.Value, s),
		}
	case *ast.ArrayType:
		return &ast.ArrayType{
			Lbrack: token.Pos(int(x.Lbrack) + s.Offset()),
			Len:    Expr(x.Len, s),
			Elt:    Expr(x.Elt, s),
		}
	case *ast.StructType:
		return &ast.StructType{
			Struct:     token.Pos(int(x.Struct) + s.Offset()),
			Fields:     FieldList(x.Fields, s),
			Incomplete: x.Incomplete,
		}
	case *ast.FuncType:
		return FuncType(x, s)
	case *ast.InterfaceType:
		return &ast.InterfaceType{
			Interface:  token.Pos(int(x.Interface) + s.Offset()),
			Methods:    FieldList(x.Methods, s),
			Incomplete: x.Incomplete,
		}
	case *ast.MapType:
		return &ast.MapType{
			Map:   token.Pos(int(x.Map) + s.Offset()),
			Key:   Expr(x.Key, s),
			Value: Expr(x.Value, s),
		}
	case *ast.ChanType:
		return &ast.ChanType{
			Begin: token.Pos(int(x.Begin) + s.Offset()),
			Arrow: token.Pos(int(x.Arrow) + s.Offset()),
			Dir:   x.Dir,
			Value: Expr(x.Value, s),
		}
	default:
		panic(fmt.Sprintf("invalid expr %q", x))
	}
}

// BasicLit :
func BasicLit(x *ast.BasicLit, s *State) *ast.BasicLit {
	if x == nil {
		return nil
	}
	return &ast.BasicLit{
		ValuePos: token.Pos(int(x.ValuePos) + s.Offset()),
		Kind:     x.Kind,
		Value:    x.Value,
	}

}

// BlockStmt :
func BlockStmt(x *ast.BlockStmt, s *State) *ast.BlockStmt {
	if x == nil {
		return nil
	}
	return &ast.BlockStmt{
		Lbrace: token.Pos(int(x.Lbrace) + s.Offset()),
		List:   Stmts(x.List, s),
		Rbrace: token.Pos(int(x.Rbrace) + s.Offset()),
	}
}

// CommentGroups :
func CommentGroups(xs []*ast.CommentGroup, s *State) []*ast.CommentGroup {
	ys := make([]*ast.CommentGroup, len(xs))
	for i := range xs {
		ys[i] = CommentGroup(xs[i], s)
	}
	return ys
}

// CommentGroup :
func CommentGroup(x *ast.CommentGroup, s *State) *ast.CommentGroup {
	if x == nil || len(x.List) == 0 {
		return nil
	}
	ys := make([]*ast.Comment, len(x.List))
	for i, x := range x.List {
		ys[i] = Comment(x, s)
	}
	// fmt.Printf("ccc %q %d -> %d\n", x.Text(), x.Pos(), ys[0].Pos())
	return &ast.CommentGroup{List: ys}
}

// Comment :
func Comment(x *ast.Comment, s *State) *ast.Comment {
	return &ast.Comment{
		Slash: token.Pos(int(x.Slash) + s.Offset()),
		Text:  x.Text,
	}
}
