package apimgr

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_getPackagePath(t *testing.T) {
	require := require.New(t)
	require.NotNil(require)

	pkgpath := getPackagePath(reflect.ValueOf("github.com/viethqc/gogstash/KDGoLib/apimgr"))
	require.Contains(pkgpath, "github.com/viethqc/gogstash/KDGoLib/apimgr")

	pkgpath = getPackagePath(reflect.ValueOf(Manager{}))
	require.Contains(pkgpath, "github.com/viethqc/gogstash/KDGoLib/apimgr")

	pkgpath = getPackagePath(reflect.ValueOf(&Manager{}))
	require.Contains(pkgpath, "github.com/viethqc/gogstash/KDGoLib/apimgr")
}
