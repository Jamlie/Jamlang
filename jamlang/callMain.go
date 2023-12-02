package jamlang

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Jamlie/Jamlang/parser"
	"github.com/Jamlie/Jamlang/runtimelang"
)

func Repl(env *runtimelang.Environment) {
	parser := parser.NewParser()
	fmt.Println("Repl mode. Type 'exit' to exit.")
	for {
		fmt.Print("> ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		text := scanner.Text()
		if text == "exit" {
			break
		}
		if text == "" {
			continue
		}

		program := parser.ProduceAST(text)
		runtimeValue, err := runtimelang.Evaluate(&program, *env)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(0)
		}
		fmt.Println(runtimeValue.Get())
	}
}

func CallMain(env *runtimelang.Environment) {
	runFlag := flag.Bool("r", false, "Run a file")
	helpFlag := flag.Bool("h", false, "Help")
	installFlag := flag.Bool("i", false, "Install a package")
	flag.Parse()
	args := flag.Args()
	if flag.NFlag() == 0 && len(args) == 0 {
		Repl(env)
	} else {
		if *runFlag {
			if len(args) < 1 {
				fmt.Println("No file specified")
				os.Exit(0)
			}

			if !strings.HasSuffix(args[0], ".jam") {
				fmt.Println("File must have .jam extension")
				os.Exit(0)
			}

			data, err := os.ReadFile(args[0])
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}

			program := parser.NewParser().ProduceAST(string(data))
			_, err = runtimelang.Evaluate(&program, *env)
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
		} else if *helpFlag {
			fmt.Println("Usage: jamlang [options] [file]")
			fmt.Println("Options:")
			fmt.Println("  -r\t\tRun a file")
			fmt.Println("  -h\t\tShow this help message")
			fmt.Println("  -i\t\tInstall a package")
			fmt.Println("  run\t\tRun a file")
			fmt.Println("  help\t\tShow this help message")
			fmt.Println("  install\tInstall a package")
		} else if *installFlag {
			if len(args) < 1 {
				fmt.Println("No package specified")
				os.Exit(0)
			}
			if _, err := os.Stat("std/" + args[0] + ".jam"); err == nil {
				fmt.Println("Library already installed")
				os.Exit(0)
			}

			if args[0] == "math" {
				getLibrary("math")
			} else if args[0] == "random" {
				getLibrary("random")
			} else if args[0] == "algorithm" {
				getLibrary("algorithm")
			} else if args[0] == "linkedlist" {
				getLibrary("linkedlist")
			} else {
				fmt.Println("Unknown library")
				os.Exit(0)
			}
		} else {
			option := args[0]
			if option == "run" {
				if len(args) < 2 {
					fmt.Println("No file specified")
					os.Exit(0)
				}

				if !strings.HasSuffix(args[1], ".jam") {
					fmt.Println("File must have .jam extension")
					os.Exit(0)
				}

				data, err := os.ReadFile(args[1])
				if err != nil {
					fmt.Println(err)
					os.Exit(0)
				}

				parser := parser.NewParser()
				program := parser.ProduceAST(string(data))

				_, err = runtimelang.Evaluate(&program, *env)
				if err != nil {
					fmt.Println(err)
					os.Exit(0)
				}
			} else if args[0] == "help" {
				fmt.Println("Usage: jamlang [options] [file]")
				fmt.Println("Options:")
				fmt.Println("  -r\t\tRun a file")
				fmt.Println("  -h\t\tShow this help message")
				fmt.Println("  -i\t\tInstall a library")
				fmt.Println("  run\t\tRun a file")
				fmt.Println("  help\t\tShow this help message")
				fmt.Println("  install\tInstall a library")
			} else if args[0] == "install" {
				if len(args) < 2 {
					fmt.Println("Usage: jamlang install [library]")
					os.Exit(0)
				}

				if _, err := os.Stat("std/" + args[1] + ".jam"); err == nil {
					fmt.Println("Library already installed")
					os.Exit(0)
				}

				if args[1] == "math" {
					getLibrary("math")
				} else if args[1] == "random" {
					getLibrary("random")
				} else if args[1] == "algorithm" {
					getLibrary("algorithm")
				} else if args[1] == "linkedlist" {
					getLibrary("linkedlist")
				} else {
					fmt.Println("Unknown library")
					os.Exit(0)
				}
			} else {
				fmt.Println("Unknown option")
				fmt.Println("Usage: jamlang [run] [file]")
				fmt.Println("Usage: jamlang -r [file]")
				fmt.Println("Usage: jamlang -i [library]")
				fmt.Println("Usage: jamlang install [library]")
				fmt.Println("Usage: jamlang -h")
				fmt.Println("Usage: jamlang help")
				fmt.Println("Usage: jamlang")
			}
		}
	}
}

func getLibrary(name string) {
	link := fmt.Sprintf("https://raw.githubusercontent.com/Jamlie/Jamlang/main/std/%s.jam", name)
	resp, err := http.Get(link)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer resp.Body.Close()

	if _, err := os.Stat("std"); os.IsNotExist(err) {
		os.Mkdir("std", 0755)
	}

	file, err := os.Create("std/" + name + ".jam")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
