package store

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
)

type GeneratePath func(string) Path

const defaultRootDir = "dfs_disk"

func GenerateCASPath(key string) Path {
	hash := sha1.Sum([]byte(key))
	hashstr := hex.EncodeToString(hash[:])
	blocksize := 5
	path := hashstr[0:min(len(hashstr), blocksize)]
	for i := blocksize; i+blocksize <= len(hashstr); i += blocksize {
		path = path + "/" + hashstr[i:i+blocksize]
	}
	return Path{
		dirPath:  path,
		fileName: hashstr,
	}
}

func DefaultGeneratePath(key string) Path {
	blocksize := 5
	path := key[0:min(len(key), blocksize)]
	for i := blocksize; i+blocksize <= len(key); i += blocksize {
		path = path + "/" + key[i:i+blocksize]
	}
	return Path{
		dirPath:  path,
		fileName: key,
	}
}

type Path struct {
	dirPath  string
	fileName string
}

func (path Path) fullPath() string {
	return fmt.Sprintf("%s/%s", path.dirPath, path.fileName)
}

type StoreOptions struct {
	RootDir      string
	GeneratePath GeneratePath
}

type Store struct {
	StoreOptions
}

func NewStore(opts StoreOptions) *Store {
	if len(opts.RootDir) == 0 {
		opts.RootDir = defaultRootDir
	}
	if opts.GeneratePath == nil {
		opts.GeneratePath = DefaultGeneratePath
	}
	if opts.RootDir[len(opts.RootDir)-1] != '/' {
		opts.RootDir = opts.RootDir + "/"
	}
	return &Store{
		StoreOptions: opts,
	}
}

func (s *Store) Has(key string) bool {
	path := s.GeneratePath(key)
	_, err := os.Stat(s.RootDir + path.fullPath())
	return !os.IsNotExist(err)
}

func (s *Store) Write(key string, r io.Reader) error {
	path := s.GeneratePath(key)
	if err := os.MkdirAll(s.RootDir+path.dirPath, os.ModePerm); err != nil {
		log.Printf("error writing file %s", err)
		return err
	}
	fullPath := s.RootDir + path.fullPath()

	file, err := os.Create(fullPath)
	if err != nil {
		log.Printf("error Creating file %s", err)
		return err
	}
	defer file.Close()
	n, err := io.Copy(file, r)
	if err != nil {
		log.Printf("error Copying file %s", err)
		return err
	}

	log.Printf("written %d bytes successfully to disk at %s", n, fullPath)

	return nil
}

func (s *Store) readToBuff(key string) (io.Reader, error) {
	path := s.GeneratePath(key)
	r, err := os.Open(s.RootDir + path.fullPath())
	if err != nil {
		log.Printf("failed to open file %s", err)
		return nil, err
	}
	defer r.Close()
	buff := new(bytes.Buffer)
	_, err = io.Copy(buff, r)
	if err != nil {
		log.Printf("error copying to buffer %s", err)
		return nil, err
	}
	return buff, err

}

func (s *Store) Read(key string) (io.Reader, error) {
	return s.readToBuff(key)
}

func (s *Store) Delete(key string) error {
	path := s.GeneratePath(key)
	if err := os.Remove(s.RootDir + path.fullPath()); err != nil {
		return err
	}
	pathstr := path.dirPath
	for len(pathstr) > 0 {
		files, err := os.ReadDir(s.RootDir + pathstr)
		if err != nil {
			return err
		}
		if len(files) == 0 {
			if err := os.Remove(s.RootDir + pathstr); err != nil {
				return err
			}
		}
		//blocksize = 5 + /(1)
		pathstr = pathstr[0:max(0, len(pathstr)-6)]
	}
	return nil
}

func (s *Store) Clear() error {
	return os.RemoveAll(s.RootDir)
}
