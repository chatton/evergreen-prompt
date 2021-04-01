package flagutil

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestBuildVariantFlagExtraction(t *testing.T){
	bv := GetBuildVariantValue("--buildvariant myvalue")
	assert.Equal(t, "myvalue", bv)
}


func TestTaskFlagExtraction(t *testing.T){
	task := GetTaskValue("--task this_is_my_task")
	assert.Equal(t, "this_is_my_task", task)
}


func TestDescriptionFlagExtraction(t *testing.T){
	desc := GetDescriptionValue(`--description "this is my description"`)
	assert.Equal(t, `"this is my description"`, desc)
}


//\s+ID\s:\s([a-zA-Z0-9]+)