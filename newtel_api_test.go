package newtel_api

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetData(t *testing.T) {

	data := CallPasswordData{
		DstNumber: "+79081234567",
		Async:     "0",
		Pin:       "1" + "1234",
		Timeout:   "30",
	}
	binary , _ := json.Marshal(data)
	out:= "{\"dstNumber\":\"+79081234567\",\"async\":\"0\",\"pin\":\"11234\",\"timeout\":\"30\"}"

	require.Equal(t, out, string(binary))
}
