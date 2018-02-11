package patchwork

import (
	"go/ast"
	"sync"

	"github.com/podhmo/astknife/lookup"
)

// scope :
type scope struct {
	elems map[string]*lookup.Result // lazily allocated
	files map[string]*ast.File
	sync.RWMutex
}

func newscope() *scope {
	return &scope{
		elems: map[string]*lookup.Result{},
		files: map[string]*ast.File{},
	}
}

// AddFile :
func (s *scope) AddFile(filename string, file *ast.File) {
	s.Lock()
	defer s.Unlock()
	s.files[filename] = file
}

// Lookup :
func (s *scope) Lookup(name string, file *ast.File, lookupfn func(*ast.File, string) *lookup.Result) *lookup.Result {
	s.RLock()
	r, ok := s.elems[name]
	s.RUnlock()
	if ok {
		return r
	}
	if file != nil {
		if r := lookupfn(file, name); r != nil {
			s.Lock()
			s.elems[name] = r
			s.Unlock()
			return r
		}
	}
	for _, f := range s.files {
		if r := lookupfn(f, name); r != nil {
			s.Lock()
			s.elems[name] = r
			s.Unlock()
			return r
		}
	}
	return nil
}
