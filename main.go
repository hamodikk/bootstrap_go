package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"sort"
	"time"
)

// Generate a population with given mean and standard deviation
func generatePopulation(size int, mean, stddev float64) []float64 {
	population := make([]float64, size)
	// Set seed for reproducibility
	seed := int64(42)
	rand.Seed(seed)
	for i := 0; i < size; i++ {
		population[i] = rand.NormFloat64()*stddev + mean
	}
	return population
}

// Calculate the median
func median(data []float64) float64 {
	sort.Float64s(data)
	n := len(data)
	if n%2 == 0 {
		return (data[n/2-1] + data[n/2]) / 2.0
	}
	return data[n/2]
}

// Calculate the mean
func mean(data []float64) float64 {
	sum := 0.0
	for _, value := range data {
		sum += value
	}
	return sum / float64(len(data))
}

// Calculate the standard error
func standardError(data []float64) float64 {
	mean := mean(data)
	variance := 0.0
	for _, value := range data {
		variance += math.Pow(value-mean, 2)
	}
	stdDev := math.Sqrt(variance / float64(len(data)-1))
	return stdDev / math.Sqrt(float64(len(data)))
}

// Perform bootstrap resampling and calculate the standard error of the mean and median
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

// Perform Central Limit Theorem
func centralLimitTheorem(population []float64, sampleSize, numSamples int) float64 {
	sampleMeans := make([]float64, numSamples)
	for i := 0; i < numSamples; i++ {
		sample := make([]float64, sampleSize)
		for j := 0; j < sampleSize; j++ {
			sample[j] = population[rand.Intn(len(population))]
		}
		sampleMeans[i] = mean(sample)
	}
	return standardError(sampleMeans)
}

func main() {
	// Start pprof for profiling
	go func() {
		err := http.ListenAndServe("localhost:6060", nil)
		if err != nil {
			log.Fatalf("Failed to start pprof server: %v", err)
		}
	}()

	// Set up logging
	logFile, err := os.OpenFile("bootstrap.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Parameters
	populationSize := 10000
	mean := 100.0
	stddev := 10.0
	bootstrapSamples := 100
	sampleSizes := []int{25, 100, 225, 400}
	numSamples := 1000

	// Memory usage before generating population and bootstrapping
	var memBefore, memAfter runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	// Generate population
	population := generatePopulation(populationSize, mean, stddev)
	log.Println("Generated population")

	// Perform bootstrap for each sample size
	for _, sampleSize := range sampleSizes {
		// Perform Central Limit Theorem
		cltSeMean := centralLimitTheorem(population, sampleSize, numSamples)
		seMean, seMedian := bootstrap(population, bootstrapSamples, sampleSize)
		log.Printf("Samples of size n = %d\n", sampleSize)
		log.Printf("  SE Mean from Central Limit Theorem for n = %d: %.2f\n", sampleSize, cltSeMean)
		log.Printf("  SE Mean from Bootstrap Samples: %.2f\n", seMean)
		log.Printf("  SE Median from Bootstrap Samples: %.2f\n", seMedian)
		fmt.Printf("Samples of size n = %d\n", sampleSize)
		fmt.Printf("  SE Mean from Central Limit Theorem for n = %d: %.2f\n", sampleSize, cltSeMean)
		fmt.Printf("  SE Mean from Bootstrap Samples: %.2f\n", seMean)
		fmt.Printf("  SE Median from Bootstrap Samples: %.2f\n", seMedian)
	}

	// Memory usage after generating population and bootstrapping
	runtime.ReadMemStats(&memAfter)
	log.Printf("Memory used for population generation and bootstrapping: %d bytes\n", memAfter.Alloc-memBefore.Alloc)

	// Add delay to the program to allow profiling access
	time.Sleep(5 * time.Minute)
}
