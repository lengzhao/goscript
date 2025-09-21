// Hello world script
package main

import "strings"

func main() {
    s := "Hello, World!"
    if strings.HasPrefix(s, "Hello") {
        return "YES"
    } else {
        return "NO"
    }
}