package ilog

type IOption interface {
	do(ILog)
}

type OptionFunc func(ILog)

func (o OptionFunc) do(l ILog)  {
	o(l)
}

func EnableRotate() IOption {
	return OptionFunc(func(l ILog) {
		if f, ok := l.(*File); ok {
			f.enableRotate()
		}
	})
}

func EnableDivide() IOption {
	return OptionFunc(func(l ILog) {
		if f, ok := l.(*File); ok {
			f.enableDivide()
		}
	})
}

func WithFile(file string) IOption {
	return OptionFunc(func(l ILog) {
		if f, ok := l.(*File); ok {
			f.withFile(file)
		}
	})
}

func FileKeep(d int) IOption {
	return OptionFunc(func(l ILog) {
		if f, ok := l.(*File); ok {
			f.fileKeep(d)
		}
	})
}