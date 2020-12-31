package intcode

type IO interface {
	Input() int
	Output(int)
	Done()
}

type SliceIO struct {
	Inputs  []int
	Outputs []int
	IsDone  bool
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

func (s SliceIO) Done() { s.IsDone = true }

type ChanIO struct {
	Inputs  chan int
	Outputs chan int
	IsDone  chan struct{}
}

func (c *ChanIO) Input() int {
	v, ok := <-c.Inputs
	if !ok {
		panic("empty channel")
	}
	return v
}

func (c *ChanIO) Output(v int) {
	c.Outputs <- v
}

func (c *ChanIO) Done() {
	close(c.Outputs)
	if c.IsDone != nil {
		close(c.IsDone)
	}
}
