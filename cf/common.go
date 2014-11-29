package cf

const (
	WAITING = iota //waiting for command
	INPUT          //performing a command and waiting for input from user
	COMMAND        //performing a command, no input from user is expected
	ADMIN          //performing command by admin, no input from user
)

type StatusType struct {
	Job    int
	Detail string
}
