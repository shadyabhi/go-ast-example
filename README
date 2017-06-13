# AST Example

This is an example code to modify AST on the fly to modify code.

Here, the code modifies other golang source code under `testdata/code.go` to add filename, function name, line number in all logging messages

```
➜ $?=0 shadyabhi/go-ast-example [ 5:05PM] % cat testcode/code.go                                                                                 master ✚ ✱ ◼
package examplecode

import (
        "sync"

        log "github.com/Sirupsen/logrus"
)

func bob() {
        log.Infof("Hello World from Bob")
}

func main() {
        log.Infof("Hello World")
        var wg sync.WaitGroup
        wg.Add(1)
        go func() {
                bob()
                wg.Done()
        }()
        wg.Wait()
}

>>>  0s elasped...
➜ $?=0 /Users/arastogi/go_workspace/src/github.com/shadyabhi/go-ast-example [ 5:05PM] % go run main.go                                                                                       master ✚ ✱ ◼
package examplecode

import (
        "sync"
        log "github.com/Sirupsen/logrus"
)

func bob() {
        log.Infof("code.go:bob:10 - Hello World from Bob")
}
func main() {
        log.Infof("code.go:main:14 - Hello World")
        var wg sync.WaitGroup
        wg.Add(1)
        go func() {
                bob()
                wg.Done()
        }()
        wg.Wait()
}

>>>  0s elasped...
➜ $?=0 shadyabhi/go-ast-example [ 5:07PM] %
```
