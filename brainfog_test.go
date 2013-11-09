package main

import (
	"testing"
)

func Test_Helloworld(t *testing.T) {
	testProg := `
    +++++ +++++             initialize counter (cell #0) to 10
    [                       use loop to set the next four cells to 70/100/30/10
        > +++++ ++              add  7 to cell #1
        > +++++ +++++           add 10 to cell #2 
        > +++                   add  3 to cell #3
        > +                     add  1 to cell #4
        <<<< -                  decrement counter (cell #0)
    ]                   
    > ++ .                  print 'H'
    > + .                   print 'e'
    +++++ ++ .              print 'l'
    .                       print 'l'
    +++ .                   print 'o'
    > ++ .                  print ' '
    << +++++ +++++ +++++ .  print 'W'
    > .                     print 'o'
    +++ .                   print 'r'
    ----- - .               print 'l'
    ----- --- .             print 'd'
    > + .                   print '!'
    > .                     print '\n'`

	bf := newBrainfog([]byte(testProg))
	var res []byte
	for c := range bf.outCh {
		res = append(res, c)
	}

	wantedRes := "Hello World!\n"

	if string(res) != wantedRes {
		t.Errorf("Got %q, want %q", string(res), wantedRes)
	}
}
