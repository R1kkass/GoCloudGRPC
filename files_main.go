package main

import (
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