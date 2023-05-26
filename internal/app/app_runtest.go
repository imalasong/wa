// 版权 @2023 凹语言 作者。保留所有权利。

package app

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/tetratelabs/wazero/sys"
	"wa-lang.org/wa/internal/app/apputil"
	"wa-lang.org/wa/internal/backends/compiler_wat"
	"wa-lang.org/wa/internal/loader"
)

func (p *App) RunTest(pkgpath string, appArgs ...string) error {
	cfg := p.opt.Config()
	cfg.UnitTest = true
	prog, err := loader.LoadProgram(cfg, pkgpath)
	if err != nil {
		return err
	}

	startTime := time.Now()
	mainPkg := prog.Pkgs[prog.Manifest.MainPkg]

	if len(mainPkg.TestInfo.Files) == 0 {
		fmt.Printf("?    %s [no test files]\n", prog.Manifest.MainPkg)
		return nil
	}

	var firstError error
	for _, t := range mainPkg.TestInfo.Tests {
		output, err := compiler_wat.New().Compile(prog, t.Name)
		if err != nil {
			return err
		}

		if err = os.WriteFile("a.out.wat", []byte(output), 0666); err != nil {
			return err
		}

		stdout, stderr, err := apputil.RunWasmEx(cfg, "a.out.wat", appArgs...)

		stdout = bytes.TrimSpace(stdout)
		bOutputOK := t.Output == string(stdout)

		if err == nil && bOutputOK {
			continue
		}

		if err != nil {
			if firstError == nil {
				firstError = err
			}
			if _, ok := err.(*sys.ExitError); ok {
				fmt.Printf("---- %s.%s\n", prog.Manifest.MainPkg, t.Name)
				if s := sWithPrefix(string(stdout), "    "); s != "" {
					fmt.Println(s)
				}
				if s := sWithPrefix(string(stderr), "    "); s != "" {
					fmt.Println(s)
				}
			} else {
				fmt.Println(err)
			}
		}

		if t.Output != "" {
			if expect, got := t.Output, string(stdout); expect != got {
				if firstError == nil {
					firstError = fmt.Errorf("expect = %q, got = %q", expect, got)
				}
				fmt.Printf("---- %s.%s\n", prog.Manifest.MainPkg, t.Name)
				fmt.Printf("    expect = %q, got = %q\n", expect, got)
			}
		}
	}
	for _, t := range mainPkg.TestInfo.Examples {
		output, err := compiler_wat.New().Compile(prog, t.Name)
		if err != nil {
			return err
		}

		if err = os.WriteFile("a.out.wat", []byte(output), 0666); err != nil {
			return err
		}

		stdout, stderr, err := apputil.RunWasmEx(cfg, "a.out.wat", appArgs...)

		stdout = bytes.TrimSpace(stdout)
		bOutputOK := t.Output == string(stdout)

		if err == nil && bOutputOK {
			continue
		}

		if err != nil {
			if firstError == nil {
				firstError = err
			}
			if _, ok := err.(*sys.ExitError); ok {
				fmt.Printf("---- %s.%s\n", prog.Manifest.MainPkg, t.Name)
				if s := sWithPrefix(string(stdout), "    "); s != "" {
					fmt.Println(s)
				}
				if s := sWithPrefix(string(stderr), "    "); s != "" {
					fmt.Println(s)
				}
			} else {
				fmt.Println(err)
			}
		}

		if t.Output != "" {
			if expect, got := t.Output, string(stdout); expect != got {
				if firstError == nil {
					firstError = fmt.Errorf("expect = %q, got = %q", expect, got)
				}
				fmt.Printf("---- %s.%s\n", prog.Manifest.MainPkg, t.Name)
				fmt.Printf("    expect = %q, got = %q\n", expect, got)
			}
		}
	}
	if firstError != nil {
		fmt.Printf("FAIL %s %v\n", prog.Manifest.MainPkg, time.Now().Sub(startTime).Round(time.Millisecond))
		os.Exit(1)
	}

	fmt.Printf("ok   %s %v\n", prog.Manifest.MainPkg, time.Now().Sub(startTime).Round(time.Millisecond))

	return nil
}

func sWithPrefix(s, prefix string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for i, line := range lines {
		lines[i] = prefix + strings.TrimSpace(line)
	}
	return strings.Join(lines, "\n")
}
