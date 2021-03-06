// Code generated by "esc -o assert.star.go -pkg starlarktest assert.star"; DO NOT EDIT.

package starlarktest

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	if !f.isDir {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is not directory", f.name)
	}

	fis, ok := _escDirs[f.local]
	if !ok {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is directory, but we have no info about content of this dir, local=%s", f.name, f.local)
	}
	limit := count
	if count <= 0 || limit > len(fis) {
		limit = len(fis)
	}

	if len(fis) == 0 && count > 0 {
		return nil, io.EOF
	}

	return fis[0:limit], nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func Dir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		_ = f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(useLocal bool, name string) []byte {
	b, err := FSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(useLocal bool, name string) (string, error) {
	b, err := FSByte(useLocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(useLocal bool, name string) string {
	return string(FSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/assert.star": {
		name:    "assert.star",
		local:   "assert.star",
		size:    1526,
		modtime: 1587336348,
		compressed: `
H4sIAAAAAAAC/3RUTXPjNgy961e82uOptOPmB3jqUw+9dXroPcOIkMVZmnQAaCX313dASok3SX2xAAKP
Dx+Pe/zN5KmPjsnjZQpRfwtJMGSGjkFwzX6KdGr2zR7EnLm9yqU7gemWWeFS9SIk/Jl/FSiJYmB3pTnz
d8xBxzwpRhc1pAtooX7SkNNTswf+sRuCwAfRkHpLzFfoSG9MMLgQMUypt6Qj5jH0Y0GTn7F6p/3YDt0J
9MPFySlhaDu45MGkEydBsJx6GPLG+0oi7kJHhAEu3Zs9roZE0oryETenSpzeC55H0pEYorxFgukyRceg
5cYkYuBrnjGrHWy/ffs+O75Id4JDn5MoT71mLq12a5TFPw9M9C+1S3dC/SwdMeKEpVREP4jvOlpDmVw/
updItXdBn8qo/oiBkgqu7o5JCoLQWx8FmuFpCKmcBEaeUxmdQboXUXY18KlpPA14ptd2OeLenRoA1qsF
v5xxr6b96m7sDmz+A+9wQM3oVoREnxDO/4Nw/hpBeaK2z8kfcZULztg5EeIyTFsT8rt39JTVuuw/4tv2
rnhRPzCynHbB72++R1plUS0ikgh0dAkH+Uyyz0ldSPIB+l5SQ8LyFbLPJBtjy/4S2UqUdnjYyYK0NqEe
oxqFnz4uex6KHGqQibLslNyoD0MgX7k87Qpibe4mqK2C4j3jr5zoUwkPF8nU90SePKZEy416JR/vaGeX
dJOcIduiHbizMrd6CirFdRCbDNfTMvPPY/lCeu1BOvjg31HWe+2gtPUnyK5pVpWdN+kB+/KwLaZ48rjE
/OIi5jHLpsMgH56pmrgJrGnqIHDe5P84q92xWOVpO1dy1UOvRoJeq5UKpUTVsu032/6rJxr8c9RqbYtn
vu37/Z7iLh/Hpmv+CwAA///xOAnW9gUAAA==
`,
	},
}

var _escDirs = map[string][]os.FileInfo{}
