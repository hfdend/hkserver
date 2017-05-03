package conf

import (
    "git.125i.cn/hfdend/hkserver.git/global"
    "flag"
    "os"
    "log"
    "io/ioutil"
    "gopkg.in/yaml.v2"
)

func init() {
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
    global.ConfigFile = flag.String("config-file", "./config/config.yml", "config file")
    flag.Parse()
    f, err := os.OpenFile(*global.ConfigFile, os.O_RDONLY, 444)
    if err != nil {
        log.Fatalln(err)
    }
    b, err := ioutil.ReadAll(f)
    if err != nil {
        log.Fatalln(err)
    }
    err = yaml.Unmarshal(b, &global.Config)
    if err != nil {
        log.Fatalln(err)
    }
}
