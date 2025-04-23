package service

import (
	"context"
	"errors"
	"helloworld/internal/biz"

	pb "helloworld/api/file/service/v1"

	"github.com/go-kratos/kratos/v2/log"
)

type FileServiceService struct {
	pb.UnimplementedFileServiceServer
	file *biz.FileBiz
}

func NewFileServiceService() *FileServiceService {
	return &FileServiceService{}
}

func (s *FileServiceService) UploadFile(ctx context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {

	log.Infof("UploadFile: %s", req.GetFilename())
	log.Infof("FileData: %s", req.GetFileData()[:10])

	if req.GetFileData() == nil || len(req.GetFileData()) == 0 {
		return nil, errors.New("file data is required")
	}

	uploadResp, err := s.file.Upload(ctx, req)
	if err != nil {
		return nil, errors.New("upload file failed")
	}

	return uploadResp, nil
}
func (s *FileServiceService) DownloadFile(ctx context.Context, req *pb.DownloadRequest) (*pb.DownloadResponse, error) {

	log.Infof("DownloadFile: %s", req.GetFilename())

	if req.GetFilename() == "" {
		return nil, errors.New("filename is required")
	}

	downloadResponse, err := s.file.Download(ctx, req)
	if err != nil {
		return nil, errors.New("download file failed")
	}
	log.Infof("DownloadFile: %s", downloadResponse.GetFileData()[:10])

	return downloadResponse, nil
}
