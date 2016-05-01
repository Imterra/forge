package files

import (
	"../../proto"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

type File struct {
	Rootdir string
}

func (f *File) SendFile(args *proto.FileRequest, resp *proto.File) error {
	full_path := filepath.Join(f.Rootdir, args.Filename)

	data, err := ioutil.ReadFile(full_path)
	if err != nil {
		return err
	}

	fi, err := os.Stat(full_path)
	if err != nil {
		return err
	}

	resp.Filename = args.Filename
	resp.Content = data
	resp.Mode = fi.Mode()
	return nil
}

func (f *File) RecvFile(args *proto.File, resp *proto.FileResponse) error {
	full_path := filepath.Join(f.Rootdir, args.Filename)
	full_dir := filepath.Dir(full_path)

	var mode os.FileMode = os.ModeDir + 0755
	err := os.MkdirAll(full_dir, mode)
	if err != nil {
		return errors.New("Creating directory: " + err.Error())
	}

	_, err = os.Stat(full_path)

	resp.Overwritten = (err != nil)

	err = ioutil.WriteFile(full_path, args.Content, args.Mode)
	if err != nil {
		return err
	}

	resp.Filename = args.Filename
	resp.Checksum = GetDataChecksum(args.Content)
	return nil
}
