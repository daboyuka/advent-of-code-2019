package intcode

type IO interface {
	Input() int
	Output(int)
}

type SliceIO struct {
	Inputs  []int
	Outputs []int
}

func (s *SliceIO) Input() int {
	if len(s.Inputs) == 0 {
		panic("ran out of input")
	}
	v := s.Inputs[0]
	s.Inputs = s.Inputs[1:]
	return v
}

func (s *SliceIO) Output(v int) {
	s.Outputs = append(s.Outputs, v)
}
