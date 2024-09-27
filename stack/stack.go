package stack

type Stack struct {
	items []rune
}

func (s *Stack) Push(item rune) {
	s.items = append(s.items, item)
}

func (s *Stack) Pop() (rune, bool) {
	if s.IsEmpty() {
		var zeroValue rune
		return zeroValue, false
	}
	index := len(s.items) - 1
	item := s.items[index]
	s.items = s.items[:index]
	return item, true
}

func (s *Stack) Peek() (rune, bool) {
	if s.IsEmpty() {
		var zeroValue rune
		return zeroValue, false
	}
	return s.items[len(s.items)-1], true
}

func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *Stack) Size() int {
	return len(s.items)
}
