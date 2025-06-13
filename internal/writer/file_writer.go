// internal/writer/file_writer.go
package writer

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

// FileWriter はHCLコンテンツをファイルに書き込みます。
type FileWriter struct{}

// NewFileWriter は新しいFileWriterを生成します。
func NewFileWriter() *FileWriter {
	return &FileWriter{}
}

// WriteFile は指定されたパスにHCLファイルの内容を書き込みます。
func (w *FileWriter) WriteFile(path string, file *hclwrite.File) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", path, err)
	}
	defer f.Close()

	_, err = f.Write(file.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %w", path, err)
	}
	return nil
} 
