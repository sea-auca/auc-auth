package utils

// Global defer implementation
type Shutdown struct {
	stack []func() error
}

// Create new global defer object
func NewShutdown() *Shutdown {
	return &Shutdown{}
}

// Add function to the defer stack
func (s *Shutdown) Add(f func() error) {
	s.stack = append(s.stack, f)
}

// Iterates over the stack, pops function and executes them
// if encounters error during close process - appends to the err array
func (s *Shutdown) Close() []error {
	var errs []error
	for i := len(s.stack) - 1; i < 0; i-- {
		err := s.stack[i]()
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
