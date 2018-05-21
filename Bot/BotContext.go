package bot

type Context struct {
    Admins  map[string]bool
    Waiting chan bool
}

func (this* Context) IsAdmin(UserId string) (bool) {
    flag, ok := this.Admins[UserId]
    if ok == false {
        return false
    }

    return flag
}