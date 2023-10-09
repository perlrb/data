//go:build mage

package main

import (
	"encoding/csv"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
	"strings"
)

type Product struct {
	SKUs         []string `yaml:"skus" json:"skus,omitempty"`
	Name         string   `yaml:"name" json:"name,omitempty"`
	Red          int      `yaml:"red" json:"red,omitempty"`
	Green        int      `yaml:"green" json:"green,omitempty"`
	Blue         int      `yaml:"blue" json:"blue,omitempty"`
	Hex          string   `yaml:"hex" json:"hex,omitempty"`
	Manufacturer string   `yaml:"manufacturer,omitempty" json:"manufacturer,omitempty"`
}

type Manufacturer struct {
	Name     string    `yaml:"name" json:"name,omitempty"`
	Products []Product `yaml:"products" json:"products,omitempty"`
}

func Seed() error {
	var manufacturersMap = make(map[string]*Manufacturer)

	f, err := os.Open("seed.csv")
	if err != nil {
		return err
	}
	defer f.Close()
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	for i, record := range records {
		if i == 0 {
			// ignore header
			continue
		}
		manufacturerName := record[0]
		if _, ok := manufacturersMap[strings.ToLower(manufacturerName)]; !ok {
			newManufacturer := &Manufacturer{
				Name:     manufacturerName,
				Products: make([]Product, 0),
			}
			manufacturersMap[strings.ToLower(manufacturerName)] = newManufacturer
		}
		r, err := strconv.Atoi(record[4])
		if err != nil {
			r = -1
		}
		g, err := strconv.Atoi(record[5])
		if err != nil {
			g = -1
		}
		b, err := strconv.Atoi(record[6])
		if err != nil {
			b = -1
		}
		tmpSkus := []string{record[1], record[2]}
		finalSkus := make([]string, 0)
		for _, str := range tmpSkus {
			if str != "" {
				finalSkus = append(finalSkus, str)
			}
		}
		manufacturersMap[strings.ToLower(manufacturerName)].Products = append(manufacturersMap[strings.ToLower(manufacturerName)].Products, Product{
			SKUs:  finalSkus,
			Name:  record[3],
			Red:   r,
			Green: g,
			Blue:  b,
			Hex:   record[7],
		})
	}
	yamlBytes, err := yaml.Marshal(manufacturersMap)
	if err != nil {
		return err
	}
	f, err = os.Create("products.yaml")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(yamlBytes)
	if err != nil {
		return err
	}
	return nil
}
