package contract_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

func captureOutput(f func()) string {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	stdout := os.Stdout
	os.Stdout = w
	defer func() {
		os.Stdout = stdout
	}()

	f()
	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}

func Test_TabFn(T *testing.T) {
	log := log.NewPlain(1)
	tab := contract.Tab(3)

	expected := "      "
	actual := captureOutput(func() {
		tab(log)
	})
	assert.Equal(T, expected, actual)

	tab = tab.Shift().Shift()
	expected = "          "
	actual = captureOutput(func() {
		tab(log)
	})
	assert.Equal(T, expected, actual)
}
