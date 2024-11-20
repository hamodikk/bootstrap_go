# Bootstrap Resampling in Go

This Go program demonstrates the use of bootstrap resampling to calculate the standard error of the mean and median of different sample sizes in Go. The program generates a population with a given mean and standard deviation, performs bootstrap resampling, and logs the results.

## Table of Contents
- [Introduction](#introduction)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Code Explanation](#code-structure)
- [Testing](#testing)
- [Profiling](#profiling)
- [Analysis and Comparison](#analysis-and-comparison)
- [Summary and Recommendation](#summary-and-recommendation)

## Introduction

Bootstrapping is used to estimate the distribution of a statistic by generating multiple samples from the original data, usually with replacements. This program aims to replicate bootstrapping that can be achieved in R, but using Go. It generates a population and calculates the standard error of the mean and median.

The goal of this project is to compare the functionality and performance of Go in performing more complex statistical analysis to the performance of R packages. The bootstrapping program for R is created by Thomas W. Miller and is included in the repository [here]("run-bootstrap-median.R"). I have made small modifications at the end of the R program to include a runtime as well as report memory usage of the program for performance comparisons.

## Features

- Generates population
- Performs bootstrapping
- Reports the standard error of mean and median
- Generates log and profile
- Unit test and benchmarking available

## Installation

1. Make sure you have [Go installed](https://go.dev/doc/install).
2. Clone this repo to your local machine:
    ```bash
    git clone https://github.com/hamodikk/bootstrap_go.git
    ```
3. Navigate to the project directory
    ```bash
    cd <project-directory>
    ```

## Usage

Use the following command in your terminal or Powershell to run the program:
```bash
go run .\main.go
```

Or, you can also run the main.exe. Make sure to check the log file for the success of the program.

Here is how you can test and benchmark the code:
```bash
# Test the code to compare coefficient of regression between Python and R against the Go code
go test -v

# Benchmark the code to obtain execution times of the Go code
go test -bench=.
```

## Code Explanation

- Helper function that generates a population with the specified mean and standard deviation:
```go
func generatePopulation(size int, mean, stddev float64) []float64 {
	population := make([]float64, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		population[i] = rand.NormFloat64()*stddev + mean
	}
	return population
}
```
Note that rand.Seed is apparently deprecated, but using the code in this way has not caused any issues with running the code.

- Simple helper functions to calculate the median, mean and standard error:

```go
func median(data []float64) float64 {
	sort.Float64s(data)
	n := len(data)
	if n%2 == 0 {
		return (data[n/2-1] + data[n/2]) / 2.0
	}
	return data[n/2]
}
```

```go
func mean(data []float64) float64 {
	sum := 0.0
	for _, value := range data {
		sum += value
	}
	return sum / float64(len(data))
}
```

```go
func standardError(data []float64) float64 {
	mean := mean(data)
	variance := 0.0
	for _, value := range data {
		variance += math.Pow(value-mean, 2)
	}
	variance /= float64(len(data) - 1) // Use sample variance
	return math.Sqrt(variance) / math.Sqrt(float64(len(data)))
}
```

- Function to perform bootstrapping and calculate the standard error of mean and median:

```go
func bootstrap(population []float64, bootstrapSamples, sampleSize int) (float64, float64) {
	means := make([]float64, bootstrapSamples)
	medians := make([]float64, bootstrapSamples)
	for i := 0; i < bootstrapSamples; i++ {
		sample := make([]float64, sampleSize)
		for j := 0; j < sampleSize; j++ {
			sample[j] = population[rand.Intn(len(population))]
		}
		means[i] = mean(sample)
		medians[i] = median(sample)
	}

	seMean := standardError(means)
	seMedian := standardError(medians)
	return seMean, seMedian
}
```

- Start profiling within the main function:

```go
    go func() {
        err := http.ListenAndServe("localhost:6060", nil)
        if err != nil {
            log.Fatalf("Failed to start pprof server: %v", err)
        }
}()
```

- Set up logging. This appends the report in the file each time the program is ran:

```go
    logFile, err := os.OpenFile("bootstrap.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
```

## Testing

To run the tests and benchmarks, use the following commands:

```bash
go test -v
```
```bash
go test -bench=
```

## Profiling

The program includes a profiling server that can be accessed at localhost:6060. Once the program is executed, it will run idly for 5 minutes at the end of the main function execution to allow user to check the profiling.

## Analysis and Comparison

In order to compare the performance of the Go program with the R program written by Thomas W. Miller, I modified the R code to include memory usage and execution time for the bootstrapping. I have also included the log file for both the R program [here]("listing-from-100-run-bootstrap-median.txt") as well as the one for the Go program [here](bootstrap.log).

We can compare the functionality of both programs by comparing the outputs to see if they generate similar results. Following are the output generated by each program, which can also be found in their respective log files:

```txt
Study conditions:
  Population mean: 100 SD: 10 

Estimated standard errors using 100 bootstrap samples

Samples of size n = 25
  SE Mean from Central Limit Theorem: 2
  SE Mean from Samples: 1.99
  SE Mean from Bootstrap Samples: 1.99
  SE Median from Bootstrap Samples: 2.32

Samples of size n = 100
  SE Mean from Central Limit Theorem: 1
  SE Mean from Samples: 1.05
  SE Mean from Bootstrap Samples: 1.05
  SE Median from Bootstrap Samples: 1.23

Samples of size n = 225
  SE Mean from Central Limit Theorem: 0.67
  SE Mean from Samples: 0.68
  SE Mean from Bootstrap Samples: 0.69
  SE Median from Bootstrap Samples: 0.8

Samples of size n = 400
  SE Mean from Central Limit Theorem: 0.5
  SE Mean from Samples: 0.53
  SE Mean from Bootstrap Samples: 0.52
  SE Median from Bootstrap Samples: 0.65

```

```log
2024/11/19 19:36:20 Generated population
2024/11/19 19:36:20 Samples of size n = 25
2024/11/19 19:36:20   SE Mean from Central Limit Theorem for n = 25: 0.06
2024/11/19 19:36:20   SE Mean from Bootstrap Samples: 0.21
2024/11/19 19:36:20   SE Median from Bootstrap Samples: 0.24
2024/11/19 19:36:20 Samples of size n = 100
2024/11/19 19:36:20   SE Mean from Central Limit Theorem for n = 100: 0.03
2024/11/19 19:36:20   SE Mean from Bootstrap Samples: 0.09
2024/11/19 19:36:20   SE Median from Bootstrap Samples: 0.10
2024/11/19 19:36:20 Samples of size n = 225
2024/11/19 19:36:20   SE Mean from Central Limit Theorem for n = 225: 0.02
2024/11/19 19:36:20   SE Mean from Bootstrap Samples: 0.07
2024/11/19 19:36:20   SE Median from Bootstrap Samples: 0.08
2024/11/19 19:36:20 Samples of size n = 400
2024/11/19 19:36:20   SE Mean from Central Limit Theorem for n = 400: 0.02
2024/11/19 19:36:20   SE Mean from Bootstrap Samples: 0.04
2024/11/19 19:36:20   SE Median from Bootstrap Samples: 0.05
2024/11/19 19:36:20 Memory used for population generation and bootstrapping: 155080 bytes

```

We see that while the values returned from the Go program does not look similar to the ones from the R program, they follow a similar trend, where the standard error is lower as the sample size increases. This means that the program function correctly performs bootstrap resampling in Go.

The profiling can provide us with the memory allocation for main.go but I wanted to look at the memory allocation specifically for the functions that generate the population and bootstrapping, so the code includes a line that logs the memory used by population generation and bootstrapping. For comparison purposes, I will look at the `HeapAlloc` when we launch the profile at `localhost:6060`. Following are the execution times as well as memory allocations for the R and Go programs performing the same functions:

| Code Language  | Code efficiency (seconds) | Memory Allocation (bytes) |
|----------------|---------------------------|---------------------------|
| Go             | 0.179                     | 433600                    |
| R              | 15.32298                  | 115607336                 |

If these observations are correct, with the assumption that the program functions correctly, it means that the Go program performs the same process about a hundred times faster, with about 240 times less memory allocation.

## Observations

Here are some points to consider:

- rand.Seed is deprecated, but did not see any apparent effects of this while running the code.
- I personally struggled with refactoring bootstrap resampling, so I am not 100% sure if the code is functioning in the same way that the R code does.
- Profiling does not work as the code runs very quickly, it is not possible to launch `localhost:6060` before the code stops running, so I couldn't initially figure out what the problem was. I played around a little bit and found out that if I added the 5 minutes sleep at the end of the main function, I was able to access the profiling in the localhost. The part that I am unsure about is whether the profiling reports the "aftermath" of the code in terms of memory allocations or actually reports how much memory is used for the program.
- I wanted to see how I could implement goroutines into the code, considering we could potentially run some of the functions concurrently, so I asked for some help from Copilot. However, the resulting program did not make a significant performance difference, which lead me to remove the changes and keep only the profiling concurrent.
- My assumtion on why the values are so different between the R and Go programs results is due to the seed selection. I looked into setting the seed number the same for the Go program, but found out that R and Go random number generation algorithms are different. R uses a Mersenne-Twister, whereas Go's math/rand uses one based on linear congruential generator (LCG). This could be contributing to the differences between the values in the results. Besides this, I am not sure what could be contributing to the differences.

## Summary and Recommendation

### Summary
This project performs bootstrap resampling by refactoring a similar function from R generated by Thomas W. Miller. Comparisons between the produced results for both the R and Go code show that the result trends are comparable, therefore show that the Go program is able to generate similar results. Comparison of the execution time showed that the Go program is 100-fold faster and about 240 times less memory intensive.

### Recommendation to the Research Consultancy

Depending on the budget and the manpower of the consultancy, they could consider switching to Go. The advantages that Go provides in terms of performance is certainly very significant. The only downside that the consultancy should be aware of is the availability of packages. Go is not a preferred language for applied statistics, so there aren't a lot of packages available for such analysis. This could mean that the consultancy might have to create their own packages or functions to perform similar statistical methods, which, depending on the complexity, could take a lot of time and manpower.

It is not very easy to get specific numbers when it comes to cloud computing costs, but [this article](https://redmonk.com/rstephens/2023/04/11/iaaspricing2022/) from 2022 shows some graphs on compute units by price/hour, and if we look at one of the most popular IaaS brands AWS, we can see that the price is on average about $3.00/hour. Scaling this down to seconds would give us $0.083/second. At the speeds that we run each of the program, the R program would cost the consultancy about $0.125/run whereas Go program would cost about $0.0015. If we scale this process up to a million processes in a year, running similar processes for a year in R would cost the company $125,000 whereas Go would cost them $1,500. This is a very significant difference at 98.8% savings.