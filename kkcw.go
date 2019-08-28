// license and copyright => see help()
package main

import (
    "fmt"
    "os"
    "os/exec"
    "bufio"
    "io"
    "strings"
    "regexp"
    "strconv"
)

const VERSION = "0.3.3"

func help() {
    fmt.Fprintf(os.Stderr, "kkcw: kkc wrapper %s\n", VERSION)
    fmt.Fprintln(os.Stderr, "Copyright (C) 2019 Ryuichi Ueda.")
    fmt.Fprintln(os.Stderr, "")
    fmt.Fprintln(os.Stderr, "usage: echo <string> | kkcw")
    fmt.Fprintln(os.Stderr, "")
    fmt.Fprintln(os.Stderr, "Released under the GPLv3")
    fmt.Fprintln(os.Stderr, "https://github.com/ryuichiueda/kkcw")
}

type Token struct {
    token_id int
    results []string
}

func parse(text string, token_id int) Token {
    resultmask := regexp.MustCompile(`[^:]*: `)
    resultline := regexp.MustCompile(`\d+: *`)
    orgstrpart := regexp.MustCompile(`/[^>]+>`)

    r := Token{}
    r.token_id = token_id

    for _, str := range strings.Split(text, "\n") {
	str = strings.Replace(str, ">> ", "", -1)
        if ! resultline.MatchString(str) {
	   continue
        }

	str = resultmask.ReplaceAllString(str, "")
	str = orgstrpart.ReplaceAllString(str, "")
	str = strings.Replace(str, "<", "", -1)

	r.results = append(r.results, descape(str))
    }
    return r
}

func kkc(token string) string {
    cmd := exec.Command("kkc")
    stdin, _ := cmd.StdinPipe()
    io.WriteString(stdin, token)
    stdin.Close()

    out, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Fprintln(os.Stderr, "Command Exec Error.")
        fmt.Fprintln(os.Stderr, string(out))
	os.Exit(1)
    }

    return string(out)
}

func escape(text string) string {
    tmp1 := strings.Replace(text, "&", "&amp;", -1)
    tmp2 := strings.Replace(tmp1, ">", "&gt;", -1)
    return strings.Replace(tmp2, "<", "&lt;", -1)
}

func descape(text string) string {
    tmp1 := strings.Replace(text, "&lt;", "<", -1)
    tmp2 := strings.Replace(tmp1, "&gt;", ">", -1)
    return strings.Replace(tmp2, "&amp;", "&", -1)
}

func concat(tokens []Token, rank int) string {
    ans := tokens[0].results[rank]
    for _, str := range tokens[1:] {
        ans += " " + str.results[rank]
    }
    return ans + "\n"
}

func mainProc(str string, candidate_num int) string {
    tokens := make([]Token, 0)

    for n, str := range strings.Split(str, " ") {
        out := kkc(escape(str) + " " + strconv.Itoa(candidate_num))
	tokens = append(tokens, parse(out, n))
    }

    var ans string
    for n := 0; n < candidate_num; n++ {
        ans += concat(tokens, n)
    }
    return ans
}

func candNum() int {
    if len(os.Args) != 3 {
        return 1
    }

    if os.Args[1] != "-n" {
        fmt.Fprintln(os.Stderr, "Invalid Option")
        os.Exit(1)
    }

    cand_num, err := strconv.Atoi(os.Args[2])
    if err != nil {
        fmt.Fprintln(os.Stderr, "Invalid Number")
        os.Exit(1)
    }
    return cand_num
}

func main() {
    if len(os.Args) == 2 || len(os.Args) > 3 {
        help()
        os.Exit(0)
    }

    stdin := bufio.NewScanner(os.Stdin)
    cand_num := candNum()
    for stdin.Scan() {
        result := mainProc(stdin.Text(), cand_num)
        fmt.Print(result)
    }
}
