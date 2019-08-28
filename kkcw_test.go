package main

import "testing"

func check(input string, output string, candnum int, t *testing.T){
  result := mainProc(input, candnum)
  if result != output {
    t.Fatalf("failed test " + input + "\n" + result + " != " + output)
  }
}

func TestWords(t *testing.T) {
  check("やまだたろう", "山田太郎\n", 1, t)
  check("やまだたかお", "山田孝雄\n山田貴雄\n", 2, t)
}
