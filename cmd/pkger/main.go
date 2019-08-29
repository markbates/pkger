package main

import (
	"log"
	"os"

	"github.com/markbates/pkger/cmd/pkger/cmd"
)

// type command interface {
// 	Name() string
// 	execer
// 	flagger
// }
//
// type execer interface {
// 	Exec([]string) error
// }
//
// type flagger interface {
// 	Flags() *flag.FlagSet
// }
// type arrayFlags []string
//
// func (i arrayFlags) String() string {
// 	return fmt.Sprintf("%s", []string(i))
// }
//
// func (i *arrayFlags) Set(value string) error {
// 	*i = append(*i, value)
// 	return nil
// }

func main() {
	args := os.Args[1:]
	if err := cmd.Main(args, cmd.NewOptions(cmd.StdIO())); err != nil {
		log.Fatal(err)
	}
	//
	// defer func() {
	// 	c := exec.Command("go", "mod", "tidy")
	// 	c.Run()
	// }()
	//
	// root := &packCmd{
	// 	subCmds: []command{
	// 		&readCmd{}, &serveCmd{}, &infoCmd{},
	// 	},
	// }
	// if err := root.Exec(os.Args[1:]); err != nil {
	// 	log.Fatal(err)
	// }
}

// does not computee
