package version

import "runtime/debug"

var version = "dev"

func Full() string {
    if version != "dev" { 
        return version
    }
    if bi, ok := debug.ReadBuildInfo(); ok {
        v := bi.Main.Version
        if v != "" && v != "(devel)" {
            return v 
        }
    }
    return version
}