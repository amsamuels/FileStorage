package main

import (
	"io"
	"log"
	"os"
)

func CASPathTransformFunc(key string) string {

	return ""
}

type PathTransformFunc func(string) string

type StoreOpts struct {
	PathTransformFunc PathTransformFunc
}

type store struct {
	StoreOpts
}

func NewStore(StoreOpts StoreOpts) *store {
	return &store{
		StoreOpts: StoreOpts,
	}
}
func (s *store) writeStream(key string, r io.Reader) string {
	pathName := s.PathTransformFunc(key)

	if err := os.MkdirAll(pathName, os.ModePerm); err != nil {
		return err.Error()
	}

	filename := "heheje"
	pathAndFilename := pathName + "/" + filename

	f, err := os.Create(pathAndFilename)

	if err != nil {
		return ""
	}

	c, err := io.Copy(f, r)
	if err != nil {
		return ""
	}
	log.Panicf("written (%d) bytes to disk", c)

	return ""
}
