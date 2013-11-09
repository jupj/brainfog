// brainfog project main.go
package main

import (
	"bytes"
	"os"
	"path/filepath"
	"fmt"
	"io/ioutil"
)

// brainfog is an interpreter for the brainf*** language
type brainfog struct {
	ip      int         // instruction pointer
	program []byte      // instruction sequence
	cp      int         // cell pointer
	cell    [30000]byte // program memory
	inCh    chan byte   // input channel
	outCh   chan byte   // output channel
}

// newBrainfog creates a new *brainfog for the source code bfSrc
func newBrainfog(bfSrc []byte) *brainfog {
	bf := &brainfog{inCh: make(chan byte), outCh: make(chan byte)}

	// Pick the instructions from the source and add them to the program
	instructions := []byte("+-<>,.[]")
	for _, c := range bfSrc {
		if bytes.Contains(instructions, []byte{c}) {
			bf.program = append(bf.program, c)
		}
	}

	// Run the program
	go bf.run()
	return bf
}

// doBranch executes all the instructions of a branch/loop
func (bf *brainfog) doBranch() error {
	if bf.program[bf.ip] != '[' {
		return fmt.Errorf("doBranch: invalid start index: %d", bf.ip)
	}

	// store start and end indices for the loop
	start := bf.ip
	end, err := bf.findEnd(start)
	if err != nil {
		return err
	}

	for bf.ip <= end {
		if bf.ip == start {
			// At the beginning of the loop
			if bf.cell[bf.cp] == 0 {
				// No flag: Jump out of the loop
				bf.ip = end
				break
			} else {
				// Enter the loop
				bf.ip++
			}
		}

		if bf.ip == end {
			// End of loop: jump back to start of loop
			bf.ip = start
		} else {
			// Normal instruction
			err = bf.doInstruction()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// findEnd finds the index of the corresponding ] for a given [
func (bf *brainfog) findEnd(start int) (int, error) {
	if bf.program[start] != '[' {
		return 0, fmt.Errorf("findEnd: invalid start index: %d", start)
	}

	var err error

	for i := start + 1; i < len(bf.program); i++ {
		switch bf.program[i] {
		case ']':
			return i, nil
		case '[':
			// Found inner loop, call findEnd recursively for this loop
			i, err = bf.findEnd(i)
			if err != nil {
				return 0, err
			}
		}
	}
	return 0, fmt.Errorf("findEnd: no matching ] found for index: %d", start)
}

// doInstruction executes an instrution and increments the instrution pointer
func (bf *brainfog) doInstruction() error {
	switch bf.program[bf.ip] {
	case '+':
		bf.cell[bf.cp]++
	case '-':
		bf.cell[bf.cp]--
	case '<':
		if bf.cp == 0 {
			return fmt.Errorf("Cell pointer underflow at instruction %d", bf.ip)
		}
		bf.cp--
	case '>':
		if bf.cp == (len(bf.cell) - 1) {
			return fmt.Errorf("Cell pointer overflow at instruction %d", bf.ip)
		}
		bf.cp++
	case '.':
		bf.outCh <- bf.cell[bf.cp]
	case ',':
		bf.cell[bf.cp] = <-bf.inCh
	case '[':
		err := bf.doBranch()
		if err != nil {
			return err
		}
	case ']':
		return fmt.Errorf("Unmatched ] at index %d", bf.ip)
	}
	bf.ip++

	return nil
}

// run executes the instructions of the program
func (bf *brainfog) run() {
	for bf.ip < len(bf.program) {
		err := bf.doInstruction()
		if err != nil {
			panic(err)
		}
	}
	close(bf.outCh)
}

func main() {
	if len(os.Args) > 1 {
		// Read bf code from file
		filename := os.Args[1]
		code, err := ioutil.ReadFile(filename)
		if err != nil {
			panic(fmt.Errorf("Unable to open file %v", filename))
		}

		bf := newBrainfog(code)
		for c := range bf.outCh {
			fmt.Printf("%c", c)
		}
	} else {
		fmt.Printf("Usage: %s <sourcefile>\n", filepath.Base(os.Args[0]))
	}
}
