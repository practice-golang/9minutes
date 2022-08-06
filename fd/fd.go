package fd

import (
	"9minutes/router"
	"io/fs"
	"os"
	"sort"
)

const (
	_       = iota
	NOTSORT // 1 - do not sort
	NAME    // filename
	SIZE    // filesize
	TIME    // filetime
	ASC     // ascending
	DESC    // descending
)

func sortByName(a, b fs.FileInfo) bool {
	switch true {
	case a.IsDir() && !b.IsDir():
		return true
	case !a.IsDir() && b.IsDir():
		return false
	default:
		return a.Name() < b.Name()
	}
}

func sortBySize(a, b fs.FileInfo) bool {
	switch true {
	case a.IsDir() && !b.IsDir():
		return true
	case !a.IsDir() && b.IsDir():
		return false
	default:
		return a.Size() < b.Size()
	}
}

func sortByTime(a, b fs.FileInfo) bool {
	// switch true {
	// case a.IsDir() && !b.IsDir():
	// 	return true
	// case !a.IsDir() && b.IsDir():
	// 	return false
	// default:
	// 	return a.ModTime().Format("20060102150405") < b.ModTime().Format("20060102150405")
	// }

	return a.ModTime().Format("20060102150405") < b.ModTime().Format("20060102150405")
}

func Dir(path string, sortby, direction int) ([]fs.DirEntry, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if sortby > 1 && sortby < 5 {
		sort.Slice(files, func(a, b int) bool {
			aInfo, _ := files[a].Info()
			bInfo, _ := files[b].Info()

			switch sortby {
			case NAME:
				if direction == DESC {

					return !sortByName(aInfo, bInfo)
				}
				return sortByName(aInfo, bInfo)
			case SIZE:
				if direction == DESC {
					return !sortBySize(aInfo, bInfo)
				}
				return sortBySize(aInfo, bInfo)
			case TIME:
				if direction == DESC {
					return !sortByTime(aInfo, bInfo)
				}
				return sortByTime(aInfo, bInfo)
			default:
				return false
			}
		})
	}

	return files, nil
}
func CheckFileExists(path string, isEmbed bool) (result bool) {
	result = false

	switch isEmbed {
	case true:
		ef, err := fs.Stat(router.Content, path)
		if err == nil && ef != nil {
			result = true
		}
	case false:
		f, err := os.Stat(path)
		if err == nil && f != nil {
			result = true
		}
	}

	return
}
