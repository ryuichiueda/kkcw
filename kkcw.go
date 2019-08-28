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

const VERSION = "0.2.0"

func help() {
    fmt.Fprintf(os.Stderr, "kkcw: kkc wrapper %s\n", VERSION)
    fmt.Fprintln(os.Stderr, "Copyright (C) 2019 Ryuichi Ueda.");
    fmt.Println()
    fmt.Fprintln(os.Stderr, "usage: echo <string> | kkcw");
    fmt.Println()
    fmt.Fprintln(os.Stderr, "Released under the GPLv3")
    fmt.Fprintln(os.Stderr, "https://github.com/ryuichiueda/kkcw")
}

type Token struct {
    token_id int
    results []string
}

func readline () string {
    stdin := bufio.NewScanner(os.Stdin)
    if stdin.Scan() {
        return stdin.Text()
    }else{
        return ""
    }
}

func parse(text string, token_id int) Token {
    slice := strings.Split(text, "\n")

    resultmask := regexp.MustCompile(`[^:]*: `)

    resultline := regexp.MustCompile(`\d+: *`)
    orgstrpart := regexp.MustCompile(`/[^>]+>`)

    r := Token{}
    r.token_id = token_id

    for _, str := range slice {
	str = strings.Replace(str, ">> ", "", -1)
        if ! resultline.MatchString(str) {
	   continue
        }

	str = resultmask.ReplaceAllString(str, "")
	str = orgstrpart.ReplaceAllString(str, "")
	str = strings.Replace(str, "<", "", -1)

	r.results = append(r.results, str)
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
    ans := strings.Replace(tmp2, "<", "&lt;", -1)

    return ans
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

    slice := strings.Split(str, " ")
    for n, str := range slice {
        str_clean := escape(str)
        out := kkc(str_clean + " " + strconv.Itoa(candidate_num))
	tokens = append(tokens, parse(out, n))
    }

    var ans string
    for n := 0; n < candidate_num; n++ {
        ans += concat(tokens, n)
    }
    return ans
}

func main() {
    var err error
    cand_num := 1

    switch len(os.Args) {
    case 2:
        help()
        os.Exit(0)
    case 3:
        if os.Args[1] != "-n" {
            fmt.Fprintln(os.Stderr, "Invalid Option")
	    os.Exit(1)
        }

        cand_num, err = strconv.Atoi(os.Args[2])
        if err != nil {
            fmt.Fprintln(os.Stderr, "Invalid Number")
	    os.Exit(1)
        }
    }

    result := mainProc(readline(), cand_num)
    fmt.Print(result)
}
