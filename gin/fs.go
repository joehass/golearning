package gin

import (
	"net/http"
	"os"
)

type onlyFilesFS struct {
	fs http.FileSystem
}

type netueredReaddirFile struct {
	http.File
}

//返回一个http.FileSystem，在router.Static()内部使用
//ListDirectory = true,返回相同的http.Dir()，否则将会返回一个阻止列出文件目录的filesystem
func Dir(root string, ListDirectory bool) http.FileSystem {
	fs := http.Dir(root)
	if ListDirectory {
		return fs
	}

	return &onlyFilesFS{fs}
}

//打开文件
func (fs onlyFilesFS) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return netueredReaddirFile{f}, nil
}

func (f netueredReaddirFile) Readdir(int) ([]os.FileInfo, error) {
	return nil, nil
}
