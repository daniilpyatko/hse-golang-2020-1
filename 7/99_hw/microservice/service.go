package main

import (
	context "context"
	"encoding/json"
	"errors"
	"net"
	"regexp"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	status "google.golang.org/grpc/status"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные

type Service struct {
	Allowed     map[string][]string
	Addr        string
	WriteTo     []chan Event
	WriteStatTo []chan Stat
	ClientAddr  []string
	NameToAddr  map[string]string
}

func (s *Service) authUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	curAr := md.Get("consumer")
	if len(curAr) == 0 {
		return nil, status.Error(codes.Unauthenticated, "")
	}
	consumer := md.Get("consumer")[0]
	curMethod := info.FullMethod
	matched := false
	for _, cur := range s.Allowed[consumer] {
		if ok, _ := regexp.MatchString(cur, curMethod); ok {
			matched = true
		}
	}
	if !matched {
		return nil, status.Error(codes.Unauthenticated, "")
	}
	reply, err := handler(ctx, req)
	return reply, err
}

func (s *Service) authStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	md, _ := metadata.FromIncomingContext(ss.Context())
	curAr := md.Get("consumer")
	if len(curAr) == 0 {
		return status.Error(codes.Unauthenticated, "")
	}
	consumer := md.Get("consumer")[0]
	curMethod := info.FullMethod
	matched := false
	for _, cur := range s.Allowed[consumer] {
		if ok, _ := regexp.MatchString(cur, curMethod); ok {
			matched = true
		}
	}
	if !matched {
		return status.Error(codes.Unauthenticated, "")
	}
	err := handler(srv, ss)
	return err
}

func (s *Service) Start(ctx context.Context, lis net.Listener, service *Service) error {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(s.authUnaryInterceptor),
		grpc.StreamInterceptor(s.authStreamInterceptor),
	)
	RegisterBizServer(server, NewBiz(service))
	RegisterAdminServer(server, NewAdmin(service))
	go server.Serve(lis)
	defer func() {
		server.Stop()
		lis.Close()
	}()
	<-ctx.Done()
	return nil
}

func StartMyMicroservice(ctx context.Context, addr string, ACLData string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	var allowed map[string][]string
	err = json.Unmarshal([]byte(ACLData), &allowed)
	if err != nil {
		lis.Close()
		return errors.New("Couldn't unmarshal ACLData")
	}
	service := Service{
		Allowed: allowed,
		Addr:    addr,
		NameToAddr: map[string]string{
			"Logging":    "/main.Admin/Logging",
			"Check":      "/main.Biz/Check",
			"Test":       "/main.Biz/Test",
			"Add":        "/main.Biz/Add",
			"Statistics": "/main.Admin/Statistics",
		},
	}
	go service.Start(ctx, lis, &service)
	return nil
}

type Admin struct {
	service *Service
}

func NewAdmin(s *Service) *Admin {
	return &Admin{service: s}
}

func (s *Service) Log(from string, consumer string) {
	writeString := s.NameToAddr[from]

	// Writing to Log
	for ind, ch := range s.WriteTo {
		ev := Event{}
		ev.Consumer = consumer
		ev.Method = writeString
		ev.Host = s.ClientAddr[ind]
		ch <- ev
	}

	// Writing to Stat
	newStat := Stat{
		ByMethod: map[string]uint64{
			writeString: 1,
		},
		ByConsumer: map[string]uint64{
			consumer: 1,
		},
	}
	for _, ch := range s.WriteStatTo {
		ch <- newStat
	}

}

func (a *Admin) Logging(nothing *Nothing, stream Admin_LoggingServer) error {
	md, _ := metadata.FromIncomingContext(stream.Context())
	c, _ := peer.FromContext(stream.Context())
	a.service.Log("Logging", md.Get("consumer")[0])
	cur := make(chan Event)
	a.service.WriteTo = append(a.service.WriteTo, cur)
	a.service.ClientAddr = append(a.service.ClientAddr, c.Addr.String())
	for {
		curEvent := <-cur
		stream.Send(&curEvent)
	}
	return nil
}

func (a *Admin) Statistics(interval *StatInterval, stream Admin_StatisticsServer) error {
	md, _ := metadata.FromIncomingContext(stream.Context())
	a.service.Log("Statistics", md.Get("consumer")[0])
	tm := time.Duration(interval.GetIntervalSeconds())
	ticker := time.NewTicker(tm * time.Second)
	curStat := Stat{
		ByConsumer: make(map[string]uint64),
		ByMethod:   make(map[string]uint64),
	}
	ch := make(chan Stat)
	a.service.WriteStatTo = append(a.service.WriteStatTo, ch)
	for {
		select {
		case <-ticker.C:
			stream.Send(&curStat)
			curStat = Stat{
				ByConsumer: make(map[string]uint64),
				ByMethod:   make(map[string]uint64),
			}
		case newStat := <-ch:
			for k, v := range newStat.ByConsumer {
				curStat.ByConsumer[k] += v
			}
			for k, v := range newStat.ByMethod {
				curStat.ByMethod[k] += v
			}
		}
	}
	return nil
}

type Biz struct {
	service *Service
}

func NewBiz(s *Service) *Biz {
	return &Biz{service: s}
}

func (b *Biz) Check(ctx context.Context, nothing *Nothing) (*Nothing, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	b.service.Log("Check", md.Get("consumer")[0])
	return &Nothing{}, nil
}
func (b *Biz) Add(ctx context.Context, nothing *Nothing) (*Nothing, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	b.service.Log("Add", md.Get("consumer")[0])
	return &Nothing{}, nil
}
func (b *Biz) Test(ctx context.Context, nothing *Nothing) (*Nothing, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	b.service.Log("Test", md.Get("consumer")[0])
	return &Nothing{}, nil
}
