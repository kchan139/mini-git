package repository

import "os"

func (r *Repository) AddToIndex(path, hash string, info os.FileInfo) error {
	return nil
}
