package lib

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
)

// ExtractTar extracts the named tar file into
// the directory named by target.
func ExtractTar(f string, target string) error {
	tf, err := os.Open(f)
	if err != nil {
		return err
	}
	// parse the file as a tar
	tr := tar.NewReader(tf)
	var h *tar.Header

	// loop over the files/folders in the
	// tar one at a time
	for h, err = tr.Next(); err == nil; h, err = tr.Next() {
		fi := h.FileInfo()
		path := filepath.Join(target, h.Name)
		if fi.Mode()&os.ModeSymlink != 0 {
			// it's a symlink; create it
			lpath := h.Linkname
			if !filepath.IsAbs(lpath) {
				lpath = filepath.Join(target, lpath)
			}
			err := os.Symlink(lpath, path)
			if err != nil {
				return err
			}
		} else if fi.IsDir() {
			// use os.MkdirAll instead of os.Mkdir
			// in case the directory already exists
			// (especially since many tar files
			// include ./)
			err := os.MkdirAll(path, fi.Mode())
			if err != nil {
				return err
			}
			// just in case the user's umask
			// masked out some of the mode bits
			err = os.Chmod(path, fi.Mode())
			if err != nil {
				return err
			}
		} else {
			// it's a normal file; create it with
			// the appropriate mode bits and write
			// the contents to it
			f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fi.Mode())
			if err != nil {
				return err
			}
			_, err = io.Copy(f, tr)
			if err != nil {
				return err
			}
			err = f.Sync()
			if err != nil {
				return err
			}
			err = f.Close()
			if err != nil {
				return err
			}
			// just in case the user's umask
			// masked out some of the mode bits
			err = os.Chmod(path, fi.Mode())
			if err != nil {
				return err
			}
		}
	}
	if err != io.EOF {
		return err
	}
	return nil
}
