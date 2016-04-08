package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestJSONMarshalling(t *testing.T) {
	j := []byte(`
  {
    "t":"p",
    "d":"2016-04-08T02:24:14.726Z",
    "p": {
      "a": "123abc",
      "p": [{"diffs":[[1,"helo"]],"start1":0,"start2":0,"length1":0,"length2":4}]
    }
  }`)
	m := new(message)
	json.Unmarshal(j, m)
	fmt.Println(m)
}
