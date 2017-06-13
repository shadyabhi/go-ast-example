# AST Example

This is an example code to modify AST on the fly to modify code.

Here, the code modifies other golang source code under `testdata/code.go` to add filename, function name, line number in all logging messages

```
➜ $?=1 shadyabhi/go-ast-example [ 5:18PM] % cat testcode/code.go                                                                                 ⏎ master ✚ ◼
package examplecode

import (
        "sync"

        log "github.com/Sirupsen/logrus"
)

func bob() {
        log.Debugf("Hello World from Bob: %d", 0)
}

func main() {
        log.Infof("Hello World from main")
        var wg sync.WaitGroup
        wg.Add(1)
        go func() {
                bob()
                wg.Done()
        }()
        wg.Wait()
}

>>>  1s elasped...
➜ $?=0 /Users/arastogi/go_workspace/src/github.com/shadyabhi/go-ast-example [ 5:18PM] % go run main.go                                                                                         master ✚ ◼
package examplecode

import (
        "sync"
        log "github.com/Sirupsen/logrus"
)

func bob() {
        log.Debugf("code.go:bob:10 - Hello World from Bob: %d", 0)
}
func main() {
        log.Infof("code.go:main:14 - Hello World from main")
        var wg sync.WaitGroup
        wg.Add(1)
        go func() {
                bob()
                wg.Done()
        }()
        wg.Wait()
}

>>>  0s elasped...
➜ $?=0 shadyabhi/go-ast-example [ 5:18PM] %
```
