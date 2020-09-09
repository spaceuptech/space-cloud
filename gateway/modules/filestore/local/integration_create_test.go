// +build file_integration

package local

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

func Test_CreateFile(t *testing.T) {
	ctx := context.Background()
	type args struct {
		req  *model.CreateFileRequest
		data io.Reader
	}
	type test struct {
		name         string
		WantFileName string
		WantFileData []byte
		args         args
		wantErr      bool
	}
	testCases := []test{
		{
			name: "Create a text file at root level path doesn't start with slash(/)",
			args: args{
				req: &model.CreateFileRequest{
					Path: "",
					Type: "file",
					Name: "creds.txt",
				},
				data: bytes.NewBuffer([]byte("Die always like a fantastic lieutenant commander.")),
			},
			WantFileName: "creds.txt",
			WantFileData: []byte("Die always like a fantastic lieutenant commander."),
			wantErr:      false,
		},
		{
			name: "Create a text file at root level",
			args: args{
				req: &model.CreateFileRequest{
					Path: "/",
					Type: "file",
					Name: "creds.txt",
				},
				data: bytes.NewBuffer([]byte("Die always like a fantastic lieutenant commander.")),
			},
			WantFileName: "creds.txt",
			WantFileData: []byte("Die always like a fantastic lieutenant commander."),
			wantErr:      false,
		},
		{
			name: "Create a text file in a single level nested folder where path doesn't start with slash(/)",
			args: args{
				req: &model.CreateFileRequest{
					Path:    "websites/",
					Type:    "file",
					Name:    "creds.txt",
					MakeAll: true,
				},
				data: bytes.NewBuffer([]byte("Die always like a fantastic lieutenant commander.")),
			},
			WantFileName: "creds.txt",
			WantFileData: []byte("Die always like a fantastic lieutenant commander."),
			wantErr:      false,
		},
		{
			name: "Create a text file in a single level nested folder where path doesn't end with slash(/)",
			args: args{
				req: &model.CreateFileRequest{
					Path:    "/websites",
					Type:    "file",
					Name:    "creds.txt",
					MakeAll: true,
				},
				data: bytes.NewBuffer([]byte("Die always like a fantastic lieutenant commander.")),
			},
			WantFileName: "creds.txt",
			WantFileData: []byte("Die always like a fantastic lieutenant commander."),
			wantErr:      false,
		},
		{
			name: "Create a text file in a single level nested folder",
			args: args{
				req: &model.CreateFileRequest{
					Path:    "/websites/",
					Type:    "file",
					Name:    "creds.txt",
					MakeAll: true,
				},
				data: bytes.NewBuffer([]byte("Die always like a fantastic lieutenant commander.")),
			},
			WantFileName: "creds.txt",
			WantFileData: []byte("Die always like a fantastic lieutenant commander."),
			wantErr:      false,
		},
		{
			name: "Create a text file in a single level nested folder where the folder doesn't exists",
			args: args{
				req: &model.CreateFileRequest{
					Path: "/websites",
					Type: "file",
					Name: "creds.txt",
				},
				data: bytes.NewBuffer([]byte("Die always like a fantastic lieutenant commander.")),
			},
			WantFileName: "creds.txt",
			WantFileData: []byte("Die always like a fantastic lieutenant commander."),
			wantErr:      true,
		},
	}

	path := fmt.Sprintf("%s/space_cloud_test", os.ExpandEnv("$HOME"))
	file, err := Init(path)
	if err != nil {
		t.Fatalf("Create() Couldn't initialize local store for path (%s)", path)
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			err = file.CreateFile(ctx, tt.args.req, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				data, err := ioutil.ReadFile(path + "/" + strings.TrimSuffix(strings.TrimPrefix(tt.args.req.Path, "/"), "/") + "/" + tt.args.req.Name)
				if err != nil {
					t.Errorf("Create() unable to read created file (%v)", err)
					return
				}

				if string(data) != string(tt.WantFileData) {
					t.Errorf("Create() file contains wrong data (%v)", err)
					return
				}
			}

			// clear data
			if err := RemoveContents(path); err != nil {
				helpers.Logger.LogInfo(helpers.GetRequestID(ctx), "Couldn't clean inside data generated by tests", nil)
				return
			}
		})
	}
	// clear data
	if err := os.RemoveAll(path); err != nil {
		helpers.Logger.LogInfo(helpers.GetRequestID(ctx), "Couldn't clean data generated by tests", nil)
		return
	}
}

func Test_CreateFolder(t *testing.T) {
	ctx := context.Background()
	type args struct {
		req  *model.CreateFileRequest
		data io.Reader
	}
	type test struct {
		name        string
		WantDirName string
		args        args
		wantErr     bool
	}
	testCases := []test{
		{
			name: "Create a directory at root level path doesn't start with slash(/)",
			args: args{
				req: &model.CreateFileRequest{
					Path: "",
					Type: "dir",
					Name: "websites",
				},
			},
			WantDirName: "websites",
			wantErr:     false,
		},
		{
			name: "Create a directory at root level",
			args: args{
				req: &model.CreateFileRequest{
					Path: "/",
					Type: "dir",
					Name: "websites",
				},
			},
			WantDirName: "websites",
			wantErr:     false,
		},
		{
			name: "Create a directory at single nested level path doesn't start with slash(/)",
			args: args{
				req: &model.CreateFileRequest{
					Path:    "websites/",
					Type:    "dir",
					Name:    "netlify",
					MakeAll: true,
				},
			},
			WantDirName: "netlify",
			wantErr:     false,
		},
		{
			name: "Create a directory at single nested level path doesn't end with slash(/)",
			args: args{
				req: &model.CreateFileRequest{
					Path:    "websites",
					Type:    "dir",
					Name:    "netlify",
					MakeAll: true,
				},
			},
			WantDirName: "netlify",
			wantErr:     false,
		},
		{
			name: "Create a directory at single nested level",
			args: args{
				req: &model.CreateFileRequest{
					Path:    "/websites/",
					Type:    "dir",
					Name:    "netlify",
					MakeAll: true,
				},
			},
			WantDirName: "netlify",
			wantErr:     false,
		},
		{
			name: "Create a directory at single nested level in which provided path doesn't exists",
			args: args{
				req: &model.CreateFileRequest{
					Path: "/websites/",
					Type: "dir",
					Name: "netlify",
				},
			},
			WantDirName: "netlify",
			wantErr:     true,
		},
	}

	path := fmt.Sprintf("%s/space_cloud_test", os.ExpandEnv("$HOME"))
	file, err := Init(path)
	if err != nil {
		t.Fatalf("Create() Couldn't initialize local store for path (%s)", path)
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			err = file.CreateDir(ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				_, err = ioutil.ReadDir(path + "/" + strings.TrimSuffix(strings.TrimPrefix(tt.args.req.Path, "/"), "/") + "/" + tt.args.req.Name)
				if err != nil {
					t.Errorf("Create() unable to read created dir (%v)", err)
					return
				}
			}

			// clear data
			if err := RemoveContents(path); err != nil {
				helpers.Logger.LogInfo(helpers.GetRequestID(ctx), "Couldn't clean inside data generated by tests", nil)
				return
			}
		})
	}
	// clear data
	if err := os.RemoveAll(path); err != nil {
		helpers.Logger.LogInfo(helpers.GetRequestID(ctx), "Couldn't clean data generated by tests", nil)
		return
	}
}

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
