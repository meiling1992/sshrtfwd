package single

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrorAleadyRunningInstance = errors.New("another instance is already running")
)

type Single struct {
	Name string
	Path string
	File *os.File
}

func New(name string, opts ...Option) (*Single, error) {
	if name != "" {
		s := &Single{
			Name: name,
		}
		for _, opt := range opts {
			opt(s)
		}
		if s.Path == "" {
			s.Path = os.TempDir()
		}

		return s, nil
	}
	return nil, fmt.Errorf("name cannot be empty")
}

func (s *Single) lockfile() string {
	return filepath.Join(s.Path, fmt.Sprintf("%s.lock", s.Name))
}

func (s *Single) Lock() error {
	if err := os.Remove(s.lockfile()); err != nil && !os.IsNotExist(err) {
		return ErrorAleadyRunningInstance
	}
	file, err := os.OpenFile(s.lockfile(), os.O_EXCL|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	s.File = file
	return nil
}
func (s *Single) UnLocl() error {
	if err := s.File.Close(); err != nil {
		return fmt.Errorf("failed to close the file :%w", err)
	}
	if err := os.Remove(s.lockfile()); err != nil {
		return fmt.Errorf("failed to remove the lock file :%w", err)
	}
	return nil
}
