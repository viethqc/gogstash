package env

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_GetString(t *testing.T) {
	var (
		assert = assert.New(t)
		value  string
	)

	value = GetString("IMPOSSIBLE_ENV_KEY", "Impossible value !")
	assert.Equal("Impossible value !", value, "Get ENV 'IMPOSSIBLE_ENV_KEY'")

	path := os.Getenv("PATH")
	value = GetString("PATH", "")
	assert.Equal(path, value, "Get ENV 'PATH'")
}

func Test_GetBool(t *testing.T) {
	var (
		assert = assert.New(t)
		value  bool
	)

	value = GetBool("IMPOSSIBLE_ENV_KEY", true)
	assert.Equal(true, value, "Get ENV 'IMPOSSIBLE_ENV_KEY'")

	value = GetBool("IMPOSSIBLE_ENV_KEY", false)
	assert.Equal(false, value, "Get ENV 'IMPOSSIBLE_ENV_KEY'")

	os.Setenv("IMPOSSIBLE_ENV_KEY_TRUE", "true")
	value = GetBool("IMPOSSIBLE_ENV_KEY_TRUE", true)
	assert.Equal(true, value, "Get ENV 'IMPOSSIBLE_ENV_KEY_TRUE'")

	value = GetBool("IMPOSSIBLE_ENV_KEY_TRUE", false)
	assert.Equal(true, value, "Get ENV 'IMPOSSIBLE_ENV_KEY_TRUE'")

	os.Setenv("IMPOSSIBLE_ENV_KEY_FALSE", "false")
	value = GetBool("IMPOSSIBLE_ENV_KEY_FALSE", true)
	assert.Equal(false, value, "Get ENV 'IMPOSSIBLE_ENV_KEY_FALSE'")

	value = GetBool("IMPOSSIBLE_ENV_KEY_FALSE", false)
	assert.Equal(false, value, "Get ENV 'IMPOSSIBLE_ENV_KEY_FALSE'")
}

func Test_GetInt(t *testing.T) {
	var (
		assert = assert.New(t)
		value  int
	)

	value = GetInt("IMPOSSIBLE_ENV_KEY", 9527)
	assert.Equal(9527, value, "Get ENV 'IMPOSSIBLE_ENV_KEY'")

	os.Setenv("IMPOSSIBLE_ENV_KEY_9527", "9527")
	value = GetInt("IMPOSSIBLE_ENV_KEY_9527", 168)
	assert.Equal(9527, value, "Get ENV 'IMPOSSIBLE_ENV_KEY_9527'")

	os.Setenv("IMPOSSIBLE_ENV_KEY_9527", "Not a number")
	value = GetInt("IMPOSSIBLE_ENV_KEY_FALSE", 168)
	assert.Equal(168, value, "Get ENV 'IMPOSSIBLE_ENV_KEY_9527'")
}
