package main

import (
	"errors"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"strings"

	"github.com/Tomoka64/go-pkg-seeker/model"
)

//importPkg imports a path to a directory where the fname is located.
func importPkg(fname, dir string) (*build.Package, error) {
	p, err := build.Import(fname, dir, build.ImportComment)
	if err != nil {
		return &build.Package{}, err
	}
	if p.BinaryOnly {
		return &build.Package{}, errors.New("it consists of binary only")
	}
	if p.IsCommand() {
		return &build.Package{}, errors.New("the package is considered a command to be installed (not just a library)")
	}
	return p, nil
}

//extractWord gets filename and pattern and looks for the result accordingly and puts the found results
//into datas([]model.Result) and returns it.
func extractWord(fname, pattern string, datas []model.Result) ([]model.Result, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fname, nil, parser.ParseComments)
	if err != nil {
		return []model.Result{}, err
	}

	cmap := ast.NewCommentMap(fset, f, f.Comments)
	for n, cgs := range cmap {
		f := fset.File(n.Pos())
		for _, cg := range cgs {
			t := cg.Text()
			if strings.Contains(t, pattern) {
				a := f.Position(cg.Pos()).Line
				datas = append(datas, model.NewResult(fname, pattern, t, a))
			}
		}
	}

	return datas, nil
}
