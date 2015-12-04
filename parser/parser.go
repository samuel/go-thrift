// Copyright 2012-2015 Samuel Stauffer. All rights reserved.
// Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

package parser

//go:generate pigeon -o grammar.peg.go ./grammar.peg
//go:generate goimports -w ./grammar.peg.go

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Filesystem interface {
	Open(filename string) (io.ReadCloser, error)
	Abs(path string) (string, error)
}

type Parser struct {
	Filesystem Filesystem // For handling includes. Can be set to nil to fall back to os package.
}

func (p *Parser) Parse(r io.Reader, opts ...Option) (*Thrift, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	name := "<reader>"
	if named, ok := r.(namedReader); ok {
		name = named.Name()
	}
	t, err := Parse(name, b, opts...)
	if err != nil {
		return nil, err
	}
	return t.(*Thrift), nil
}

func (p *Parser) ParseFile(filename string) (map[string]*Thrift, string, error) {
	files := make(map[string]*Thrift)

	absPath, err := p.abs(filename)
	if err != nil {
		return nil, "", err
	}

	path := absPath
	for path != "" {
		rd, err := p.open(path)
		if err != nil {
			return nil, "", err
		}
		thrift, err := p.Parse(rd)
		if err != nil {
			return nil, "", err
		}
		files[path] = thrift

		basePath := filepath.Dir(path)
		for incName, incPath := range thrift.Includes {
			p, err := p.abs(filepath.Join(basePath, incPath))
			if err != nil {
				return nil, "", err
			}
			thrift.Includes[incName] = p
		}

		// Find path for next unparsed include
		path = ""
		for _, th := range files {
			for _, incPath := range th.Includes {
				if files[incPath] == nil {
					path = incPath
					break
				}
			}
		}
	}

	return files, absPath, nil
}

func (p *Parser) open(path string) (io.ReadCloser, error) {
	if p.Filesystem == nil {
		return os.Open(path)
	}
	return p.Filesystem.Open(path)
}

func (p *Parser) abs(path string) (string, error) {
	if p.Filesystem == nil {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return "", err
		}
		return filepath.Clean(absPath), nil
	}
	return p.Filesystem.Abs(path)
}

type namedReader interface {
	Name() string
}
