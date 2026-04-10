package tui

import (
	"fmt"
	"os"
	"todo/internal/store"
	"todo/internal/todo"

	"golang.org/x/term"
)

func Run(s *store.Store) error {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("Canot open terminal in raw mode: %s\n", err)
	}

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return err
	}

	defer term.Restore(int(os.Stdin.Fd()), oldState)

	fmt.Print("\x1b[?1049h")
	defer fmt.Print("\x1b[?1049l")

	buf := make([]byte, 3)
	selected := 0
	drawSplitLayout(s.Todos, selected, width, height)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			return fmt.Errorf("Error while reading input :%s\n", err)
		}

		if n == 1 {
			switch buf[0] {
			case 'q':
				return nil
			case 'd':
				if len(s.Todos) > 0 {
					s.Delete(s.Todos[selected].Id)
					if len(s.Todos) == 0 {
						selected = 0
					} else if selected >= len(s.Todos) {
						selected = len(s.Todos) - 1
					}
				}
			case 'a':
				title := readInput()
				if title != "" {
					s.Add(title)
				}
			case '\r', '\n':
				if len(s.Todos) > 0 {
					s.MarkDone(s.Todos[selected].Id)
				}
			}
		} else if n == 3 && buf[0] == '\x1b' && buf[1] == '[' {
			switch buf[2] {
			case 'A':
				if selected > 0 {
					selected -= 1
				}
			case 'B':
				if selected < len(s.Todos)-1 {
					selected += 1
				}
			}
		}

		drawSplitLayout(s.Todos, selected, width, height)
	}
}

func readInput() string {
	fmt.Print("\x1b[2J\x1b[H")
	fmt.Print("Add new todo:\r\n")
	fmt.Print("> ")

	input := []rune{}
	buf := make([]byte, 1)

	for {
		_, err := os.Stdin.Read(buf)
		if err != nil {
			return ""
		}

		ch := buf[0]

		switch ch {
		case '\r', '\n':
			return string(input)

		case '\x1b':
			return ""
		case 127, '\b':
			if len(input) != 0 {
				input = input[:len(input)-1]
				fmt.Print("\b \b")
			}

		default:
			input = append(input, rune(ch))
			fmt.Printf("%c", ch)
		}
	}
}

func drawMenu(todos []todo.Todo, selected, width, height int) {
	// Clear screen and move to top
	fmt.Print("\x1b[2J\x1b[H")
	fmt.Printf("Terminal size is %dx%d\r\n", width, height)

	fmt.Print("TODO List (↑/↓ to navigate, Enter to toggle, d to delete, q to quit)\r\n\r\n")

	if len(todos) == 0 {
		fmt.Print("No todos yet. Press 'a' to add one.\r\n")
		return
	}

	for i, t := range todos {

		status := "[ ]"
		if t.Done {
			status = "[✓]"
		}

		if i == selected {
			fmt.Print("\x1b[7m")
		}

		fmt.Printf("%s %d. %s\r\n", status, t.Id, t.Title)

		fmt.Print("\x1b[0m")
	}
}

func drawSplitLayout(todos []todo.Todo, selected, width, height int) {
	fmt.Print("\x1b[2J\x1b[H")

	// panle widths
	leftWidth := width / 2
	// rightWidth := width - leftWidth

	fmt.Printf("TODO lists%s\r\n", spcaes(leftWidth-9))

	for col := 1; col <= width; col++ {
		fmt.Printf("\x1b[2;%dH_", col)
	}

	if len(todos) == 0 {
		fmt.Printf("\x1b[3;3H")
		fmt.Print("No todos yet. Press 'a' to add one.")
		// Still draw the divider and status bar
		for row := 1; row <= height-1; row++ {
			fmt.Printf("\x1b[%d;%dH│", row, leftWidth+1)
		}
		fmt.Printf("\x1b[%d;1H", height)
		fmt.Print("q:quit a:add")
		return
	}

	maxRows := height - 3
	for i := 0; i < maxRows && i < len(todos); i++ {
		fmt.Printf("\x1b[%d;1H", i+3)
		status := "[ ]"
		if todos[i].Done {
			status = "[✓]"
		}

		if i == selected {
			fmt.Print("\x1b[7m")
		}
		titileMaxwidth := leftWidth - 10
		title := todos[i].Title
		if len(title) > titileMaxwidth {
			title = title[:titileMaxwidth-3] + "..."
		}
		fmt.Printf("%s %d. %s", status, todos[i].Id, title)

		fmt.Print("\x1b[0m")
	}
	fmt.Printf("\x1b[1;%dH", leftWidth+2) // Row 1, right side
	fmt.Print("DETAILS")

	if len(todos) > 0 && selected < len(todos) {
		t := todos[selected]

		fmt.Printf("\x1b[3;%dH", leftWidth+2)
		fmt.Printf("Title: %s", t.Title)

		fmt.Printf("\x1b[4;%dH", leftWidth+2)
		if t.Done {
			fmt.Print("Status: Done")
		} else {
			fmt.Print("Status: Not done")
		}

		fmt.Printf("\x1b[5;%dH", leftWidth+2)
		fmt.Printf("Created: %s", t.CreatedAt.Format("02/01/2006"))
	}

	//draws line inbetween
	for row := 1; row <= height-1; row++ {
		fmt.Printf("\x1b[%d;%dH│", row, leftWidth+1)
	}
	fmt.Printf("\x1b[%d;1H", height)
	fmt.Print("q:quit ↑↓:navigate Enter:toggle d:delete a:add")

}

func spcaes(n int) string {
	if n <= 0 {
		return ""
	}
	s := make([]byte, n)
	for i := range s {
		s[i] = ' '
	}
	return string(s)
}
