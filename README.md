[![GoDoc](https://godoc.org/github.com/bunyk/require?status.svg)](https://godoc.org/github.com/bunyk/require)

**This package is obsolete.** Since Go 1.16 you could use [`embed`](https://pkg.go.dev/embed) from standard library for the same purpose.

[See my blog post about this package and embed](http://bunyk.github.io/posts/go_embed)

----

The purpose of this module is to allow you to "require" text files,
so they could be hardcoded in single executable and all filesystem access
except of course process loading happens during the build. 

# Installation
```
go get github.com/bunyk/require
go install github.com/bunyk/require/hardcode
```

# Usage

Write `require.File("file.txt")` anywhere in your code to get contents of given file as a string. Then you run 
`hardcode` giving it a list of files to process, and it will give you go code to include in your program as output.

# Tutorial
Say you have file with some famous quotes, and you want to build
a small fortune program that prints one of them randomly:

```
There are two ways of constructing a software design: One way is to make it so simple that there are obviously no deficiencies and the other way is to make it so complicated that there are no obvious deficiencies.  — C.A.R. Hoare, The 1980 ACM Turing Award Lecture
The cheapest, fastest, and most reliable components are those that aren’t there.  — Gordon Bell
One of my most productive days was throwing away 1000 lines of code.  — Ken Thompson
Deleted code is debugged code.  — Jeff Sickel
```

You could write something like this:

```
package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

func main() {
	contents, err := ioutil.ReadFile("fortunes.txt")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fortunes := strings.Split(string(contents), "\n")

	rand.Seed(time.Now().UTC().UnixNano())

	fmt.Println(fortunes[rand.Intn(len(fortunes))])
}
```

But then, when you forget to include "fortunes.txt" into your software package, 
you could get this error:

```
 $ ./fortune 
open fortunes.txt: no such file or directory
```

So, you `go get "github.com/bunyk/require"`, and rewrite code like this:

```go
package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/bunyk/require"
)

func main() {
	fortunes := strings.Split(require.File("fortunes.txt"), "\n")

	rand.Seed(time.Now().UTC().UnixNano())

	fmt.Println(fortunes[rand.Intn(len(fortunes))])
}
```

Building this will produce incorrect program,
because now we need one additional preprocessing step. Do:

```bash
go install github.com/bunyk/require/hardcode
hardcode --package="main" $(find . | grep "\.go$") > fortunes.go
```

`hardcode` will go over source files you provided, find all calls to `require.File()` there
open contents of files that are required and generate to standart output 
additional go source file for package chosen in arguments, with all the files content hardcoded.

When you build your project now, `require.File()` will return you string 
with that file contents without looking for any files on filesystem.

