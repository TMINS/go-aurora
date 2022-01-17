package services

import "context"

type TestService struct {
}

func (t TestService) TestPRC(ctx context.Context, request *TestRequest) (*TestResponse, error) {
	//TODO implement me
	return &TestResponse{Status: "ok"}, nil
}

func (t TestService) mustEmbedUnimplementedTestServiceServer() {
	//TODO implement me
	panic("implement me")
}
