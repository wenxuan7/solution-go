package main

import (
	"context"
	"fmt"
	settingsproto "github.com/wenxuan7/oms-grpc/settings"
	"github.com/wenxuan7/solution/settings"
	"github.com/wenxuan7/solution/utils"
)

type SettingsReaderServer struct {
	s *settings.Service
	settingsproto.UnimplementedReaderServer
}

func NewSettingsReaderServer() *SettingsReaderServer {
	s, err := settings.NewServiceWithLCache()
	if err != nil {
		panic(err)
	}
	return &SettingsReaderServer{s: s}
}

func (srs *SettingsReaderServer) Get(ctx context.Context, req *settingsproto.Req) (*settingsproto.Resp, error) {
	ctx = utils.WithCompanyId(ctx, uint(req.CompanyId))
	ctx = utils.WithTraceId(ctx, req.TraceId)
	get, err := srs.s.Get(ctx, req.GetK())
	if err != nil {
		return nil, fmt.Errorf("server: fail to settings.Get in SettingsReaderServer: %w", err)
	}
	resp := &settingsproto.Resp{V: get}
	return resp, nil
}

func (srs *SettingsReaderServer) Gets(ctx context.Context, req *settingsproto.MultiReq) (*settingsproto.MultiResp, error) {
	ctx = utils.WithCompanyId(ctx, uint(req.CompanyId))
	gets, err := srs.s.Gets(ctx, req.Ks)
	if err != nil {
		return nil, fmt.Errorf("server: fail to settings.Gets in SettingsReaderServer: %w", err)
	}
	resp := &settingsproto.MultiResp{Vs: gets}
	return resp, nil
}
