package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
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
	err := os.Mkdir("testOutput", 0777)
	if err != nil {
		log.Println(err)
	}

	fmt.Print("# Benchmark Results (OS:", runtime.GOOS, ")\n")

	for _, benchmark := range Benchmarks {
		var buffer bytes.Buffer
		cmd := exec.Command("go", "test", "-bench="+benchmark, "-benchmem", "-cpuprofile", "testOutput/"+benchmark+"_profile.out", "-memprofile", "testOutput/"+benchmark+"_memprofile.out")
		cmd.Stdout = &buffer
		cmd.Stderr = &buffer
		err := cmd.Run()
		if err != nil {
			panic(err)
		}

		var CPUImageBuffer bytes.Buffer
		cmd = exec.Command("go", "tool", "pprof", "-png", "./testOutput/"+benchmark+"_profile.out")
		cmd.Stdout = &CPUImageBuffer
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			panic(err)
		}

		var MemImageBuffer bytes.Buffer
		cmd = exec.Command("go", "tool", "pprof", "-png", "./testOutput/"+benchmark+"_memprofile.out")
		cmd.Stdout = &MemImageBuffer
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			panic(err)
		}

		cpuURL, _, err := uploadImage(CPUImageBuffer.Bytes(), benchmark+"_profile.png")
		if err != nil {
			panic(err)
		}

		memURL, _, err := uploadImage(MemImageBuffer.Bytes(), benchmark+"_memprofile.png")
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
}
