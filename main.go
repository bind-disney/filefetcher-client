package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strings"

	"path/filepath"

	"github.com/bind-disney/filefetcher-core/cli"
)

const (
	defaultHost        = ""
	defaultPort        = 27518
	defaultDirectory   = "downloads"
	defaultPermissions = 0775
	commandPrompt      = "$ "
	logPrefix          = "Filefetcher Client: "
)

var (
	hostOption      string
	portOption      uint
	directoryOption string
	helpOption      bool
	remoteDirectory string
	commandsList    string
	commandsHelp    string
	bufferSize      int64
)

func init() {
	initFlags()
	initLogger()
	initCommandsHelp()
}

func main() {
	if flag.NFlag() == 1 && helpOption {
		flag.Usage()
		os.Exit(0)
	}

	if err := validateRemoteDirectory(); err != nil {
		cli.FatalError("File directory", err)
	}

	remoteAddress := fmt.Sprintf("%s:%d", hostOption, portOption)
	// We need to keep listener separately in more low-level way instead of rpc.Dial,
	// because of local tcp connection address is needed in order to make RPC requests
	listener, err := net.Dial("tcp", remoteAddress)
	if err != nil {
		cli.FatalError("Connection", err)
	}
	client := rpc.NewClient(listener)
	// Listener will be closed automatically by the client' Close() method
	defer client.Close()
	localAddress := listener.LocalAddr().String()

	if err = getRemoteDirectory(localAddress, client); err != nil {
		cli.FatalError("RPC", err)
	}

	if err = getRemoteBufferSize(localAddress, client); err != nil {
		cli.FatalError("RPC", err)
	}

	printCommandsList()
	printPrompt()

	var (
		line            string
		command         string
		targetDirectory string
	)

	scanner := bufio.NewScanner(os.Stdin)

REPL:

	for scanner.Scan() {
		line = strings.TrimSpace(scanner.Text())
		input := strings.Fields(line)
		inputSize := len(input)

		if inputSize == 0 {
			printPrompt()
			continue
		}

		command = strings.ToLower(input[0])
		err = nil // Clear previous error in the current loop iteration

		switch command {
		case "list", "ls":
			if inputSize >= 2 {
				targetDirectory = input[1]
			} else {
				targetDirectory = ""
			}

			err = listRemoteFiles(localAddress, client, targetDirectory)
		case "chdir", "cd":
			if inputSize < 2 {
				fmt.Println("New directory is not specified")
				printPrompt()
				continue
			}

			targetDirectory = input[1]
			err = changeRemoteDirectory(localAddress, client, targetDirectory)
		case "get", "download", "fetch":
			switch inputSize {
			case 1:
				fmt.Println("Filename is not specified")
				printPrompt()
				continue
			case 2:
				targetDirectory, err = filepath.Abs(directoryOption)
				if err != nil {
					cli.LogError("Download", err)
					printPrompt()
					continue
				}

				fileName := input[1]
				err = getRemoteFile(localAddress, client, fileName, targetDirectory)
			default:
				targetDirectory, err = filepath.Abs(input[2])
				if err != nil {
					cli.LogError("Download", err)
					printPrompt()
					continue
				}

				fileName := input[1]
				err = getRemoteFile(localAddress, client, fileName, targetDirectory)
			}
		case "help", "h":
			printCommandsHelp()
		case "exit", "quit", "q":
			break REPL
		default:
			printCommandsList()
		}

		if err == nil {
			printPrompt()
		}
	}
}

func initFlags() {
	flag.StringVar(&hostOption, "H", defaultHost, "Server hostname or IP address")
	flag.UintVar(&portOption, "P", defaultPort, "Server port")
	flag.StringVar(&directoryOption, "D", defaultDirectory, "Target directory")
	flag.BoolVar(&helpOption, "h", false, "Display this help")
	flag.Usage = cli.ShowUsage

	flag.Parse()
}

func initLogger() {
	log.SetOutput(os.Stderr)
	log.SetPrefix(logPrefix)
}

func printPrompt() {
	fmt.Printf("%s %s", remoteDirectory, commandPrompt)
}
