package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchAvailableSubnet(t *testing.T) {
	// case normal
	result, err := SearchAvailableSubnet("192.168.0.0/16", []string{
		"192.168.136.0/21",
		"192.168.1.0/24",
		"192.168.0.0/24",
		"192.168.2.0/30",
		"10.0.0.0/24",
		"192.168.128.0/20",
	}, 8)

	assert.Nil(t, err)
	assert.Equal(t, "192.168.3.0/24", result.String())

	// case normal
	result, err = SearchAvailableSubnet("192.168.0.0/16", []string{
		"192.168.0.0/20",
	}, 8)

	assert.Nil(t, err)
	assert.Equal(t, "192.168.16.0/24", result.String())

	// case no exists
	result, err = SearchAvailableSubnet("192.168.0.0/16", []string{}, 8)
	assert.Equal(t, "192.168.0.0/24", result.String())
	assert.Nil(t, err)

	// case error maskInc
	result, err = SearchAvailableSubnet("192.168.0.0/16", []string{
		"10.0.0.0/24",
		"192.168.128.0/20",
	}, 20)
	assert.NotNil(t, err)
}
