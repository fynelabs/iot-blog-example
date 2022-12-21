package main

import (
	"sort"
	"testing"

	"github.com/ssimunic/gosensors"
	"github.com/stretchr/testify/assert"
)

func TestStringsFromSensors(t *testing.T) {
	sensors, err := gosensors.NewFromFile("sensors.out")
	assert.NoError(t, err)
	assert.NotEmpty(t, sensors.String())

	sensors2, err := gosensors.NewFromFile("sensors2.out")
	assert.NoError(t, err)
	assert.NotEmpty(t, sensors.String())

	try1 := stringsFromSensors(sensors)
	try2 := stringsFromSensors(sensors)

	assert.Equal(t, try1, try2)

	try3 := stringsFromSensors(sensors2)
	assert.NotEqual(t, try1, try3)

	sort.SliceIsSorted(try1, func(i, j int) bool {
		return try1[i] < try1[j]
	})
	sort.SliceIsSorted(try3, func(i, j int) bool {
		return try3[i] < try3[j]
	})
}
