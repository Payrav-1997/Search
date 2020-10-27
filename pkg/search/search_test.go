package search

import (
	"context"
	"testing"
)

func TestAll(t *testing.T) {
	root := context.Background()
	files := make([]string, 0)
	files = append(files, "text1.txt")
	files = append(files, "text2.txt")
	files = append(files, "text3.txt")

	results := All(root, "concurrency", files)

	if len(results) != 0 {
		t.Error("Error")
	}
}

func TestAny(t *testing.T) {
	root := context.Background()
	files := make([]string, 0)
	files = append(files, "text1.txt")
	files = append(files, "text2.txt")
	files = append(files, "text3.txt")
	
		results := Any(root, "Tests", files)
	
		if (len(results))>1{
			t.Error("Error")
		}
}