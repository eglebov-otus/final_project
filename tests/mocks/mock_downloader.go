// Code generated by MockGen. DO NOT EDIT.
// Source: image-previewer/internal/domain (interfaces: Downloader)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	image "image"
	dto "image-previewer/internal/domain/dto"
	reflect "reflect"
)

// MockDownloader is a mock of Downloader interface
type MockDownloader struct {
	ctrl     *gomock.Controller
	recorder *MockDownloaderMockRecorder
}

// MockDownloaderMockRecorder is the mock recorder for MockDownloader
type MockDownloaderMockRecorder struct {
	mock *MockDownloader
}

// NewMockDownloader creates a new mock instance
func NewMockDownloader(ctrl *gomock.Controller) *MockDownloader {
	mock := &MockDownloader{ctrl: ctrl}
	mock.recorder = &MockDownloaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDownloader) EXPECT() *MockDownloaderMockRecorder {
	return m.recorder
}

// Download mocks base method
func (m *MockDownloader) Download(arg0 string, arg1 dto.ImageDimensions) (image.Image, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Download", arg0, arg1)
	ret0, _ := ret[0].(image.Image)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Download indicates an expected call of Download
func (mr *MockDownloaderMockRecorder) Download(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Download", reflect.TypeOf((*MockDownloader)(nil).Download), arg0, arg1)
}
