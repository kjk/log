package log

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	err := os.RemoveAll("test_dir")
	assert.Nil(t, err)
	pathFormat := filepath.Join("test_dir", "second", "2006-01-02.txt")
	pathExp := time.Now().Format(pathFormat)
	f, err := NewDailyRotateFile(pathFormat)
	assert.Nil(t, err)
	assert.NotNil(t, f)
	assert.Equal(t, pathExp, f.path)
	n, err := f.WriteString("hello\n")
	assert.Nil(t, err)
	assert.Equal(t, n, 6)
	n, err = f.Write([]byte("bar\n"))
	assert.Nil(t, err)
	assert.Equal(t, n, 4)
	f.Close()
	d, err := ioutil.ReadFile(pathExp)
	assert.Nil(t, err)
	assert.Equal(t, string(d), "hello\nbar\n")
	err = os.RemoveAll("test_dir")
	assert.Nil(t, err)
}
