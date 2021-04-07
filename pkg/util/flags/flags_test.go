package flags

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestBuildVariantFlagExtraction(t *testing.T) {
	bv := GetBuildVariantValue("--buildvariant myvalue")
	assert.Equal(t, "myvalue", bv)
}

func TestTaskFlagExtraction(t *testing.T) {
	t.Run("Test task as only flag", func(t *testing.T) {
		task := GetTaskValue("--task this_is_my_task")
		assert.Equal(t, "this_is_my_task", task)
	})
	t.Run("Test task as first flag", func(t *testing.T) {
		task := GetTaskValue("--task this_is_my_task --buildvariant this_is_my_bv")
		assert.Equal(t, "this_is_my_task", task)
	})
	t.Run("Test task as second flag", func(t *testing.T) {
		task := GetTaskValue("--buildvariant this_is_my_bv --task this_is_my_task")
		assert.Equal(t, "this_is_my_task", task)
	})
}

func TestDescriptionFlagExtraction(t *testing.T) {
	desc := GetDescriptionValue(`--description "this is my description"`)
	assert.Equal(t, `"this is my description"`, desc)
}

func TestProjectFlagExtraction(t *testing.T) {
	project := GetProjectValue("--project ops-manager-kubernetes")
	assert.Equal(t, "ops-manager-kubernetes", project)
}

func TestExtractFlags(t *testing.T) {

	t.Run("Test even number of values", func(t *testing.T) {
		input := "patch create --task this_is_my_task --buildvariant this_is_my_bv"
		flags := extractFlags(input, patchCreate)

		assert.Equal(t, "this_is_my_task", flags["--task"])
		assert.Equal(t, "this_is_my_bv", flags["--buildvariant"])
	})

	t.Run("Test with flags that have no value as last item", func(t *testing.T) {
		input := "patch create --task this_is_my_task --buildvariant this_is_my_bv --uncommited"
		flags := extractFlags(input, patchCreate)

		assert.Equal(t, "this_is_my_task", flags["--task"])
		assert.Equal(t, "this_is_my_bv", flags["--buildvariant"])
		assert.Equal(t, "", flags["--uncommited"])
	})

	t.Run("Test with flags that have no value as middle item", func(t *testing.T) {
		input := "patch create --task this_is_my_task --buildvariant this_is_my_bv --uncommited --priority 100"
		flags := extractFlags(input, patchCreate)

		assert.Equal(t, "this_is_my_task", flags["--task"])
		assert.Equal(t, "this_is_my_bv", flags["--buildvariant"])
		assert.Equal(t, "", flags["--uncommited"])
		assert.Equal(t, "100", flags["--priority"])
	})

}
