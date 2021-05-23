package jsre

import (
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/dop251/goja"
)

type testNativeObjectBinding struct {
	vm *goja.Runtime
}

type msg struct {
	Msg string
}

func (no *testNativeObjectBinding) TestMethod(call goja.FunctionCall) goja.Value {
	m := call.Argument(0).ToString().String()
	return no.vm.ToValue(&msg{m})
}

func newWithTestJS(t *testing.T, testjs string) (*JSRE, string) {
	dir, err := ioutil.TempDir("", "jsre-test")
	if err != nil {
		t.Fatal("cannot create temporary directory:", err)
	}
	if testjs != "" {
		if err := ioutil.WriteFile(path.Join(dir, "test.js"), []byte(testjs), os.ModePerm); err != nil {
			t.Fatal("cannot create test.js:", err)
		}
	}
	jsre := New(dir, os.Stdout)
	return jsre, dir
}

func TestExec(t *testing.T) {
	jsre, dir := newWithTestJS(t, `msg = "testMsg"`)
	defer os.RemoveAll(dir)

	err := jsre.Exec("test.js")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	val, err := jsre.Run("msg")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if val.ExportType().Kind() != reflect.String {
		t.Errorf("expected string value, got %v", val)
	}
	exp := "testMsg"
	got := val.ToString().String()
	if exp != got {
		t.Errorf("expected '%v', got '%v'", exp, got)
	}
	jsre.Stop(false)
}

func TestNatto(t *testing.T) {
	jsre, dir := newWithTestJS(t, `setTimeout(function(){msg = "testMsg"}, 1);`)
	defer os.RemoveAll(dir)

	err := jsre.Exec("test.js")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	time.Sleep(100 * time.Millisecond)
	val, err := jsre.Run("msg")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if val.ExportType().Kind() != reflect.String {
		t.Errorf("expected string value, got %v", val)
	}
	exp := "testMsg"
	got := val.ToString().String()
	if exp != got {
		t.Errorf("expected '%v', got '%v'", exp, got)
	}
	jsre.Stop(false)
}

func TestBind(t *testing.T) {
	jsre := New("", os.Stdout)
	defer jsre.Stop(false)

	jsre.Set("no", &testNativeObjectBinding{vm: jsre.vm})

	_, err := jsre.Run(`no.TestMethod("testMsg")`)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestLoadScript(t *testing.T) {
	jsre, dir := newWithTestJS(t, `msg = "testMsg"`)
	defer os.RemoveAll(dir)

	_, err := jsre.Run(`loadScript("test.js")`)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	val, err := jsre.Run("msg")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if val.ExportType().Kind() != reflect.String {
		t.Errorf("expected string value, got %v", val)
	}
	exp := "testMsg"
	got := val.ToString().String()
	if exp != got {
		t.Errorf("expected '%v', got '%v'", exp, got)
	}
	jsre.Stop(false)
}
