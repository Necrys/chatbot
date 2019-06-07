package commands

import "../CmdProcessor"
import "strings"
import "math/rand"
import "time"
import "fmt"

type CmdRoll struct {
    rng *rand.Rand
}

func NewCmdRoll() (*CmdRoll) {
    this := &CmdRoll { rng: rand.New(rand.NewSource(time.Now().UnixNano())) }
    return this
}

func (this* CmdRoll) HandleCommand(cmdCtx cmdprocessor.CommandCtxIf) (bool) {
    tokens := strings.Split(cmdCtx.Args(), ",")
    for i, _ := range tokens {
        tokens[i] = strings.Trim(tokens[i], " \n\t")
    }

    cmdCtx.Reply(fmt.Sprintf("%s", tokens[this.rng.Intn(len(tokens))]))

    return true
}
