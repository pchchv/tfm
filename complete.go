package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pchchv/golog"
)

var (
	genCmdWords = []string{
		"set",
		"map",
		"cmap",
		"cmd",
		"quit",
		"up",
		"half-up",
		"page-up",
		"scroll-up",
		"down",
		"half-down",
		"page-down",
		"scroll-down",
		"updir",
		"open",
		"jump-next",
		"jump-prev",
		"top",
		"bottom",
		"high",
		"middle",
		"low",
		"toggle",
		"invert",
		"unselect",
		"glob-select",
		"glob-unselect",
		"calcdirsize",
		"copy",
		"cut",
		"paste",
		"clear",
		"sync",
		"draw",
		"redraw",
		"load",
		"reload",
		"echo",
		"echomsg",
		"echoerr",
		"cd",
		"select",
		"delete",
		"rename",
		"source",
		"push",
		"read",
		"shell",
		"shell-pipe",
		"shell-wait",
		"shell-async",
		"find",
		"find-back",
		"find-next",
		"find-prev",
		"search",
		"search-back",
		"search-next",
		"search-prev",
		"filter",
		"setfilter",
		"mark-save",
		"mark-load",
		"mark-remove",
		"tag",
		"tag-toggle",
		"cmd-escape",
		"cmd-complete",
		"cmd-menu-complete",
		"cmd-menu-complete-back",
		"cmd-menu-accept",
		"cmd-enter",
		"cmd-interrupt",
		"cmd-history-next",
		"cmd-history-prev",
		"cmd-left",
		"cmd-right",
		"cmd-home",
		"cmd-end",
		"cmd-delete",
		"cmd-delete-back",
		"cmd-delete-home",
		"cmd-delete-end",
		"cmd-delete-unix-word",
		"cmd-yank",
		"cmd-transpose",
		"cmd-transpose-word",
		"cmd-word",
		"cmd-word-back",
		"cmd-delete-word",
		"cmd-capitalize-word",
		"cmd-uppercase-word",
		"cmd-lowercase-word",
	}

	genOptWords = []string{
		"anchorfind",
		"noanchorfind",
		"anchorfind!",
		"autoquit",
		"noautoquit",
		"autoquit!",
		"cursorfmt",
		"cursorpreviewfmt",
		"dircache",
		"nodircache",
		"dircache!",
		"dircounts",
		"nodircounts",
		"dircounts!",
		"dirfirst",
		"nodirfirst",
		"dirfirst!",
		"dironly",
		"nodironly",
		"dironly!",
		"dirpreviews",
		"nodirpreviews",
		"dirpreviews!",
		"drawbox",
		"nodrawbox",
		"drawbox!",
		"globsearch",
		"noglobsearch",
		"globsearch!",
		"hidden",
		"nohidden",
		"hidden!",
		"icons",
		"noicons",
		"icons!",
		"ignorecase",
		"noignorecase",
		"ignorecase!",
		"ignoredia",
		"noignoredia",
		"ignoredia!",
		"incsearch",
		"noincsearch",
		"incsearch!",
		"incfilter",
		"noincfilter",
		"incfilter!",
		"mouse",
		"nomouse",
		"mouse!",
		"number",
		"nonumber",
		"number!",
		"preview",
		"nopreview",
		"preview!",
		"relativenumber",
		"norelativenumber",
		"relativenumber!",
		"reverse",
		"noreverse",
		"reverse!",
		"smartcase",
		"nosmartcase",
		"smartcase!",
		"smartdia",
		"nosmartdia",
		"smartdia!",
		"waitmsg",
		"wrapscan",
		"nowrapscan",
		"wrapscan!",
		"wrapscroll",
		"nowrapscroll",
		"wrapscroll!",
		"findlen",
		"period",
		"scrolloff",
		"tabstop",
		"errorfmt",
		"filesep",
		"hiddenfiles",
		"history",
		"ifs",
		"info",
		"previewer",
		"cleaner",
		"promptfmt",
		"ratios",
		"selmode",
		"shell",
		"shellflag",
		"shellopts",
		"sortby",
		"timefmt",
		"tempmarks",
		"tagfmt",
		"infotimefmtnew",
		"infotimefmtold",
		"truncatechar",
	}
)

func matchLongest(s1, s2 []rune) []rune {
	i := 0
	for ; i < len(s1) && i < len(s2); i++ {
		if s1[i] != s2[i] {
			break
		}
	}
	return s1[:i]
}

func matchWord(s string, words []string) (matches []string, longest []rune) {
	for _, w := range words {
		if !strings.HasPrefix(w, s) {
			continue
		}

		matches = append(matches, w)

		if len(longest) != 0 {
			longest = matchLongest(longest, []rune(w))
		} else if s != "" {
			longest = []rune(w + " ")
		}
	}

	if len(longest) == 0 {
		longest = []rune(s)
	}

	return
}

func matchExec(s string) (matches []string, longest []rune) {
	var words []string

	paths := strings.Split(envPath, string(filepath.ListSeparator))

	for _, p := range paths {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			continue
		}

		files, err := ioutil.ReadDir(p)
		if err != nil {
			golog.Info("reading path: %s", err)
		}

		for _, f := range files {
			if !strings.HasPrefix(f.Name(), s) {
				continue
			}

			f, err = os.Stat(filepath.Join(p, f.Name()))
			if err != nil {
				golog.Info("getting file information: %s", err)
				continue
			}

			if !f.Mode().IsRegular() || !isExecutable(f) {
				continue
			}

			golog.Info(f.Name())
			words = append(words, f.Name())
		}
	}

	sort.Strings(words)

	if len(words) > 0 {
		uniq := words[:1]
		for i := 1; i < len(words); i++ {
			if words[i] != words[i-1] {
				uniq = append(uniq, words[i])
			}
		}
		words = uniq
	}

	return matchWord(s, words)
}

