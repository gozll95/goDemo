package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	fmt.Println("On Unix:")
	fmt.Println(filepath.Join("a", "b", "c"))
	fmt.Println(filepath.Join("a", "b/c"))
	fmt.Println(filepath.Join("a/b", "c"))
	fmt.Println(filepath.Join("a/b", "/c"))
}

// On Unix:
// a/b/c
// a/b/c
// a/b/c
// a/b/c



func Abs(path string) (string, error)

package main
import (
    "path/filepath"
    "os"
    "fmt"
)

func main() {
    //
    pName := os.Args[0]
    absName, err := filepath.Abs(pName)
    if err != nil{
        fmt.Println(err)
    }

    fmt.Println(absName)
}
func Base(path string) string

package main
import (
    "path/filepath"
    "fmt"
)

func main() {
    baseName := filepath.Base("/a/b/c/e.txt")
    fmt.Println(baseName)
}
func Clean(path string) string

package main
import (
    "path/filepath"
    "fmt"
)

func main() {
    p := "../..//././//a/b/c.txt"
    pc := filepath.Clean(p)
    fmt.Println(pc)
}
func Dir(path string) string

package main
import (
    "path/filepath"
    "fmt"
)

func main() {
    d := filepath.Dir("/a/b/c/d.txt")
    fmt.Println(d)
}
func EvalSymlinks(path string) (string, error)

package main
import (
    "path/filepath"
    "fmt"
)

func main() {
    e, _:= filepath.EvalSymlinks("/Users/hyhu/SourcePrj/mysourceprj")
    fmt.Println(e)
}
func Ext(path string) string

package main
import (
    "path/filepath"
    "fmt"
)

func main() {
    e := filepath.Ext("/Users/1.txt")
    fmt.Println(e)
}
func FromSlash(path string) string

package main
import (
    "path/filepath"
    "fmt"
)

func main() {
    //windows下有效果
    r := filepath.FromSlash("/a//b/c/d.txt")
    fmt.Println(r)
}
func Glob(pattern string) (matches []string, err error)

package main
import (
    "path/filepath"
    "fmt"
)

func main() {
    m,_ := filepath.Glob("/usr/*")
    fmt.Println(m)
}
func HasPrefix(p, prefix string) bool

Go1.7中已经废弃使用
func IsAbs(path string) bool

package main
import (
    "path/filepath"
    "fmt"
)

func main() {
    b := filepath.IsAbs("/a/b/c/d.txt")
    fmt.Println(b)
    b = filepath.IsAbs("d.txt")
    fmt.Println(b)
}
func Join(elem ...string) string

package main

import (
    "fmt"
    "path/filepath"
)

func main() {
    fmt.Println("On Unix:")
    fmt.Println(filepath.Join("a", "b", "c"))
    fmt.Println(filepath.Join("a", "b/c"))
    fmt.Println(filepath.Join("a/b", "c"))
    fmt.Println(filepath.Join("a/b", "/c"))
}
func Match(pattern, name string) (matched bool, err error)

package main
import (
    "path/filepath"
    "fmt"
)

func main() {
    //windows下有效果
    m,_ := filepath. Match("/usr/*", "/usr/local")
    fmt.Println(m)
}
func Rel(basepath, targpath string) (string, error)

package main

import (
    "fmt"
    "path/filepath"
)

func main() {
    paths := []string{
        "/a/b/c",
        "/b/c",
        "./b/c",
    }
    base := "/a"

    fmt.Println("On Unix:")
    for _, p := range paths {
        rel, err := filepath.Rel(base, p)
        fmt.Printf("%q: %q %v\n", p, rel, err)
    }

}
func Split(path string) (dir, file string)

package main

import (
    "fmt"
    "path/filepath"
)

func main() {
    paths := []string{
        "/home/arnie/amelia.jpg",
        "/mnt/photos/",
        "rabbit.jpg",
        "/usr/local//go",
    }
    fmt.Println("On Unix:")
    for _, p := range paths {
        dir, file := filepath.Split(p)
        fmt.Printf("input: %q\n\tdir: %q\n\tfile: %q\n", p, dir, file)
    }
}
func SplitList(path string) []string

package main

import (
    "fmt"
    "path/filepath"
)

func main() {
    fmt.Println("On Unix:", filepath.SplitList("/a/b/c:/usr/bin"))
}
func ToSlash(path string) string

package main
import (
    "path/filepath"
    "fmt"
)

func main() {
    //windows下有效果
    r := filepath.ToSlash("\\a\\b\\c/d.txt")
    fmt.Println(r)
}
func VolumeName(path string) string

windows平台下
C:\foo\bar 对应的结果 C:

\\host\share\foo 对应的结果 \\host\share

package main
import (
    "path/filepath"
    "fmt"
)

func main() {
    v := filepath.VolumeName("C:\foo\bar")
    fmt.Println(v)
    v = filepath.VolumeName(`\\host\share\foo"`)
    fmt.Println(v)
}
func Walk(root string, walkFn WalkFunc) error

同下
type WalkFunc

package main
import (
    "path/filepath"
    "fmt"
    "os"
)

func MyWalkFunc (path string, info os.FileInfo, err error) error{
    var e error

    fmt.Println(path, info.Name(), info.Size(), info.Mode(), info.ModTime())

    return e
}

func main() {
    filepath.Walk("/", MyWalkFunc)
}

作者：小包子你
链接：http://www.jianshu.com/p/0e191bde42e3
來源：简书
著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。