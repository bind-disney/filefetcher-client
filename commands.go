package main

import (
	"fmt"
	"net/rpc"
	"os"
	"path"
	"path/filepath"

	"io"

	"net"

	"github.com/bind-disney/filefetcher-core/cli"
	fileRpc "github.com/bind-disney/filefetcher-core/rpc"
)

type (
	commandAliases []string

	commandReference struct {
		aliases commandAliases
		help    string
	}
)

const (
	rpcListFilesMethod        = "Server.ListFiles"
	rpcGetFileMethod          = "Server.GetFile"
	rpcCurrentDirectoryMethod = "Server.CurrentDirectory"
	rpcChangeDirectoryMethod  = "Server.ChangeDirectory"
	rpcBufferSizeMethod       = "Server.BufferSize"
)

func listRemoteFiles(address string, client *rpc.Client, targetDirectory string) error {
	method := rpcListFilesMethod
	files := make([]string, 10)
	request := fileRpc.NewFileSystemRequest(address, targetDirectory)
	if err := client.Call(method, request, &files); err != nil {
		handleCommandError(method, err)
		return err
	}

	if len(files) == 0 {
		fmt.Println("No files in this directory")
	} else {
		for _, entry := range files {
			fmt.Println(entry)
		}
	}

	return nil
}

func getRemoteFile(address string, client *rpc.Client, fileName, targetDirectory string) error {
	var (
		receivedBytes    int64
		downloadResponse fileRpc.DownloadResponse
	)

	method := rpcGetFileMethod
	request := fileRpc.NewFileSystemRequest(address, fileName)
	if err := client.Call(method, request, &downloadResponse); err != nil || downloadResponse.Address == "" {
		handleCommandError(method, err)
		return err
	}

	connection, err := net.Dial("tcp", downloadResponse.Address)
	if err != nil {
		handleCommandError(method, err)
		return err
	}
	defer connection.Close()

	fileSize := downloadResponse.FileSize
	filePath := filepath.Join(targetDirectory, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		handleCommandError(method, err)
		return err
	}
	defer file.Close()

	for {
		if (fileSize - receivedBytes) < bufferSize {
			if _, err = io.CopyN(file, connection, (fileSize - receivedBytes)); err != nil {
				handleCommandError(method, err)
				return err
			}

			connection.Read(make([]byte, (receivedBytes+bufferSize)-fileSize))
			break
		}

		if _, err = io.CopyN(file, connection, bufferSize); err != nil {
			handleCommandError(method, err)
			return err
		}

		receivedBytes += bufferSize
	}

	return nil
}

func getRemoteDirectory(address string, client *rpc.Client) error {
	method := rpcCurrentDirectoryMethod
	request := fileRpc.NewClientRequest(address)
	if err := client.Call(method, request, &remoteDirectory); err != nil {
		handleCommandError(method, err)
		return err
	}

	return nil
}

func changeRemoteDirectory(address string, client *rpc.Client, targetDirectory string) error {
	method := rpcChangeDirectoryMethod
	request := fileRpc.NewFileSystemRequest(address, targetDirectory)
	if err := client.Call(method, request, &remoteDirectory); err != nil {
		handleCommandError(method, err)
		return err
	}

	fmt.Println(remoteDirectory)
	return nil
}

func validateRemoteDirectory() error {
	directory, err := filepath.Abs(directoryOption)
	if err != nil {
		return err
	}

	directory = path.Clean(directory)
	fileInfo, err := os.Stat(directory)

	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if os.IsNotExist(err) {
		if err = os.Mkdir(directory, os.ModeDir|defaultPermissions); err != nil {
			return err
		}

		// Re-fetch info for the newly created directory' further processing
		fileInfo, err = os.Stat(directory)
		if err != nil {
			return err
		}
	}

	if !fileInfo.IsDir() {
		return fmt.Errorf("%s is not a directory", directory)
	}

	return nil
}

func getRemoteBufferSize(address string, client *rpc.Client) error {
	method := rpcBufferSizeMethod
	request := fileRpc.NewClientRequest(address)
	if err := client.Call(method, request, &bufferSize); err != nil {
		handleCommandError(method, err)
		return err
	}

	return nil
}

func handleCommandError(command string, err error) {
	prefix := cli.FormatCommandPrefix(command)
	cli.LogError(prefix, err)
	printPrompt()
}
