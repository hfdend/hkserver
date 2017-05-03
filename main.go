package main

import (
    "net/http"
    "log"
    "io/ioutil"
    "encoding/json"
    "git.125i.cn/hfdend/hkserver.git/global"
    _ "git.125i.cn/hfdend/hkserver.git/conf"
    "fmt"
    "os/exec"
    "os/user"
    "strings"
    "syscall"
    "strconv"
	"github.com/gogits/go-gogs-client"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if b, err := ioutil.ReadAll(r.Body); err != nil {
            log.Println(err)
        } else {
            parse(b)
        }
    })
    fmt.Println("listen: ", global.Config.Addr)
    if err := http.ListenAndServe(global.Config.Addr, nil); err != nil {
        log.Fatalln(err)
    }
}

func parse(b []byte) {
    var data gogs.PushPayload
    if err := json.Unmarshal(b, &data); err != nil {
        log.Println(err)
        return
    }
    do(data)
}

func do(data gogs.PushPayload) {
    if data.Repo == nil {
        return
    }
    if len(global.Config.Hooks) == 0 {
        return
    }
    for _, h := range global.Config.Hooks {
        if h.Resp == data.Repo.Name && h.Branch == data.Branch() {
            action(h.Commands)
        }
    }
}

func action(commands []global.Command) {
    for _, v := range commands {
        for _, s := range v.Args {
			ary := strings.Split(s, " ")
			name := ary[0]
			var args []string
			if len(ary) > 1 {
				args = ary[1:]
			}
			cmd := exec.Command(name, args...)
			if v.User != "" {
				u, err := user.Lookup(v.User)
				if err != nil {
					log.Println(err)
					return
				}
				uid, _ := strconv.Atoi(u.Uid)
				gid, _ := strconv.Atoi(u.Gid)
				cmd.SysProcAttr = &syscall.SysProcAttr{}
				cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)}
			}

			if v.Dir != "" {
				cmd.Dir = v.Dir
			}

			if v.Path != "" {
				cmd.Path = v.Path
			}

			stdout, err := cmd.StdoutPipe()
			if err != nil {
				log.Println(err)
				return
			}
			stderr, err := cmd.StderrPipe()
			if err != nil {
				log.Println(err)
				return
			}
			if err != nil {
				log.Println(err)
				return
			}
			if err := cmd.Start(); err != nil {
				log.Println(err)
				return
			}

			bytesErr, err := ioutil.ReadAll(stderr)
			if err != nil {
				log.Println(err)
				return
			}

			if len(bytesErr) != 0 {
				fmt.Println(string(bytesErr))
			}

			bytes, err := ioutil.ReadAll(stdout)
			if err != nil {
				log.Println(err)
				return
			}
			if len(bytes) > 0 {
				fmt.Println(string(bytes))
			}
			if err := cmd.Wait(); err != nil {
				log.Println(err)
				return
			}
		}
    }
}