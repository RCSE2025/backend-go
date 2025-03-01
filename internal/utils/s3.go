package utils

type S3WorkerAPI struct {
	BucketName  string
	S3WorkerURL string
}

func NewS3WorkerAPI(bucketName string, s3WorkerURL string) *S3WorkerAPI {
	return &S3WorkerAPI{
		BucketName:  bucketName,
		S3WorkerURL: s3WorkerURL,
	}
}

func (s *S3WorkerAPI) NewBucket() error {
	return nil
}

func (s *S3WorkerAPI) UploadFile(filename string) error {
	return nil
}
