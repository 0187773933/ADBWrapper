package utils

import (
	// "os"
	"os/exec"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"strings"
	"time"
	"syscall"
	"bytes"
)

func ExecProcess( bash_command string , arguments ...string ) ( result string ) {
	command := exec.Command( bash_command , arguments... )
	//command.Env = append( os.Environ() , "DISPLAY=:0.0" )
	out , err := command.Output()
	if err != nil {
		fmt.Println( bash_command )
		fmt.Println( arguments )
		fmt.Sprintf( "%s\n" , err )
	}
	result = string( out[:] )
	return
}

// https://github.com/dohzya/timeout/blob/master/main.go
func ExecProcessWithTimeout( timeout_duration time.Duration , bash_command string , arguments ...string ) ( result string ) {
	command := exec.Command( bash_command , arguments... )
	// command.Stdin = os.Stdin
	// command.Stdout = os.Stdout
	// command.Stderr = os.Stderr
	var outb , errb bytes.Buffer
	command.Stdout = &outb
	command.Stderr = &errb
	command.Start()
	// if err != nil {
	// 	fmt.Println( bash_command )
	// 	fmt.Println( arguments )
	// 	fmt.Sprintf( "%s\n" , err )
	// }
	time.AfterFunc( timeout_duration , func() {
		command.Process.Signal( syscall.SIGTERM )
	})
	command.Wait()
	// status := command.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
	// fmt.Println( "out:" , outb.String() , "err:" , errb.String() )
	result = outb.String()
	return
}

// https://stackoverflow.com/questions/13578416/read-binary-stdout-data-from-adb-shell/31401447#31401447
func ExecProcessWithTimeoutGetBytes( timeout_duration time.Duration , bash_command string , arguments ...string ) ( result []byte ) {
	command := exec.Command( bash_command , arguments... )
	// command.Stdin = os.Stdin
	// command.Stdout = os.Stdout
	// command.Stderr = os.Stderr
	var outb , errb bytes.Buffer
	command.Stdout = &outb
	command.Stderr = &errb
	command.Start()
	// if err != nil {
	// 	fmt.Println( bash_command )
	// 	fmt.Println( arguments )
	// 	fmt.Sprintf( "%s\n" , err )
	// }
	time.AfterFunc( timeout_duration , func() {
		command.Process.Signal( syscall.SIGTERM )
	})
	command.Wait()
	// status := command.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
	// fmt.Println( "out:" , outb.String() , "err:" , errb.String() )
	result = outb.Bytes()
	return
}


func ExecProcessWithTimeoutAndGetOutputLines( timeout_duration time.Duration , bash_command string , arguments ...string ) ( lines []string ) {
	command := exec.Command( bash_command , arguments... )
	// command.Stdin = os.Stdin
	// command.Stdout = os.Stdout
	// command.Stderr = os.Stderr
	var outb , errb bytes.Buffer
	command.Stdout = &outb
	command.Stderr = &errb
	command.Start()
	// if err != nil {
	// 	fmt.Println( bash_command )
	// 	fmt.Println( arguments )
	// 	fmt.Sprintf( "%s\n" , err )
	// }
	time.AfterFunc( timeout_duration , func() {
		command.Process.Signal( syscall.SIGTERM )
	})
	command.Wait()
	// status := command.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
	// fmt.Println( "out:" , outb.String() , "err:" , errb.String() )
	result := outb.String()
	fmt.Println( result )
	non_empty_lines := strings.Replace( result , "\n\n" , "\n" , -1 )
	lines = strings.Split( non_empty_lines , "\n" )
	return
}

func ExecProcessAndGetOutputLines( bash_command string , arguments ...string ) ( lines []string ) {
	command := exec.Command( bash_command , arguments... )
	out , _ := command.Output()
	result := string( out[:] )
	non_empty_lines := strings.Replace( result , "\n\n" , "\n" , -1 )
	lines = strings.Split( non_empty_lines , "\n" )
	return
}


func WriteJSON( filePath string , data interface{} ) {
	file, _ := json.MarshalIndent( data , "" , " " )
	_ = ioutil.WriteFile( filePath , file , 0644 )
}

func RemoveEmpties( list []string ) ( result []string ) {
	for _ , item := range list {
		if item != "" {
			result = append( result , item )
		}
	}
	return
}