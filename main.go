package main

import (
	"io"
	"os"
)

type src struct {
	code []byte
	pos  int

	bkt  []int
	prev int
}

func newSrc(code string) *src {
	return &src{
		code: []byte(code),
		bkt:  make([]int, 0),
	}
}

func (s *src) Next() (byte, error) {
	if s.pos >= len(s.code) {
		return ' ', io.EOF
	}

	c := s.code[s.pos]
	if c == '[' {
		s.bkt = append(s.bkt, s.pos)
	}
	if c == ']' {
		s.prev = s.bkt[len(s.bkt)-1]
		s.bkt = s.bkt[:len(s.bkt)-1]
	}

	s.pos++
	return c, nil
}

func (s *src) Skip() {
	for f := s.bkt[len(s.bkt)-1]; f != s.prev; s.Next() {
	}
}

func (s *src) Back() {
	s.pos = s.prev
}

type brainfuck struct {
	*src
	mem [1 * 1024 * 1024]byte
	ptr int
}

func newBrainfuck(src string) *brainfuck {
	return &brainfuck{
		src: newSrc(src),
	}
}

func (bf *brainfuck) Loop() {
	var b byte
	var err error

	for {
		b, err = bf.Next()
		if err != nil {
			break
		}
		switch b {
		case '>':
			bf.ptr++
		case '<':
			bf.ptr--
		case '+':
			bf.mem[bf.ptr]++
		case '-':
			bf.mem[bf.ptr]--
		case '.':
			os.Stdin.Write(bf.mem[bf.ptr : bf.ptr+1])
		case ',':
			os.Stdout.Read(bf.mem[bf.ptr : bf.ptr+1])
		case '[':
			if bf.mem[bf.ptr] == 0 {
				bf.Skip()
			}
		case ']':
			if bf.mem[bf.ptr] != 0 {
				bf.Back()
			}
		}
	}
}

func main() {
	bf := newBrainfuck(os.Args[1])
	bf.Loop()
}
