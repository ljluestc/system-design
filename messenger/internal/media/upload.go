package media

import (
    "errors"
    "github.com/google/uuid"
    "messenger/internal/utils"
)

type Service struct {
    bucket string
}

func NewService(bucket string) *Service {
    return &Service{bucket: bucket}
}

func (s *Service) UploadFile(file []byte, fileType string) (string, error) {
    // Validate file
    if len(file) == 0 {
        return "", errors.New("empty file")
    }
    fileID := uuid.New().String()
    // Pseudocode: Upload to S3
    err := utils.UploadToS3(s.bucket, fileID, file)
    if err != nil {
        return "", err
    }
    return fileID, nil
}