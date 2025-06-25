package cli

import (
	"github.com/dubbersthehoser/internal/app"

)

func RegisterCommands(c *Commands) {
	c.Register("help", handlerHelp, "")
}

func handlerHelp(s *app.State, cmd *Commands) error {
	fmt.Printf("%s\n", Name)
	fmt.Println("COMMANDS:")
	for _, c := range s.Commands {
		fmt.Printf("* %s\t%s\n", c.Name, s.Commands.Help[c.Name])
	}
	return nil
}
