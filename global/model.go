package global

type Hook struct {
    Resp     string
    Branch   string
    Commands []Command
}

type Command struct {
    User string
    Dir  string
    Path string
    Args []string
}

type Model struct {
    Addr  string
    Hooks []Hook
}

var (
    Config Model
    ConfigFile  *string
)
