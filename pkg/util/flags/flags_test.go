package flags

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestBuildVariantFlagExtraction(t *testing.T) {
	bv := getBuildVariantValue("--buildvariant myvalue")
	assert.Equal(t, "myvalue", bv)
}

func TestTaskFlagExtraction(t *testing.T) {
	t.Run("Test task as only flag", func(t *testing.T) {
		task := getTaskValue("--task this_is_my_task")
		assert.Equal(t, "this_is_my_task", task)
	})
	t.Run("Test task as first flag", func(t *testing.T) {
		task := getTaskValue("--task this_is_my_task --buildvariant this_is_my_bv")
		assert.Equal(t, "this_is_my_task", task)
	})
	t.Run("Test task as second flag", func(t *testing.T) {
		task := getTaskValue("--buildvariant this_is_my_bv --task this_is_my_task")
		assert.Equal(t, "this_is_my_task", task)
	})
}

func TestDescriptionFlagExtraction(t *testing.T) {
	desc := getDescriptionValue(`--description "this is my description"`)
	assert.Equal(t, `"this is my description"`, desc)
}

func TestProjectFlagExtraction(t *testing.T) {
	project := getProjectValue("--project ops-manager-kubernetes")
	assert.Equal(t, "ops-manager-kubernetes", project)
}

func TestExtractFlags(t *testing.T) {

	t.Run("Test even number of values", func(t *testing.T) {
		input := "patch create --task this_is_my_task --buildvariant this_is_my_bv"
		flags := extractFlags(input, patchCreate)

		task, _ := getValueFromFlagKey("--task", flags)
		assert.Equal(t, "this_is_my_task", task)

		bv, _ := getValueFromFlagKey("--buildvariant", flags)
		assert.Equal(t, "this_is_my_bv", bv)
	})

	t.Run("Test with flags that have no value as last item", func(t *testing.T) {
		input := "patch create --task this_is_my_task --buildvariant this_is_my_bv --uncommited"
		flags := extractFlags(input, patchCreate)

		task, _ := getValueFromFlagKey("--task", flags)
		assert.Equal(t, "this_is_my_task", task)
		bv, _ := getValueFromFlagKey("--buildvariant", flags)
		assert.Equal(t, "this_is_my_bv", bv)

		_, ok := getValueFromFlagKey("--uncommited", flags)
		assert.Equal(t, true, ok)
	})

	t.Run("Test with flags that have no value as middle item", func(t *testing.T) {
		input := "patch create --task this_is_my_task --buildvariant this_is_my_bv --uncommited --priority 100"
		flags := extractFlags(input, patchCreate)

		task, _ := getValueFromFlagKey("--task", flags)
		assert.Equal(t, "this_is_my_task", task)

		bv, _ := getValueFromFlagKey("--buildvariant", flags)
		assert.Equal(t, "this_is_my_bv", bv)
		_, ok := getValueFromFlagKey("--uncommited", flags)
		assert.Equal(t, true, ok)

		priority, _ := getValueFromFlagKey("--priority", flags)
		assert.Equal(t, "100", priority)
	})

}
