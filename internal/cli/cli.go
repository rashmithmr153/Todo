package cli

import (
	"fmt"
	"os"
	"strconv"
	"todo/internal/store"
	"todo/internal/tui"
)

func Usage() string {
	use := "Usage:todo [OPTIONS] [VALUE]\n"
	use += "           add     \"Task-title\"\n"
	use += "           list                   \n"
	use += "           done     <id>          \n"
	use += "           delete  <id>          \n"

	return use
}

func Handle(s *store.Store) error {
	args := os.Args[1:]

	if len(args) == 0 {
		return tui.Run(s)
	} else {

		//CLI part
		var err error
		switch args[0] {

		case "add":
			if len(args) < 2 {
				fmt.Println(Usage())
				return fmt.Errorf("missing \"Task-title\"")
			}
			err = s.Add(args[1])

		case "list":
			s.List()

		case "done":
			if len(args) < 2 {
				fmt.Println(Usage())
				return fmt.Errorf("missing <id>")
			}
			n, e := strconv.Atoi(args[1])
			if e != nil {
				return fmt.Errorf("invalid id: %s\n%s", args[1], Usage())
			}
			err = s.MarkDone(n)

		case "delete":
			if len(args) < 2 {
				fmt.Println(Usage())
				return fmt.Errorf("missing <id>")
			}
			n, e := strconv.Atoi(args[1])
			if e != nil {
				return fmt.Errorf("invalid id: %s\n%s", args[1], Usage())
			}
			err = s.Delete(n)

		default:
			fmt.Println(Usage())
			return fmt.Errorf("invalid command")

		}
		return err
	}
}
