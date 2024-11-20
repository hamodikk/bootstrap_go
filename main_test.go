package main

import (
	"math"
	"testing"
)

func TestGeneratePopulation(t *testing.T) {
	population := generatePopulation(100, 100.0, 10.0)
	if len(population) != 100 {
		t.Errorf("Expected population size of 100, got %d", len(population))
	}
}

func TestMedian(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	expected := 3.0
	result := median(data)
	if result != expected {
		t.Errorf("Expected median %f, got %f", expected, result)
	}
}

func TestMean(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	expected := 3.0
	result := mean(data)
	if result != expected {
		t.Errorf("Expected mean %f, got %f", expected, result)
	}
}

func TestStandardError(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	expected := 0.70710678119
	result := standardError(data)
	if math.Abs(result-expected) > 1e-6 {
		t.Errorf("Expected standard error %f, got %f", expected, result)
	}
}

func TestBootstrap(t *testing.T) {
	population := generatePopulation(100, 100.0, 10.0)
	seMean, seMedian := bootstrap(population, 10, 10)
	if seMean <= 0 {
		t.Errorf("Expected positive SE Mean, got %f", seMean)
	}
	if seMedian <= 0 {
		t.Errorf("Expected positive SE Median, got %f", seMedian)
	}
}

func BenchmarkBootstrap(b *testing.B) {
	population := generatePopulation(10000, 100.0, 10.0)
	for i := 0; i < b.N; i++ {
		bootstrap(population, 100, 100)
	}
}
