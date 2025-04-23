package biz

import (
	"context"
	pb "helloworld/api/file/service/v1"
	"os"
	"path/filepath"

	"github.com/go-kratos/kratos/v2/log"
)

type FileBiz struct {
	pb.UnimplementedFileServiceServer
}

func NewFileBiz() *FileBiz {
	return &FileBiz{}
}

func (b *FileBiz) Upload(_ context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {

	// create ./tmp file dir
	if err := os.MkdirAll("./tmp", os.ModePerm); err != nil {
		return nil, err
	}

	filePath := filepath.Clean(filepath.Join("./tmp", req.GetFilename()))

	// Save the file to the specified path
	err := os.WriteFile(filePath, req.GetFileData(), 0644)
	if err != nil {
		return nil, err
	}

	return &pb.UploadResponse{
		Message: "Success",
	}, nil
}

func (b *FileBiz) Download(_ context.Context, req *pb.DownloadRequest) (*pb.DownloadResponse, error) {

	filePath := filepath.Clean(filepath.Join("./tmp", req.GetFilename()))
	log.Info(filePath)

	// check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, err
	}

	// Read the file from the specified path
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	log.Info(fileData)

	return &pb.DownloadResponse{
		FileData: fileData,
	}, nil
}
