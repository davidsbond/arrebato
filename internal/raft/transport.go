package raft

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/hashicorp/raft"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/davidsbond/arrebato/internal/clientinfo"
	raftsvc "github.com/davidsbond/arrebato/internal/proto/arrebato/raft/service/v1"
	raftpb "github.com/davidsbond/arrebato/internal/proto/arrebato/raft/v1"
)

type (
	Transport struct {
		rpcs           chan raft.RPC
		localAddr      raft.ServerAddress
		localID        raft.ServerID
		peerMutex      *sync.RWMutex
		peers          map[raft.ServerID]*grpc.ClientConn
		tls            *tls.Config
		timeout        time.Duration
		heartbeatMutex *sync.Mutex
		heartbeatFunc  func(rpc raft.RPC)
	}

	TransportConfig struct {
		ServerAddress raft.ServerAddress
		ServerID      raft.ServerID
		TLS           *tls.Config
		RPCChannel    chan raft.RPC
		Timeout       time.Duration
	}
)

func NewTransport(config TransportConfig) *Transport {
	return &Transport{
		rpcs:           config.RPCChannel,
		localAddr:      config.ServerAddress,
		localID:        config.ServerID,
		peerMutex:      &sync.RWMutex{},
		heartbeatMutex: &sync.Mutex{},
		peers:          make(map[raft.ServerID]*grpc.ClientConn),
		tls:            config.TLS,
		timeout:        config.Timeout,
	}
}

func (t *Transport) Consumer() <-chan raft.RPC {
	return t.rpcs
}

func (t *Transport) LocalAddr() raft.ServerAddress {
	return t.localAddr
}

func (t *Transport) AppendEntriesPipeline(_ raft.ServerID, _ raft.ServerAddress) (raft.AppendPipeline, error) {
	// This transport does not support pipelining. The underlying gRPC client reuses connections (keep-alive)
	// which results in only microsecond differences.
	return nil, raft.ErrPipelineReplicationNotSupported
}

func (t *Transport) AppendEntries(id raft.ServerID, target raft.ServerAddress, args *raft.AppendEntriesRequest, resp *raft.AppendEntriesResponse) error {
	client, err := t.peer(id, target)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	ctx = protocolVersionToMetadata(ctx, args.RPCHeader.ProtocolVersion)
	defer cancel()

	result, err := client.AppendEntries(ctx, &raftsvc.AppendEntriesRequest{
		Term:              args.Term,
		Leader:            args.Leader,
		PrevLogEntry:      args.PrevLogEntry,
		PrevLogTerm:       args.PrevLogTerm,
		Entries:           logsToProto(args.Entries),
		LeaderCommitIndex: args.LeaderCommitIndex,
	})
	if err != nil {
		return err
	}

	*resp = raft.AppendEntriesResponse{
		RPCHeader:      args.RPCHeader,
		Term:           result.GetTerm(),
		LastLog:        result.GetLastLog(),
		Success:        result.GetSuccess(),
		NoRetryBackoff: result.GetNoRetryBackoff(),
	}

	return nil
}

func (t *Transport) RequestVote(id raft.ServerID, target raft.ServerAddress, args *raft.RequestVoteRequest, resp *raft.RequestVoteResponse) error {
	client, err := t.peer(id, target)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	ctx = protocolVersionToMetadata(ctx, args.RPCHeader.ProtocolVersion)
	defer cancel()

	result, err := client.RequestVote(ctx, &raftsvc.RequestVoteRequest{
		Term:               args.Term,
		Candidate:          args.Candidate,
		LastLogIndex:       args.LastLogIndex,
		LastLogTerm:        args.LastLogTerm,
		LeadershipTransfer: args.LeadershipTransfer,
	})
	if err != nil {
		return err
	}

	*resp = raft.RequestVoteResponse{
		RPCHeader: args.RPCHeader,
		Term:      result.GetTerm(),
		Peers:     result.GetPeers(),
		Granted:   result.GetGranted(),
	}

	return nil
}

func (t *Transport) InstallSnapshot(id raft.ServerID, target raft.ServerAddress, args *raft.InstallSnapshotRequest, resp *raft.InstallSnapshotResponse, data io.Reader) error {
	client, err := t.peer(id, target)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	ctx = protocolVersionToMetadata(ctx, args.RPCHeader.ProtocolVersion)
	defer cancel()

	stream, err := client.InstallSnapshot(ctx)
	if err != nil {
		return err
	}

	initial := &raftsvc.InstallSnapshotRequest{
		Term:               args.Term,
		Leader:             args.Leader,
		LastLogIndex:       args.LastLogIndex,
		LastLogTerm:        args.LastLogTerm,
		Peers:              args.Peers,
		Configuration:      args.Configuration,
		ConfigurationIndex: args.ConfigurationIndex,
		Size:               args.Size,
		SnapshotVersion:    int64(args.SnapshotVersion),
	}

	if err = stream.Send(initial); err != nil {
		return err
	}

	buf := make([]byte, 1048576)
	for {
		n, err := data.Read(buf)
		if errors.Is(err, io.EOF) {
			break
		}

		if err == nil && n == 0 {
			break
		}

		if err != nil {
			return err
		}

		request := &raftsvc.InstallSnapshotRequest{
			Data: buf[:n],
		}

		if err = stream.Send(request); err != nil {
			return err
		}
	}

	result, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}

	*resp = raft.InstallSnapshotResponse{
		RPCHeader: args.RPCHeader,
		Term:      result.GetTerm(),
		Success:   result.GetSuccess(),
	}

	return nil
}

func (t *Transport) EncodePeer(_ raft.ServerID, addr raft.ServerAddress) []byte {
	return []byte(addr)
}

func (t *Transport) DecodePeer(bytes []byte) raft.ServerAddress {
	return raft.ServerAddress(bytes)
}

func (t *Transport) SetHeartbeatHandler(cb func(rpc raft.RPC)) {
	t.heartbeatMutex.Lock()
	t.heartbeatFunc = cb
	t.heartbeatMutex.Unlock()
}

func (t *Transport) TimeoutNow(id raft.ServerID, target raft.ServerAddress, args *raft.TimeoutNowRequest, resp *raft.TimeoutNowResponse) error {
	client, err := t.peer(id, target)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	ctx = protocolVersionToMetadata(ctx, args.RPCHeader.ProtocolVersion)
	defer cancel()

	if _, err = client.TimeoutNow(ctx, &raftsvc.TimeoutNowRequest{}); err != nil {
		return err
	}

	*resp = raft.TimeoutNowResponse{
		RPCHeader: args.RPCHeader,
	}

	return nil
}

func (t *Transport) peer(id raft.ServerID, target raft.ServerAddress) (raftsvc.RaftServiceClient, error) {
	t.peerMutex.RLock()
	conn, ok := t.peers[id]
	if ok && !needsDial(conn) {
		defer t.peerMutex.RUnlock()
		return raftsvc.NewRaftServiceClient(conn), nil
	}
	t.peerMutex.RUnlock()

	options := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(
			grpc.UseCompressor("gzip"),
		),
		grpc.WithChainUnaryInterceptor(
			clientinfo.UnaryClientInterceptor(string(t.localID)),
			grpc_retry.UnaryClientInterceptor(),
		),
		grpc.WithChainStreamInterceptor(
			clientinfo.StreamClientInterceptor(string(t.localID)),
			grpc_retry.StreamClientInterceptor(),
		),
	}

	if t.tls != nil {
		options = append(options, grpc.WithTransportCredentials(credentials.NewTLS(t.tls)))
	} else {
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, string(target), options...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial peer %s: %w", target, err)
	}

	t.peerMutex.Lock()
	t.peers[id] = conn
	t.peerMutex.Unlock()

	return raftsvc.NewRaftServiceClient(conn), nil
}

func needsDial(conn *grpc.ClientConn) bool {
	state := conn.GetState()

	return state == connectivity.Shutdown || state == connectivity.TransientFailure
}

func logsToProto(in []*raft.Log) []*raftpb.Log {
	out := make([]*raftpb.Log, len(in))
	for i, log := range in {
		out[i] = &raftpb.Log{
			Index:      log.Index,
			Term:       log.Term,
			Type:       typeToProto(log.Type),
			Data:       log.Data,
			Extensions: log.Extensions,
			AppendedAt: timestamppb.New(log.AppendedAt),
		}
	}

	return out
}

func typeToProto(in raft.LogType) raftpb.LogType {
	switch in {
	case raft.LogCommand:
		return raftpb.LogType_LOG_TYPE_COMMAND
	case raft.LogNoop:
		return raftpb.LogType_LOG_TYPE_NOOP
	case raft.LogAddPeerDeprecated:
		return raftpb.LogType_LOG_TYPE_ADD_PEER_DEPRECATED
	case raft.LogRemovePeerDeprecated:
		return raftpb.LogType_LOG_TYPE_REMOVE_PEER_DEPRECATED
	case raft.LogBarrier:
		return raftpb.LogType_LOG_TYPE_BARRIER
	case raft.LogConfiguration:
		return raftpb.LogType_LOG_TYPE_CONFIGURATION
	default:
		return raftpb.LogType_LOG_TYPE_UNSPECIFIED
	}
}

func protocolVersionToMetadata(ctx context.Context, pv raft.ProtocolVersion) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "X-Protocol-Version", strconv.Itoa(int(pv)))
}
