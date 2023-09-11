package web

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go-basic/webbook/internal/service"
	svcmock "go-basic/webbook/internal/service/mocks"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_SignUp(t *testing.T) {
	const signupUrl = "/users/signup"
	testCases := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		reqBuilder func(t *testing.T) *http.Request
		wantCode   int
		wantBody   string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userService := svcmock.NewMockUserService(ctrl)
				// 传递任意参数都可以， 成功返回nil
				userService.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil)
				// 因为没有使用到 codeService， 所以可以初始化也可以不初始化
				codeService := svcmock.NewMockCodeService(ctrl)
				return userService, codeService
			},
			reqBuilder: func(t *testing.T) *http.Request {
				// 构造请求体数据
				reqBody := bytes.NewBuffer([]byte(`{"email": "123@qq.com","password": "hello@world123","confirmPassword": "hello@world123"}`))
				request, err := http.NewRequest(http.MethodPost, signupUrl, reqBody)
				request.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return request
			},
			wantCode: 200,
			wantBody: "注册成功",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			// mock 需要的service
			userService, codeService := tc.mock(ctrl)
			// 使用mock的service来构造 userHandler
			userHandler := NewUserHandler(userService, codeService)
			// 注册路由
			server := gin.Default()
			userHandler.RegisterRoutes(server)
			// 构造http请求
			request := tc.reqBuilder(t)
			// 准备记录响应
			recorder := httptest.NewRecorder()
			// 执行
			server.ServeHTTP(recorder, request)
			// 断言
			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantBody, recorder.Body.String())
		})
	}
}
