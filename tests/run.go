package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/lemon-mint/godotenv"
)

var Benchmarks = []string{
	"BenchmarkParseRequest",
	"Benchmark_stricmp",
	"Benchmark_ContentLength_stricmp",
	"Benchmark_Net_URL_Parse",
	"Benchmark_H1_URI_Parse",
	"Benchmark_H1_URI_Query",
	"Benchmark_Request_Reader",
}

func main() {
	godotenv.Load()

	err := os.Mkdir("testOutput", 0777)
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("# Benchmark Results (OS: %s, Arch: %s, Go Version: %s)\n", runtime.GOOS, runtime.GOARCH, runtime.Version())

	for _, benchmark := range Benchmarks {
		var buffer bytes.Buffer
		cmd := exec.Command("go", "test", "-bench="+benchmark, "-benchmem", "-cpuprofile", "testOutput/"+benchmark+"_cpu_profile.out", "-memprofile", "testOutput/"+benchmark+"_mem_profile.out", "github.com/go-www/silverlining/h1")
		cmd.Stdout = &buffer
		cmd.Stderr = &buffer
		err := cmd.Run()
		if err != nil {
			panic(err)
		}

		var CPUImageBuffer bytes.Buffer
		cmd = exec.Command("go", "tool", "pprof", "-png", "./testOutput/"+benchmark+"_cpu_profile.out")
		cmd.Stdout = &CPUImageBuffer
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			panic(err)
		}

		var MemImageBuffer bytes.Buffer
		cmd = exec.Command("go", "tool", "pprof", "-png", "./testOutput/"+benchmark+"_mem_profile.out")
		cmd.Stdout = &MemImageBuffer
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			panic(err)
		}

		cpuURL, _, err := uploadImage(CPUImageBuffer.Bytes(), benchmark+"_cpu_profile.png")
		if err != nil {
			panic(err)
		}

		memURL, _, err := uploadImage(MemImageBuffer.Bytes(), benchmark+"_mem_profile.png")
		if err != nil {
			panic(err)
		}

		fmt.Printf("## Benchmark %s\n", benchmark)
		fmt.Printf("\n```\n%s\n```\n", buffer.String())
		fmt.Printf("\n### CPU Profile\n")
		fmt.Printf("\n![CPU Profile](%s)\n\n", cpuURL)
		fmt.Printf("\n### Memory Profile\n")
		fmt.Printf("\n![Memory Profile](%s)\n\n", memURL)
		fmt.Printf("\n")
	}

	fmt.Print("\n")

	err = exec.Command("go", "build", "-o", "./pprofserver.exe.out", "./pprofserver").Run()
	if err != nil {
		panic(err)
	}

	// Start the pprof server
	cmd := exec.Command("./pprofserver.exe.out")
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	defer cmd.Process.Kill()

	// Wait for the pprof server to start
	time.Sleep(time.Second * 10)

	downloadDone := make(chan struct{})
	traceDone := make(chan struct{})
	// Download CPU profile
	go func() {
		f, err := os.Create("./testOutput/server_cpu_profile.out")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		resp, err := http.Get("http://localhost:6060/debug/pprof/profile?seconds=30")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		_, err = io.Copy(f, resp.Body)
		if err != nil {
			panic(err)
		}

		downloadDone <- struct{}{}
	}()

	// Download trace
	go func() {
		f, err := os.Create("./testOutput/server_trace.out")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		resp, err := http.Get("http://localhost:6060/debug/pprof/trace?seconds=30")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		_, err = io.Copy(f, resp.Body)
		if err != nil {
			panic(err)
		}

		traceDone <- struct{}{}
	}()

	// Start the Load Test
	cmd = exec.Command("oha", "-z", "35sec", "--no-tui", "http://localhost:8080/plaintext")
	var ohaOutput bytes.Buffer
	cmd.Stdout = &ohaOutput
	cmd.Stderr = &ohaOutput
	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	<-downloadDone
	<-traceDone

	fmt.Print("\n## Load Test Results\n")
	fmt.Printf("\n```\n%s\n```\n", ohaOutput.String())
	fmt.Printf("\n")

	// Create Benchmark Results Image
	var CPUImageBuffer bytes.Buffer
	cmd = exec.Command("go", "tool", "pprof", "-png", "./testOutput/server_cpu_profile.out")
	cmd.Stdout = &CPUImageBuffer
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	cpuURL, _, err := uploadImage(CPUImageBuffer.Bytes(), "server_cpu_profile.png")
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n### CPU Profile\n")
	fmt.Printf("\n![CPU Profile](%s)\n\n", cpuURL)
}
