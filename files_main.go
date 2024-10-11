package main

import (
	"context"
	"mypackages/controllers"
	"mypackages/proto/files"
)

type filesServer struct {
	*files.UnimplementedFilesGreeterServer
}

func (s *filesServer) DownloadFile(in *files.FileDownloadRequest, responseStream files.FilesGreeter_DownloadFileServer) error {
	return controllers.DownloadFile(in, responseStream)
}

func (s *filesServer) UploadFile(stream files.FilesGreeter_UploadFileServer) error {
	return controllers.UploadFile(stream)
}

func (s *filesServer) FindFile(context context.Context, in *files.FindFileRequest) (*files.FindFileResponse, error) {
	return controllers.FindFile(context, in)
}