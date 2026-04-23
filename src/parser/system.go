package parser

import (
	"os"
	"path/filepath"
)

func ParseSystemModel(sysDir string) (SystemModel, error) {
	info, err := os.Stat(sysDir)
	if err != nil {
		return SystemModel{}, err
	}
	if !info.IsDir() {
		return SystemModel{}, &ParseError{File: sysDir, Line: 0, Message: "sys path is not a directory"}
	}

	model := SystemModel{
		Requirements: []Requirement{},
		Components:   []Component{},
		Interfaces:   []InterfaceDefinition{},
		TestSpecs:    []TestSpec{},
	}

	reqFiles, err := globIn(sysDir, "requirements", "*.req")
	if err != nil {
		return SystemModel{}, err
	}
	for _, f := range reqFiles {
		reqs, err := ParseRequirementsFile(f)
		if err != nil {
			return SystemModel{}, err
		}
		model.Requirements = append(model.Requirements, reqs...)
	}

	ifaceFiles, err := globIn(sysDir, "interfaces", "*.iface")
	if err != nil {
		return SystemModel{}, err
	}
	for _, f := range ifaceFiles {
		ifaces, err := ParseInterfacesFile(f)
		if err != nil {
			return SystemModel{}, err
		}
		model.Interfaces = append(model.Interfaces, ifaces...)
	}

	componentsDir := filepath.Join(sysDir, "components")
	if dirExists(componentsDir) {
		entries, err := os.ReadDir(componentsDir)
		if err != nil {
			return SystemModel{}, err
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			dir := filepath.Join(componentsDir, entry.Name())
			compFiles, err := filepath.Glob(filepath.Join(dir, "*.component"))
			if err != nil {
				return SystemModel{}, err
			}
			for _, f := range compFiles {
				comps, err := ParseComponentsFile(f)
				if err != nil {
					return SystemModel{}, err
				}
				model.Components = append(model.Components, comps...)
			}

			derivedFiles, err := filepath.Glob(filepath.Join(dir, "*.req"))
			if err != nil {
				return SystemModel{}, err
			}
			for _, f := range derivedFiles {
				reqs, err := ParseRequirementsFile(f)
				if err != nil {
					return SystemModel{}, err
				}
				model.Requirements = append(model.Requirements, reqs...)
			}

			compTests, err := filepath.Glob(filepath.Join(dir, "*.testspec"))
			if err != nil {
				return SystemModel{}, err
			}
			for _, f := range compTests {
				specs, err := ParseTestSpecsFile(f)
				if err != nil {
					return SystemModel{}, err
				}
				model.TestSpecs = append(model.TestSpecs, specs...)
			}
		}
	}

	sysTests, err := globIn(sysDir, "tests", "*.testspec")
	if err != nil {
		return SystemModel{}, err
	}
	for _, f := range sysTests {
		specs, err := ParseTestSpecsFile(f)
		if err != nil {
			return SystemModel{}, err
		}
		model.TestSpecs = append(model.TestSpecs, specs...)
	}

	return model, nil
}

func globIn(root, sub, pattern string) ([]string, error) {
	dir := filepath.Join(root, sub)
	if !dirExists(dir) {
		return nil, nil
	}
	return filepath.Glob(filepath.Join(dir, pattern))
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
