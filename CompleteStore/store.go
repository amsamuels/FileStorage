package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

type PathTransformFunc func(string) PathKey

const defaultRootFolderName = "whatibuilt"

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		Pathname: key,
		Filename: key,
	}
}

type StoreOpts struct {
	// Root is the folder name of the root, containing all the folders/files of the system
	Root              string
	PathTransformFunc PathTransformFunc
}

type store struct {
	StoreOpts
}

func NewStore(StoreOpts StoreOpts) *store {
	if StoreOpts.PathTransformFunc == nil {
		StoreOpts.PathTransformFunc = DefaultPathTransformFunc
	}

	if len(StoreOpts.Root) == 0 {
		StoreOpts.Root = defaultRootFolderName
	}
	return &store{
		StoreOpts: StoreOpts,
	}
}

func (s *store) Has(key string) bool {
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())
	_, err := os.Stat(fullPathWithRoot)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func (s *store) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)
	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.Filename)
	}()
	firstPathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.firstPathName())
	return os.RemoveAll(firstPathNameWithRoot)
}

func (s *store) Read(key string) (io.Reader, error) {
	f, err := s.ReadStream(key)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)

	return buf, err
}
func (s *store) ReadStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	pathKeyWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())
	return os.Open(pathKeyWithRoot)
}

func (s *store) Write(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.Pathname)
	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return err
	}

	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())
	f, err := os.Create(fullPathWithRoot)
	if err != nil {
		return err
	}

	c, err := io.Copy(f, r)
	if err != nil {
		return err
	}
	log.Printf("written (%d) bytes to disk", c)
	return nil
}
