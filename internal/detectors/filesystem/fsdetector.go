package filesystem

import (
	"github.com/mlw157/Scout/internal/detectors"
	"github.com/mlw157/Scout/internal/models"
	"io/fs"
	"path/filepath"
)

type FSDetector struct {
	filePatterns []detectors.FilePattern
}

func NewFSDetector() *FSDetector {
	return &FSDetector{}
}

func (detector *FSDetector) DetectFiles(root string, excludeDirs []string, ecosystems []string) ([]models.File, error) {
	detector.populateFilePatterns(ecosystems)

	var detectedFiles []models.File

	// WalkDir will use the function visit on every directory entry found
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		return detector.visit(path, d, err, &detectedFiles, excludeDirs)
	})

	return detectedFiles, err
}

// this function will run on every directory entry, we need to pass the pointer to detectedFiles, in order to modify it
func (detector *FSDetector) visit(path string, d fs.DirEntry, err error, detectedFiles *[]models.File, excludeDirs []string) error {
	if err != nil {
		return err
	}

	//fmt.Println(" ", path, d.IsDir())

	// this skips excluded directories
	for _, exclude := range excludeDirs {
		if d.Name() == exclude {
			return filepath.SkipDir
		}
	}

	// for every not excluded file, try to match pattern expressions
	if !d.IsDir() {
		for _, pattern := range detector.filePatterns {
			if pattern.Regex.MatchString(d.Name()) {
				//fmt.Printf("File %s matched Pattern %v\n", d.Name(), pattern.Regex)
				*detectedFiles = append(*detectedFiles, models.File{
					Path:      path,
					Ecosystem: pattern.Ecosystem,
				})
			}
		}
	}

	return nil

}

// based on the desired ecosystems, will populate file patterns to be used in detection, if no ecosystems are passed, all file patterns will be used
func (detector *FSDetector) populateFilePatterns(ecosystems []string) {
	detector.filePatterns = []detectors.FilePattern{}

	if len(ecosystems) > 0 {
		for _, ecosystem := range ecosystems {
			if pattern, exists := detectors.DefaultFilePatterns[ecosystem]; exists {
				detector.filePatterns = append(detector.filePatterns, pattern)
			}
		}

		return
	}

	for _, pattern := range detectors.DefaultFilePatterns {
		detector.filePatterns = append(detector.filePatterns, pattern)
	}

}
