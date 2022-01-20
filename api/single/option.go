package single

type Option func(*Single)

func WithLockPath(lockpath string) Option {
	return func(s *Single) {
		s.Path = lockpath
	}
}
