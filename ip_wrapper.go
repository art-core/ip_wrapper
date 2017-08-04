package main

import (
	"bytes"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"log"
	"syscall"
)

var verbose bool
var hostname string
var tail []string
var version string
var printHelp bool
var printVersion bool

func initArgs() {
	flag.BoolVar(&printHelp, "help", false, "print help and exit")
	flag.BoolVar(&printVersion, "V", false, "print version and exit")
	flag.BoolVar(&verbose, "v", false, "be verbose")
	flag.StringVar(&hostname, "host", "", "the hostname to check")

	flag.Parse()

	if printHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if printVersion {
		fmt.Println("ip_wrapper Version:", version)
		os.Exit(0)
	}


	// this is the binary to run, if the time's right
	tail = flag.Args()
	if len(tail) < 1 {
		os.Exit(3)
	}
}

func main() {
	initArgs()

	if verbose {
		fmt.Println("host: ", hostname)
		fmt.Println("tail: ", tail)
	}

	ip_addrs, err := net.LookupHost(hostname)
	if err != nil {
		log.Fatal(err)
	}

	exitCode := 0
	total := 0
	exitCodes := []int{}
	for _, ip := range ip_addrs {
		exitCode = run(ip, tail)
		exitCodes = append(exitCodes, exitCode)
		total += exitCode
	}
	exitWithProperCode(total, exitCodes)
}

func run(ip string, tail []string) int {
	// get the real path of the executable
	if verbose {
		fmt.Printf("Checking IP: %s\n", ip)
	}
	binary, lookErr := exec.LookPath(tail[0])
	if lookErr != nil {
		fmt.Println(lookErr)
		// check not found
		return 3
	}

	// replace '%%IP%%' placeholder in the command to be executed
	reg, err := regexp.Compile("%%IP%%")
	if err != nil {
		log.Fatal(err)
		return 2
	}

	for k, v := range tail {
		tail[k] = reg.ReplaceAllString(v, ip)
	}

	// execute a sensu check, we need
	// - return value of the executed check
	// - stdout + stderr for debugging
	cmd := exec.Command(binary, tail...)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err != nil {
		log.Fatal(err)
		return 2
	}
	execErr := cmd.Start()
	if err != nil {
		log.Fatal(err)
		return 2
	}
	execErr = cmd.Wait()

	// the exit code was zero
	if execErr != nil {
		exitCode2 := 0
		if exitError, ok := execErr.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode2 = ws.ExitStatus()
			if verbose {
				fmt.Printf("Found exit code: %d\n", exitCode2)
			}
		} else {
			exitCode2 = 2
		}
		return exitCode2
	} else {
		// exit code was 0
		return 0
	}
}

func exitWithProperCode(total int, exitCodes []int) {
	msg := "Check"
	status := ""
	rest := ""
	finalExitCode := 0
	if total == 0 {
		status = "OK"
		rest = "OK"
	} else {
		if hasOnlyUnknowns(exitCodes) {
			status = "UNKNOWN"
			rest = "Check not found"
			finalExitCode = 3
		} else {
			if verbose {
				fmt.Printf("All exitCodes:")
				for _, ex := range exitCodes {
					fmt.Printf(" %d", ex)
				}
				fmt.Printf("\n")
			}
			status = "CRITICAL"
			rest = fmt.Sprintf("%d failed", len(exitCodes))
			finalExitCode = 2
		}
	}
	fmt.Printf("%s %s: %s\n", msg, status, rest)
	os.Exit(finalExitCode)
}

func hasOnlyUnknowns(exitCodes []int) bool {
	for _, v := range exitCodes {
		if v != 3 {
			return false
		}
	}
	return true
}
