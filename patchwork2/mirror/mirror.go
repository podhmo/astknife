package mirror

import (
	"fmt"
	"go/ast"
	"go/token"
)

// File :
func File(f *ast.File, offset int, ignoreComment bool) *ast.File {
	return &ast.File{
		Doc:        CommentGroup(f.Doc, offset, ignoreComment),
		Package:    token.Pos(int(f.Package) + offset),
		Name:       Ident(f.Name, offset, ignoreComment),
		Decls:      Decls(f.Decls, offset, ignoreComment),
		Scope:      ast.NewScope(nil), // todo:
		Imports:    ImportSpecs(f.Imports, offset, ignoreComment),
		Unresolved: Idents(f.Unresolved, offset, ignoreComment),
		Comments:   CommentGroups(f.Comments, offset, ignoreComment),
	}
}

// Decls :
func Decls(xs []ast.Decl, offset int, ignoreComment bool) []ast.Decl {
	ys := make([]ast.Decl, len(xs))
	for i := range xs {
		ys[i] = Decl(xs[i], offset, ignoreComment)
	}
	return ys
}

// Decl :
func Decl(decl ast.Decl, offset int, ignoreComment bool) ast.Decl {
	if decl == nil {
		return nil
	}
	switch x := decl.(type) {
	case *ast.GenDecl:
		new := *x
		new.Doc = CommentGroup(new.Doc, offset, ignoreComment)
		if new.Lparen != token.NoPos {
			new.Lparen = token.Pos(int(new.Lparen) + offset)
		}
		for i, spec := range new.Specs {
			new.Specs[i] = Spec(spec, offset, ignoreComment)
		}
		if new.Rparen != token.NoPos {
			new.Rparen = token.Pos(int(new.Rparen) + offset)
		}
		return &new
	case *ast.FuncDecl:
		return &ast.FuncDecl{
			Doc:  CommentGroup(x.Doc, offset, ignoreComment),
			Recv: FieldList(x.Recv, offset, ignoreComment),
			Name: Ident(x.Name, offset, ignoreComment),
			Type: FuncType(x.Type, offset, ignoreComment),
			Body: BlockStmt(x.Body, offset, ignoreComment),
		}
	case *ast.BadDecl:
		return &ast.BadDecl{
			From: token.Pos(int(x.From) + offset),
			To:   token.Pos(int(x.To) + offset),
		}
	default:
		panic(fmt.Sprintf("invalid decl %q", x))
	}
}

// Specs :
func Specs(xs []ast.Spec, offset int, ignoreComment bool) []ast.Spec {
	ys := make([]ast.Spec, len(xs))
	for i := range xs {
		ys[i] = Spec(xs[i], offset, ignoreComment)
	}
	return ys
}

// ImportSpecs :
func ImportSpecs(xs []*ast.ImportSpec, offset int, ignoreComment bool) []*ast.ImportSpec {
	ys := make([]*ast.ImportSpec, len(xs))
	for i := range xs {
		ys[i] = Spec(xs[i], offset, ignoreComment).(*ast.ImportSpec)
	}
	return ys
}

// Spec :
func Spec(spec ast.Spec, offset int, ignoreComment bool) ast.Spec {
	if spec == nil {
		return nil
	}
	switch x := spec.(type) {
	case *ast.ImportSpec:
		endpos := x.EndPos
		if endpos != token.NoPos {
			endpos = token.Pos(int(x.EndPos) + offset)
		}
		return &ast.ImportSpec{
			Doc:     CommentGroup(x.Doc, offset, ignoreComment),
			Name:    Ident(x.Name, offset, ignoreComment),
			Path:    BasicLit(x.Path, offset, ignoreComment),
			Comment: CommentGroup(x.Comment, offset, ignoreComment),
			EndPos:  endpos,
		}

	case *ast.ValueSpec:
		return &ast.ValueSpec{
			Doc:     CommentGroup(x.Doc, offset, ignoreComment),
			Names:   Idents(x.Names, offset, ignoreComment),
			Type:    Expr(x.Type, offset, ignoreComment),
			Values:  Exprs(x.Values, offset, ignoreComment),
			Comment: CommentGroup(x.Comment, offset, ignoreComment),
		}
	case *ast.TypeSpec:
		assign := x.Assign
		if assign != token.NoPos {
			assign = token.Pos(int(x.Assign) + offset)
		}
		return &ast.TypeSpec{
			Doc:     CommentGroup(x.Doc, offset, ignoreComment),
			Name:    Ident(x.Name, offset, ignoreComment),
			Assign:  assign,
			Type:    Expr(x.Type, offset, ignoreComment),
			Comment: CommentGroup(x.Comment, offset, ignoreComment),
		}
	default:
		panic(fmt.Sprintf("invalid spec %q", x))
	}
}

// FieldList :
func FieldList(x *ast.FieldList, offset int, ignoreComment bool) *ast.FieldList {
	if x == nil {
		return nil
	}
	opening := x.Opening
	if opening != token.NoPos {
		opening = token.Pos(int(x.Opening) + offset)
	}
	closing := x.Closing
	if closing != token.NoPos {
		closing = token.Pos(int(x.Closing) + offset)
	}
	return &ast.FieldList{
		Opening: opening,
		List:    Fields(x.List, offset, ignoreComment),
		Closing: closing,
	}
}

// Fields :
func Fields(xs []*ast.Field, offset int, ignoreComment bool) []*ast.Field {
	ys := make([]*ast.Field, len(xs))
	for i := range xs {
		ys[i] = Field(xs[i], offset, ignoreComment)
	}
	return ys
}

// Field :
func Field(x *ast.Field, offset int, ignoreComment bool) *ast.Field {
	if x == nil {
		return nil
	}
	return &ast.Field{
		Doc:     CommentGroup(x.Doc, offset, ignoreComment),
		Names:   Idents(x.Names, offset, ignoreComment),
		Type:    Expr(x.Type, offset, ignoreComment),
		Tag:     BasicLit(x.Tag, offset, ignoreComment),
		Comment: CommentGroup(x.Comment, offset, ignoreComment),
	}
}

// Idents :
func Idents(xs []*ast.Ident, offset int, ignoreComment bool) []*ast.Ident {
	ys := make([]*ast.Ident, len(xs))
	for i := range xs {
		ys[i] = Ident(xs[i], offset, ignoreComment)
	}
	return ys
}

// Ident :
func Ident(x *ast.Ident, offset int, ignoreComment bool) *ast.Ident {
	if x == nil {
		return nil
	}
	return &ast.Ident{
		NamePos: token.Pos(int(x.NamePos) + offset),
		Name:    x.Name,
		// Object
	}
}

// FuncType :
func FuncType(x *ast.FuncType, offset int, ignoreComment bool) *ast.FuncType {
	if x == nil {
		return nil
	}
	f := x.Func
	if f != token.NoPos {
		f = token.Pos(int(x.Func) + offset)
	}
	return &ast.FuncType{
		Func:    f,
		Params:  FieldList(x.Params, offset, ignoreComment),
		Results: FieldList(x.Results, offset, ignoreComment),
	}
}

// Stmts :
func Stmts(xs []ast.Stmt, offset int, ignoreComment bool) []ast.Stmt {
	ys := make([]ast.Stmt, len(xs))
	for i := range xs {
		ys[i] = Stmt(xs[i], offset, ignoreComment)
	}
	return ys
}

// Stmt :
func Stmt(stmt ast.Stmt, offset int, ignoreComment bool) ast.Stmt {
	if stmt == nil {
		return nil
	}
	switch x := stmt.(type) {
	case *ast.BadStmt:
		return &ast.BadStmt{
			From: token.Pos(int(x.From) + offset),
			To:   token.Pos(int(x.To) + offset),
		}
	case *ast.DeclStmt:
		return &ast.DeclStmt{
			Decl: Decl(x.Decl, offset, ignoreComment),
		}
	case *ast.EmptyStmt:
		return &ast.EmptyStmt{
			Semicolon: token.Pos(int(x.Semicolon) + offset),
			Implicit:  x.Implicit,
		}
	case *ast.LabeledStmt:
		return &ast.LabeledStmt{
			Label: Ident(x.Label, offset, ignoreComment),
			Colon: token.Pos(int(x.Colon) + offset),
			Stmt:  Stmt(x.Stmt, offset, ignoreComment),
		}
	case *ast.ExprStmt:
		return &ast.ExprStmt{
			X: Expr(x.X, offset, ignoreComment),
		}
	case *ast.SendStmt:
		return &ast.SendStmt{
			Chan:  Expr(x.Chan, offset, ignoreComment),
			Arrow: token.Pos(int(x.Arrow) + offset),
			Value: Expr(x.Value, offset, ignoreComment),
		}
	case *ast.IncDecStmt:
		return &ast.IncDecStmt{
			X:      Expr(x.X, offset, ignoreComment),
			TokPos: token.Pos(int(x.TokPos) + offset),
			Tok:    x.Tok,
		}
	case *ast.AssignStmt:
		return &ast.AssignStmt{
			Lhs:    Exprs(x.Lhs, offset, ignoreComment),
			TokPos: token.Pos(int(x.TokPos) + offset),
			Tok:    x.Tok,
			Rhs:    Exprs(x.Rhs, offset, ignoreComment),
		}
	case *ast.GoStmt:
		return &ast.GoStmt{
			Go:   token.Pos(int(x.Go) + offset),
			Call: CallExpr(x.Call, offset, ignoreComment),
		}
	case *ast.DeferStmt:
		return &ast.DeferStmt{
			Defer: token.Pos(int(x.Defer) + offset),
			Call:  CallExpr(x.Call, offset, ignoreComment),
		}
	case *ast.ReturnStmt:
		return &ast.ReturnStmt{
			Return:  token.Pos(int(x.Return) + offset),
			Results: Exprs(x.Results, offset, ignoreComment),
		}
	case *ast.BranchStmt:
		return &ast.BranchStmt{
			TokPos: token.Pos(int(x.TokPos) + offset),
			Tok:    x.Tok,
			Label:  Ident(x.Label, offset, ignoreComment),
		}
	case *ast.BlockStmt:
		return BlockStmt(x, offset, ignoreComment)

	case *ast.IfStmt:
		return &ast.IfStmt{
			If:   token.Pos(int(x.If) + offset),
			Init: Stmt(x.Init, offset, ignoreComment),
			Cond: Expr(x.Cond, offset, ignoreComment),
			Body: BlockStmt(x.Body, offset, ignoreComment),
			Else: Stmt(x.Else, offset, ignoreComment),
		}
	case *ast.CaseClause:
		return &ast.CaseClause{
			Case:  token.Pos(int(x.Case) + offset),
			List:  Exprs(x.List, offset, ignoreComment),
			Colon: token.Pos(int(x.Colon) + offset),
			Body:  Stmts(x.Body, offset, ignoreComment),
		}
	case *ast.SwitchStmt:
		return &ast.SwitchStmt{
			Switch: token.Pos(int(x.Switch) + offset),
			Init:   Stmt(x.Init, offset, ignoreComment),
			Tag:    Expr(x.Tag, offset, ignoreComment),
			Body:   BlockStmt(x.Body, offset, ignoreComment),
		}
	case *ast.TypeSwitchStmt:
		return &ast.TypeSwitchStmt{
			Switch: token.Pos(int(x.Switch) + offset),
			Init:   Stmt(x.Init, offset, ignoreComment),
			Assign: Stmt(x.Assign, offset, ignoreComment),
			Body:   BlockStmt(x.Body, offset, ignoreComment),
		}
	case *ast.CommClause:
		return &ast.CommClause{
			Case:  token.Pos(int(x.Case) + offset),
			Comm:  Stmt(x.Comm, offset, ignoreComment),
			Colon: token.Pos(int(x.Colon) + offset),
			Body:  Stmts(x.Body, offset, ignoreComment),
		}
	case *ast.SelectStmt:
		return &ast.SelectStmt{
			Select: token.Pos(int(x.Select) + offset),
			Body:   BlockStmt(x.Body, offset, ignoreComment),
		}
	case *ast.ForStmt:
		return &ast.ForStmt{
			For:  token.Pos(int(x.For) + offset),
			Init: Stmt(x.Init, offset, ignoreComment),
			Cond: Expr(x.Cond, offset, ignoreComment),
			Post: Stmt(x.Post, offset, ignoreComment),
			Body: BlockStmt(x.Body, offset, ignoreComment),
		}
	case *ast.RangeStmt:
		return &ast.RangeStmt{
			For:    token.Pos(int(x.For) + offset),
			Key:    Expr(x.Key, offset, ignoreComment),
			Value:  Expr(x.Value, offset, ignoreComment),
			TokPos: token.Pos(int(x.TokPos) + offset),
			Tok:    x.Tok,
			X:      Expr(x.X, offset, ignoreComment),
			Body:   BlockStmt(x.Body, offset, ignoreComment),
		}
	default:
		panic(fmt.Sprintf("invalid stmt %q", x))
	}
}

// CallExpr :
func CallExpr(x *ast.CallExpr, offset int, ignoreComment bool) *ast.CallExpr {
	if x == nil {
		return nil
	}
	return &ast.CallExpr{
		Fun:      Expr(x.Fun, offset, ignoreComment),
		Lparen:   token.Pos(int(x.Lparen) + offset),
		Args:     Exprs(x.Args, offset, ignoreComment),
		Ellipsis: token.Pos(int(x.Ellipsis) + offset),
		Rparen:   token.Pos(int(x.Rparen) + offset),
	}
}

// Exprs :
func Exprs(xs []ast.Expr, offset int, ignoreComment bool) []ast.Expr {
	ys := make([]ast.Expr, len(xs))
	for i := range xs {
		ys[i] = Expr(xs[i], offset, ignoreComment)
	}
	return ys
}

// Expr :
func Expr(expr ast.Expr, offset int, ignoreComment bool) ast.Expr {
	if expr == nil {
		return nil
	}
	switch x := expr.(type) {
	case *ast.BadExpr:
		return &ast.BadExpr{
			From: token.Pos(int(x.From) + offset),
			To:   token.Pos(int(x.To) + offset),
		}
	case *ast.Ident:
		return Ident(x, offset, ignoreComment)
	case *ast.Ellipsis:
		return &ast.Ellipsis{
			Ellipsis: token.Pos(int(x.Ellipsis) + offset),
			Elt:      Expr(x.Elt, offset, ignoreComment),
		}
	case *ast.BasicLit:
		return BasicLit(x, offset, ignoreComment)
	case *ast.FuncLit:
		return &ast.FuncLit{
			Type: FuncType(x.Type, offset, ignoreComment),
			Body: BlockStmt(x.Body, offset, ignoreComment),
		}
	case *ast.CompositeLit:
		return &ast.CompositeLit{
			Type:   Expr(x.Type, offset, ignoreComment),
			Lbrace: token.Pos(int(x.Lbrace) + offset),
			Elts:   Exprs(x.Elts, offset, ignoreComment),
			Rbrace: token.Pos(int(x.Rbrace) + offset),
		}
	case *ast.ParenExpr:
		return &ast.ParenExpr{
			Lparen: token.Pos(int(x.Lparen) + offset),
			X:      Expr(x.X, offset, ignoreComment),
			Rparen: token.Pos(int(x.Rparen) + offset),
		}
	case *ast.SelectorExpr:
		return &ast.SelectorExpr{
			X:   Expr(x.X, offset, ignoreComment),
			Sel: Ident(x.Sel, offset, ignoreComment),
		}
	case *ast.IndexExpr:
		return &ast.IndexExpr{
			X:      Expr(x.X, offset, ignoreComment),
			Lbrack: token.Pos(int(x.Lbrack) + offset),
			Index:  Expr(x.Index, offset, ignoreComment),
			Rbrack: token.Pos(int(x.Rbrack) + offset),
		}
	case *ast.SliceExpr:
		return &ast.SliceExpr{
			X:      Expr(x.X, offset, ignoreComment),
			Lbrack: token.Pos(int(x.Lbrack) + offset),
			Low:    Expr(x.Low, offset, ignoreComment),
			High:   Expr(x.High, offset, ignoreComment),
			Max:    Expr(x.Max, offset, ignoreComment),
			Slice3: x.Slice3,
			Rbrack: token.Pos(int(x.Rbrack) + offset),
		}
	case *ast.TypeAssertExpr:
		return &ast.TypeAssertExpr{
			X:      Expr(x.X, offset, ignoreComment),
			Lparen: token.Pos(int(x.Lparen) + offset),
			Type:   Expr(x.Type, offset, ignoreComment),
			Rparen: token.Pos(int(x.Rparen) + offset),
		}
	case *ast.CallExpr:
		return CallExpr(x, offset, ignoreComment)
	case *ast.StarExpr:
		return &ast.StarExpr{
			Star: token.Pos(int(x.Star) + offset),
			X:    Expr(x.X, offset, ignoreComment),
		}
	case *ast.UnaryExpr:
		return &ast.UnaryExpr{
			OpPos: token.Pos(int(x.OpPos) + offset),
			Op:    x.Op,
			X:     Expr(x.X, offset, ignoreComment),
		}
	case *ast.BinaryExpr:
		return &ast.BinaryExpr{
			X:     Expr(x.X, offset, ignoreComment),
			OpPos: token.Pos(int(x.OpPos) + offset),
			Op:    x.Op,
			Y:     Expr(x.Y, offset, ignoreComment),
		}
	case *ast.KeyValueExpr:
		return &ast.KeyValueExpr{
			Key:   Expr(x.Key, offset, ignoreComment),
			Colon: token.Pos(int(x.Colon) + offset),
			Value: Expr(x.Value, offset, ignoreComment),
		}
	default:
		panic(fmt.Sprintf("invalid expr %q", x))
	}
}

// BasicLit :
func BasicLit(x *ast.BasicLit, offset int, ignoreComment bool) *ast.BasicLit {
	if x == nil {
		return nil
	}
	return &ast.BasicLit{
		ValuePos: token.Pos(int(x.ValuePos) + offset),
		Kind:     x.Kind,
		Value:    x.Value,
	}

}

// BlockStmt :
func BlockStmt(x *ast.BlockStmt, offset int, ignoreComment bool) *ast.BlockStmt {
	if x == nil {
		return nil
	}
	return &ast.BlockStmt{
		Lbrace: token.Pos(int(x.Lbrace) + offset),
		List:   Stmts(x.List, offset, ignoreComment),
		Rbrace: token.Pos(int(x.Rbrace) + offset),
	}
}

// CommentGroups :
func CommentGroups(xs []*ast.CommentGroup, offset int, ignoreComment bool) []*ast.CommentGroup {
	ys := make([]*ast.CommentGroup, len(xs))
	for i := range xs {
		ys[i] = CommentGroup(xs[i], offset, ignoreComment)
	}
	return ys
}

// CommentGroup :
func CommentGroup(x *ast.CommentGroup, offset int, ignoreComment bool) *ast.CommentGroup {
	if ignoreComment {
		return x
	}

	if x == nil || len(x.List) == 0 {
		return nil
	}
	ys := make([]*ast.Comment, len(x.List))
	for i, x := range x.List {
		ys[i] = Comment(x, offset, ignoreComment)
	}
	return &ast.CommentGroup{List: ys}
}

// Comment :
func Comment(x *ast.Comment, offset int, ignoreComment bool) *ast.Comment {
	return &ast.Comment{
		Slash: token.Pos(int(x.Slash) + offset),
		Text:  x.Text,
	}
}
