package imports

import (
	"sync"
)

// export.go exports some type and func of golang.org/x/tools/imports

var (
	exportedGoPath   map[string]*Pkg
	exportedGoPathMu sync.RWMutex
	showLog          bool
	filterGoModule   []string
	excludeGoModule  []string
)

// ScanConfig represent scanconfig
type ScanConfig struct {
	ScanRoot    bool
	PrintLog    bool
	Filter      []string
	ExecludeMod []string
}

// GoPath returns all importable packages (abs dir path => *Pkg).
func GoPath(cfg ScanConfig) map[string]*Pkg {
	showLog = cfg.PrintLog
	filterGoModule = append(filterGoModule, cfg.Filter...)
	excludeGoModule = append(excludeGoModule, cfg.ExecludeMod...)

	exportedGoPathMu.Lock()
	defer exportedGoPathMu.Unlock()
	if exportedGoPath != nil {
		return exportedGoPath
	}
	populateIgnoreOnce.Do(populateIgnore)
	/* scanGoRootOnce.Do(scanGoRoot) // async */
	/* scanGoPathOnce.Do(scanGoPath) */
	scanGoModPathOnce.Do(scanGoMod)
	if cfg.ScanRoot {
		scanGoRootOnce.Do(scanGoRoot) // async
		<-scanGoRootDone
	}
	dirScanMu.Lock()
	defer dirScanMu.Unlock()
	exportedGoPath = exportDirScan(dirScan)
	return exportedGoPath
}

func exportDirScan(ds map[string]*pkg) map[string]*Pkg {
	r := make(map[string]*Pkg)
	for path, pkg := range ds {
		r[path] = exportPkg(pkg)
	}
	return r
}

// Pkg represents exported type of pkg.
type Pkg struct {
	Dir             string // absolute file path to Pkg directory ("/usr/lib/go/src/net/http")
	ImportPath      string // full Pkg import path ("net/http", "foo/bar/vendor/a/b")
	ImportPathShort string // vendorless import path ("net/http", "a/b")
}

func exportPkg(p *pkg) *Pkg {
	return &Pkg{Dir: p.dir, ImportPath: p.importPath, ImportPathShort: p.importPathShort}
}
