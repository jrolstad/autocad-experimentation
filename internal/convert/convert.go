package convert

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type DependencyStatus struct {
	Name  string
	Found bool
	Path  string
}

func Doctor() []DependencyStatus {
	var out []DependencyStatus
	for _, name := range []string{"dwg2dxf", "dwgread"} {
		path, err := exec.LookPath(name)
		if err != nil {
			out = append(out, DependencyStatus{Name: name, Found: false})
			continue
		}
		out = append(out, DependencyStatus{Name: name, Found: true, Path: path})
	}
	return out
}

func ToDXF(inDWG string, verbose bool) (string, func(), error) {
	if strings.ToLower(filepath.Ext(inDWG)) != ".dwg" {
		return "", func() {}, errors.New("input must be a .dwg file")
	}
	if _, err := exec.LookPath("dwg2dxf"); err != nil {
		return "", func() {}, errors.New("dwg2dxf is not installed or not on PATH; run `dwgsum doctor`")
	}

	tmpDir, err := os.MkdirTemp("", "dwgsum-*")
	if err != nil {
		return "", func() {}, fmt.Errorf("failed creating temp dir: %w", err)
	}

	base := strings.TrimSuffix(filepath.Base(inDWG), filepath.Ext(inDWG))
	outDXF := filepath.Join(tmpDir, base+".dxf")
	cmd := exec.Command("dwg2dxf", "-o", outDXF, inDWG)
	out, err := cmd.CombinedOutput()
	if err != nil {
		cleanup := func() { _ = os.RemoveAll(tmpDir) }
		return "", cleanup, fmt.Errorf("dwg2dxf failed: %w; output: %s", err, strings.TrimSpace(string(out)))
	}
	if verbose && len(out) > 0 {
		fmt.Fprintln(os.Stderr, strings.TrimSpace(string(out)))
	}

	if _, err := os.Stat(outDXF); err != nil {
		cleanup := func() { _ = os.RemoveAll(tmpDir) }
		return "", cleanup, fmt.Errorf("conversion did not produce DXF: %w", err)
	}

	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return outDXF, cleanup, nil
}
