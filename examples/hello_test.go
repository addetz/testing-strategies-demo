package hello_test

import (
	"fmt"
	"testing"
)

func SayHello(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

func TestSayHello(t *testing.T) {
	names := []string{"Anna", "Belle"}
	for _, n := range names {
		t.Run(n, func(t *testing.T) {
			want := "Hello, " + n + "!"
			got := SayHello(n)
			if got != want {
				t.Fatalf("got: %s; want: %s",
					got, want)
			}
		})
	}
}

func FuzzSayHello(f *testing.F) {
  f.Add("Anna")
  f.Add("Belle")
  f.Fuzz(func(t *testing.T, n string) {
    want := "Hello, " + n + "!"    
    got := SayHello(n)
    if got != want {
      t.Fatalf("got: %s; want: %s", 
        got, want)
    }
  })
}
