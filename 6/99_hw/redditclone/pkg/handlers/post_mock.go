// Code generated by MockGen. DO NOT EDIT.
// Source: post.go

// Package handlers is a generated GoMock package.
package handlers

import (
	gomock "github.com/golang/mock/gomock"
	comment "redditclone/pkg/comment"
	post "redditclone/pkg/post"
	user "redditclone/pkg/user"
	reflect "reflect"
)

// MockPostRepo is a mock of PostRepo interface
type MockPostRepo struct {
	ctrl     *gomock.Controller
	recorder *MockPostRepoMockRecorder
}

// MockPostRepoMockRecorder is the mock recorder for MockPostRepo
type MockPostRepoMockRecorder struct {
	mock *MockPostRepo
}

// NewMockPostRepo creates a new mock instance
func NewMockPostRepo(ctrl *gomock.Controller) *MockPostRepo {
	mock := &MockPostRepo{ctrl: ctrl}
	mock.recorder = &MockPostRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPostRepo) EXPECT() *MockPostRepoMockRecorder {
	return m.recorder
}

// GetAll mocks base method
func (m *MockPostRepo) GetAll() ([]*post.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]*post.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll
func (mr *MockPostRepoMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockPostRepo)(nil).GetAll))
}

// GetByCategory mocks base method
func (m *MockPostRepo) GetByCategory(category string) ([]*post.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByCategory", category)
	ret0, _ := ret[0].([]*post.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByCategory indicates an expected call of GetByCategory
func (mr *MockPostRepoMockRecorder) GetByCategory(category interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByCategory", reflect.TypeOf((*MockPostRepo)(nil).GetByCategory), category)
}

// GetById mocks base method
func (m *MockPostRepo) GetById(id string) (*post.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", id)
	ret0, _ := ret[0].(*post.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById
func (mr *MockPostRepoMockRecorder) GetById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockPostRepo)(nil).GetById), id)
}

// GetByUsername mocks base method
func (m *MockPostRepo) GetByUsername(username string) ([]*post.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUsername", username)
	ret0, _ := ret[0].([]*post.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUsername indicates an expected call of GetByUsername
func (mr *MockPostRepoMockRecorder) GetByUsername(username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUsername", reflect.TypeOf((*MockPostRepo)(nil).GetByUsername), username)
}

// AddPost mocks base method
func (m *MockPostRepo) AddPost(post *post.Post) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddPost", post)
}

// AddPost indicates an expected call of AddPost
func (mr *MockPostRepoMockRecorder) AddPost(post interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPost", reflect.TypeOf((*MockPostRepo)(nil).AddPost), post)
}

// AddNewPost mocks base method
func (m *MockPostRepo) AddNewPost(user *user.User, newPost *post.Post) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddNewPost", user, newPost)
}

// AddNewPost indicates an expected call of AddNewPost
func (mr *MockPostRepoMockRecorder) AddNewPost(user, newPost interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNewPost", reflect.TypeOf((*MockPostRepo)(nil).AddNewPost), user, newPost)
}

// AddCommentToPost mocks base method
func (m *MockPostRepo) AddCommentToPost(postId string, curUser *user.User, curComment comment.Comment) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddCommentToPost", postId, curUser, curComment)
}

// AddCommentToPost indicates an expected call of AddCommentToPost
func (mr *MockPostRepoMockRecorder) AddCommentToPost(postId, curUser, curComment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCommentToPost", reflect.TypeOf((*MockPostRepo)(nil).AddCommentToPost), postId, curUser, curComment)
}

// PostUpvote mocks base method
func (m *MockPostRepo) PostUpvote(postId, curUserId string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PostUpvote", postId, curUserId)
}

// PostUpvote indicates an expected call of PostUpvote
func (mr *MockPostRepoMockRecorder) PostUpvote(postId, curUserId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostUpvote", reflect.TypeOf((*MockPostRepo)(nil).PostUpvote), postId, curUserId)
}

// PostDownvote mocks base method
func (m *MockPostRepo) PostDownvote(postId, curUserId string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PostDownvote", postId, curUserId)
}

// PostDownvote indicates an expected call of PostDownvote
func (mr *MockPostRepoMockRecorder) PostDownvote(postId, curUserId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostDownvote", reflect.TypeOf((*MockPostRepo)(nil).PostDownvote), postId, curUserId)
}

// PostUnvote mocks base method
func (m *MockPostRepo) PostUnvote(postId, curUserId string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PostUnvote", postId, curUserId)
}

// PostUnvote indicates an expected call of PostUnvote
func (mr *MockPostRepoMockRecorder) PostUnvote(postId, curUserId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostUnvote", reflect.TypeOf((*MockPostRepo)(nil).PostUnvote), postId, curUserId)
}

// CommentDelete mocks base method
func (m *MockPostRepo) CommentDelete(postId, commentId, curUserId string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CommentDelete", postId, commentId, curUserId)
}

// CommentDelete indicates an expected call of CommentDelete
func (mr *MockPostRepoMockRecorder) CommentDelete(postId, commentId, curUserId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommentDelete", reflect.TypeOf((*MockPostRepo)(nil).CommentDelete), postId, commentId, curUserId)
}

// PostDelete mocks base method
func (m *MockPostRepo) PostDelete(postId string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PostDelete", postId)
}

// PostDelete indicates an expected call of PostDelete
func (mr *MockPostRepoMockRecorder) PostDelete(postId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostDelete", reflect.TypeOf((*MockPostRepo)(nil).PostDelete), postId)
}
