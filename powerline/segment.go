package powerline

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Segment struct {
	Bg     string
	Fg     string
	sepFg  string
	values []string
}

func isWritableDir(dir string) bool {
	tmpPath := path.Join(dir, ".powerline-write-test")
	_, err := os.Create(tmpPath)
	if err != nil {
		return false
	}
	os.Remove(tmpPath)
	return true
}

func LockSegment(cwd string, t Theme, s Symbols) Segment {
	if isWritableDir(cwd) {
		return Segment{values: nil}
	} else {
		return Segment{
			Bg:     t.Lock.Bg,
			Fg:     t.Lock.Fg,
			values: []string{s.Lock},
		}
	}
}

func GetCurrentWorkingDir() (string, []string) {
	dir, err := filepath.Abs(".")
	if err != nil {
		log.Fatal(err)
	}
	userDir := strings.Replace(dir, os.Getenv("HOME"), "~", 1)
	userDir = strings.TrimSuffix(userDir, "/")
	parts := strings.Split(userDir, "/")
	return dir, parts
}

func HomeSegment(cwdParts []string, t Theme) Segment {
	if cwdParts[0] == "" {
		cwdParts[0] = "/"
	}

	if cwdParts[0] == "~" {
		return Segment{
			Bg:     t.Path.Bg,
			Fg:     t.Path.Fg,
			values: []string{"~"},
		}
	} else {
		return Segment{
			Bg:     t.Path.Bg,
			Fg:     t.Path.Fg,
			values: []string{cwdParts[0]},
		}
	}
}

func PathSegment(cwdParts []string, t Theme, s Symbols) Segment {
	cwdParts = cwdParts[1:]

	length := len(cwdParts)
	if length > 3 {
		tmp := []string{}
		//tmp = append(tmp, cwdParts[0])
		tmp = append(tmp, s.Ellipsis)
		tmp = append(tmp, cwdParts[length-2])
		tmp = append(tmp, cwdParts[length-1])
		cwdParts = tmp
	}

	return Segment{
		Bg:     t.Path.Bg,
		Fg:     t.Path.Fg,
		sepFg:  t.Path.SepFg,
		values: cwdParts,
	}
}

func getGitInformation() (string, bool) {
	var status string
	var staged bool
	stdout, _ := exec.Command("git", "status", "--ignore-submodules").Output()
	reBranch := regexp.MustCompile(`^(HEAD detached at|HEAD detached from|On branch) (\S+)`)
	matchBranch := reBranch.FindStringSubmatch(string(stdout))
	if len(matchBranch) > 0 {
		if matchBranch[2] == "detached" {
			status = matchBranch[2]
		} else {
			status = matchBranch[2]
		}
	}

	reStatus := regexp.MustCompile(`Your branch is (ahead|behind).*?([0-9]+) comm`)
	matchStatus := reStatus.FindStringSubmatch(string(stdout))
	if len(matchStatus) > 0 {
		status = fmt.Sprintf("%s %s", status, matchStatus[2])
		if matchStatus[1] == "behind" {
			status = fmt.Sprintf("%s\u21E9", status)
		} else if matchStatus[1] == "ahead" {
			status = fmt.Sprintf("%s\u21E7", status)
		}
	}

	// handle diverged branches (1 behind, 1 ahead)
	reAmend := regexp.MustCompile(`have (\d+) and (\d+) different commits`)
	matchAmend := reAmend.FindStringSubmatch(string(stdout))
	if len(matchAmend) > 0 {
		status = fmt.Sprintf("%s %s\u21E9 %s\u21E7", status, matchAmend[1], matchAmend[2])
	}

	staged = !strings.Contains(string(stdout), "nothing to commit")
	if strings.Contains(string(stdout), "Untracked files") {
		status = fmt.Sprintf("%s \u271A", status)
	}

	return status, staged
}

func GitSegment(t Theme) Segment {
	gitStatus, gitStaged := getGitInformation()

	if gitStatus != "" {
		var bg string
		var fg string
		if gitStaged {
			bg = t.Git.Dirty.Bg
			fg = t.Git.Dirty.Fg
		} else {
			bg = t.Git.Clean.Bg
			fg = t.Git.Clean.Fg

		}
		return Segment{
			Bg:     bg,
			Fg:     fg,
			values: []string{gitStatus},
		}
	} else {
		return Segment{values: nil}
	}
}

func TimeSegment(t Theme) Segment {
	timestr := time.Now().Format("15:04:05")
	if os.Getenv("SSH_CLIENT") != "" {
		return Segment{
			Bg:     t.Remote.Bg,
			Fg:     t.Remote.Fg,
			values: []string{"\u262F " + timestr},
		}
	} else {
		return Segment{
			Bg:     t.Home.Bg,
			Fg:     t.Home.Fg,
			values: []string{"\u262F " + timestr},
		}
	}
}

func AuthSegment(t Theme) Segment {
	auth := os.Getenv("AUTH_ACCOUNT")
	if auth == "" {
		return Segment{values: nil}
	}

	return Segment{
		Bg:     t.Auth.Bg,
		Fg:     t.Auth.Fg,
		values: []string{"#" + auth},
	}
}

func KitchenSegment(t Theme) Segment {
	focus := os.Getenv("KITCHEN_FOCUS")
	if focus == "" {
		return Segment{values: nil}
	}

	return Segment{
		Bg:     t.Kitchen.Bg,
		Fg:     t.Kitchen.Fg,
		values: []string{"@" + focus},
	}
}

func ExitCodeSegment(code string, t Theme) Segment {
	i, err := strconv.Atoi(code)
	if err != nil || i == 0 {
		return Segment{values: nil}
	} else {
		return Segment{
			Bg:     t.Error.Bg,
			Fg:     t.Error.Fg,
			values: []string{code},
		}
	}
}

func LinefeedSegment(t Theme) Segment {
    return Segment{
        Bg:     t.ShellBg,
        Fg:     t.ShellBg,
        values: []string{"\n"},
    }
}
