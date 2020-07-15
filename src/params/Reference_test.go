package params_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/log"
	"github.com/x1n13y84issmd42/oasis/src/params"
)

func Test_AccessContent(T *testing.T) {
	d := "foobar"
	assert.Equal(T, "foobar", params.AccessContent()(d))
}

func jsonarray(js string) interface{} {
	var d []interface{}
	err := json.Unmarshal([]byte(js), &d)
	if err != nil {
		fmt.Printf("JSON Error: %s", err.Error())
	}
	return d
}

func jsondata(js string) interface{} {
	var d interface{}
	err := json.Unmarshal([]byte(js), &d)
	if err != nil {
		fmt.Printf("JSON Error: %s", err.Error())
	}
	return d
}

func unpanic(T *testing.T, expectedMsg string) {
	if r := recover(); r != nil {
		if expectedMsg != "" {
			assert.Equal(T, expectedMsg, r)
		}
	} else {
		T.Error("Should have panicked.")
	}
}

func Test_AccessArray(T *testing.T) {
	T.Run("OK", func(T *testing.T) {
		d := jsonarray(`
		[
			"foo",
			"bar",
			"qeq",
			"baz"
		]
		`)
		access := params.AccessContent()
		assert.Equal(T, "qeq", params.AccessArray(access, 2)(d))
	})

	T.Run("Out of bounds", func(T *testing.T) {
		defer unpanic(T, "Array index 200 is out of range 0-3.")

		d := jsonarray(`
		[
			"foo",
			"bar",
			"qeq",
			"baz"
		]
		`)

		access := params.AccessContent()
		assert.Equal(T, "IRRELEVANT", params.AccessArray(access, 200)(d))
	})

	T.Run("Not an array", func(T *testing.T) {
		defer unpanic(T, "Not an array.")

		d := jsondata(`{"foo": 42}`)

		access := params.AccessContent()
		assert.Equal(T, "IRRELEVANT", params.AccessArray(access, 200)(d))
	})
}

func Test_AccessObject(T *testing.T) {
	T.Run("OK", func(T *testing.T) {
		d := jsondata(`
		{
			"foo": "F00",
			"bar": "B4R"
		}
		`)
		access := params.AccessContent()
		assert.Equal(T, "F00", params.AccessObject(access, "foo")(d))
		assert.Equal(T, "B4R", params.AccessObject(access, "bar")(d))
	})

	T.Run("Out of bounds", func(T *testing.T) {
		defer unpanic(T, "Field 'FAILURE' is not found in the object.")

		d := jsondata(`
		{
			"foo": "F00",
			"bar": "B4R"
		}
		`)

		access := params.AccessContent()
		assert.Equal(T, "IRRELEVANT", params.AccessObject(access, "FAILURE")(d))
	})

	T.Run("Not an object", func(T *testing.T) {
		defer unpanic(T, "Not an object.")

		d := jsondata(`42`)

		access := params.AccessContent()
		assert.Equal(T, "IRRELEVANT", params.AccessObject(access, "IRRELEVANT")(d))
	})
}

func Test_ParseSelector(T *testing.T) {
	T.Run("[0] OK", func(T *testing.T) {
		_, actual := params.ParseSelector("[0]", log.NewPlain(0))
		assert.Equal(T, "", actual)
	})

	T.Run("[0][200][300][42] OK", func(T *testing.T) {
		_, actual := params.ParseSelector("[0][200][300][42]", log.NewPlain(0))
		assert.Equal(T, "", actual)
	})

	T.Run(".UsEr_NaMe OK", func(T *testing.T) {
		_, actual := params.ParseSelector(".UsEr_NaMe", log.NewPlain(0))
		assert.Equal(T, "", actual)
	})

	T.Run(".users[13].id OK", func(T *testing.T) {
		_, actual := params.ParseSelector(".users[13].id", log.NewPlain(0))
		assert.Equal(T, "", actual)
	})

	T.Run("[0][200][-1][42] FAIL", func(T *testing.T) {
		_, actual := params.ParseSelector("[0][200][-1][42]", log.NewPlain(0))
		assert.Equal(T, "[-1][42]", actual)
	})

	T.Run(".user name FAIL", func(T *testing.T) {
		_, actual := params.ParseSelector(".user name", log.NewPlain(0))
		assert.Equal(T, " name", actual)
	})

	T.Run(".users[13]-.id FAIL", func(T *testing.T) {
		_, actual := params.ParseSelector(".users[13]-.id", log.NewPlain(0))
		assert.Equal(T, "-.id", actual)
	})
}

func Test_Reference(T *testing.T) {

	T.Run("Cast/string", func(T *testing.T) {
		ref := params.Reference{}
		assert.Equal(T, "hello", ref.Cast("hello"))
	})

	T.Run("Cast/int64", func(T *testing.T) {
		ref := params.Reference{}
		assert.Equal(T, "42", ref.Cast(42))
	})

	T.Run("Cast/float64", func(T *testing.T) {
		ref := params.Reference{}
		assert.Equal(T, "42.000001", ref.Cast(42.000001))
		assert.Equal(T, "42.240000", ref.Cast(42.24))
	})

	T.Run("Cast/bool", func(T *testing.T) {
		ref := params.Reference{}
		assert.Equal(T, "true", ref.Cast(true))
		assert.Equal(T, "false", ref.Cast(false))
	})

	T.Run("Cast/map", func(T *testing.T) {
		ref := params.Reference{}
		assert.Equal(T, "map[string]string{\"foo\":\"F00\"}", ref.Cast(map[string]string{"foo": "F00"}))
	})

	T.Run("Value/Array", func(T *testing.T) {
		ref := params.Reference{
			Result: &contract.OperationResult{
				ResponseBytes: []byte("[1, 22, 333, 4444]"),
			},
			Selector: "[1]",
			Log:      log.NewPlain(0),
		}

		assert.Equal(T, "22", ref.Value()())
	})

	T.Run("Value/Object", func(T *testing.T) {
		ref := params.Reference{
			Result: &contract.OperationResult{
				ResponseBytes: []byte(`{
					"foo": "F00",
					"bar": "B4R"
					}`),
			},
			Selector: ".bar",
			Log:      log.NewPlain(0),
		}

		assert.Equal(T, "B4R", ref.Value()())
	})

	T.Run("Value/String", func(T *testing.T) {
		ref := params.Reference{
			Result: &contract.OperationResult{
				ResponseBytes: []byte(`"hello"`),
			},
			Selector: "",
			Log:      log.NewPlain(0),
		}

		assert.Equal(T, "hello", ref.Value()())
	})

	T.Run("Value/Number", func(T *testing.T) {
		ref := params.Reference{
			Result: &contract.OperationResult{
				ResponseBytes: []byte(`42`),
			},
			Selector: "",
			Log:      log.NewPlain(0),
		}

		assert.Equal(T, "42", ref.Value()())
	})

	T.Run("Value/Boolean", func(T *testing.T) {
		ref := params.Reference{
			Result: &contract.OperationResult{
				ResponseBytes: []byte(`true`),
			},
			Selector: "",
			Log:      log.NewPlain(0),
		}

		assert.Equal(T, "true", ref.Value()())
	})
}
