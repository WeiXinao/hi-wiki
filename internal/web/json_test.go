package web

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Stu struct {
	Name any `json:"name,omitempty"`
}

func TestMarshal(t *testing.T) {
	stu := Stu{Name: nil}
	marshal, err := json.Marshal(stu)
	assert.NoError(t, err)
	t.Log(string(marshal))
}
