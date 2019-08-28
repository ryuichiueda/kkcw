package main

import (
    "testing"
    "strings"
)

func check(input string, output string, candnum int, t *testing.T){
  result := ""
  for _, str := range strings.Split(input, "\n") {
      result += mainProc(str, candnum)
  }
  if result != output {
    t.Fatalf("failed test " + input + "\n" + result + " != " + output)
  }
}

func TestWords(t *testing.T) {
  check("やまだたろう", "山田太郎\n", 1, t)
  check("やまだたかお", "山田孝雄\n山田貴雄\n", 2, t)
  check("やまだたろう\nやまだたかお\nうえだりゅういち", "山田太郎\n山田孝雄\n上田隆一\n", 1, t)
  check("><&&><", "><&&><\n", 1, t)
}
