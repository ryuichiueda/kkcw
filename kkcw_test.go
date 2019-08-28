package main

import "testing"

func checkSimpleReplacement(input string, output string, t *testing.T){
  result := mainProc(input)
  if result != output {
    t.Fatalf("failed test " + input + "\n" + result + " != " + output)
  }
}

func TestWords(t *testing.T) {
  checkSimpleReplacement("やまだたろう", "山田太郎", t)
}
