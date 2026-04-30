package storage

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Storage struct {
	BasePath string
}

func New(basePath string) *Storage {
	return &Storage{BasePath: basePath}
}

func (s *Storage) NewAvatarFilename(mediaType string) string {
	buff := make([]byte, 32)
	_, err := rand.Read(buff)

	if err != nil {
		panic(err)
	}

	id := base64.RawURLEncoding.EncodeToString(buff)
	ext := mediaTypeToExt(mediaType)
	return id + ext
}

func (s *Storage) AvatarPath(filename string) string {
	return filepath.Join(s.BasePath, "avatars", filename)
}

func (s *Storage) AvatarURL(filepath string) string {
	return "/" + filepath
}

func (s *Storage) SaveFile(path string, src io.Reader) error {
	dst, err := os.Create(path)

	if err != nil {
		return err
	}

	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func (s *Storage) DeleteFile(path string) error {
	err := os.Remove(path)

	if os.IsNotExist(err) {
		return nil
	}

	return err
}

func mediaTypeToExt(mediaType string) string {
	parts := strings.Split(mediaType, "/")
	if len(parts) != 2 {
		return ".bin"
	}
	return "." + parts[1]
}
