package raft

import (
	"context"
	"io"
	"strconv"
	"strings"

	"github.com/hashicorp/raft"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	raftsvc "github.com/davidsbond/arrebato/internal/proto/arrebato/raft/service/v1"
	raftpb "github.com/davidsbond/arrebato/internal/proto/arrebato/raft/v1"
)

type (
	GRPC struct {
		transport *Transport
	}
)

func NewGRPC(transport *Transport) *GRPC {
	return &GRPC{transport: transport}
}

func (svr *GRPC) Register(registrar grpc.ServiceRegistrar, healthSvr *health.Server) {
	raftsvc.RegisterRaftServiceServer(registrar, svr)
	healthSvr.SetServingStatus(raftsvc.RaftService_ServiceDesc.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)
}

func (svr *GRPC) AppendEntries(ctx context.Context, request *raftsvc.AppendEntriesRequest) (*raftsvc.AppendEntriesResponse, error) {
	pv, err := protocolVersionFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	entries, err := logsFromProto(request.GetEntries())
	if err != nil {
		return nil, err
	}

	command := &raft.AppendEntriesRequest{
		RPCHeader: raft.RPCHeader{
			ProtocolVersion: pv,
		},
		Term:              request.GetTerm(),
		Leader:            request.GetLeader(),
		PrevLogEntry:      request.GetPrevLogEntry(),
		PrevLogTerm:       request.GetPrevLogTerm(),
		Entries:           entries,
		LeaderCommitIndex: request.GetLeaderCommitIndex(),
	}

	response, err := svr.awaitRPCResponse(ctx, command, nil)
	if err != nil {
		return nil, err
	}

	result, ok := response.(*raft.AppendEntriesResponse)
	if !ok {
		//TODO: error
	}

	return &raftsvc.AppendEntriesResponse{
		Term:           result.Term,
		LastLog:        result.LastLog,
		Success:        result.Success,
		NoRetryBackoff: result.NoRetryBackoff,
	}, nil
}

func (svr *GRPC) RequestVote(ctx context.Context, request *raftsvc.RequestVoteRequest) (*raftsvc.RequestVoteResponse, error) {
	pv, err := protocolVersionFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	command := &raft.RequestVoteRequest{
		RPCHeader: raft.RPCHeader{
			ProtocolVersion: pv,
		},
		Term:               request.GetTerm(),
		Candidate:          request.GetCandidate(),
		LastLogIndex:       request.GetLastLogIndex(),
		LastLogTerm:        request.GetLastLogTerm(),
		LeadershipTransfer: request.GetLeadershipTransfer(),
	}

	response, err := svr.awaitRPCResponse(ctx, command, nil)
	if err != nil {
		return nil, err
	}

	result, ok := response.(*raft.RequestVoteResponse)
	if !ok {
		//TODO: error
	}

	return &raftsvc.RequestVoteResponse{
		Term:    result.Term,
		Peers:   result.Peers,
		Granted: result.Granted,
	}, nil
}

func (svr *GRPC) TimeoutNow(ctx context.Context, _ *raftsvc.TimeoutNowRequest) (*raftsvc.TimeoutNowResponse, error) {
	pv, err := protocolVersionFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	command := &raft.TimeoutNowRequest{
		RPCHeader: raft.RPCHeader{
			ProtocolVersion: pv,
		},
	}

	if _, err = svr.awaitRPCResponse(ctx, command, nil); err != nil {
		return nil, err
	}

	return &raftsvc.TimeoutNowResponse{}, nil
}

func (svr *GRPC) InstallSnapshot(server raftsvc.RaftService_InstallSnapshotServer) error {
	ctx := server.Context()
	pv, err := protocolVersionFromMetadata(ctx)
	if err != nil {
		return err
	}

	initial, err := server.Recv()
	if err != nil {
		return err
	}

	command := raft.InstallSnapshotRequest{
		RPCHeader: raft.RPCHeader{
			ProtocolVersion: pv,
		},
		SnapshotVersion:    raft.SnapshotVersion(initial.GetSnapshotVersion()),
		Term:               initial.GetTerm(),
		Leader:             initial.GetLeader(),
		LastLogIndex:       initial.GetLastLogIndex(),
		LastLogTerm:        initial.GetLastLogTerm(),
		Peers:              initial.GetPeers(),
		Configuration:      initial.GetConfiguration(),
		ConfigurationIndex: initial.GetConfigurationIndex(),
		Size:               initial.GetSize(),
	}

	if _, err = svr.awaitRPCResponse(ctx, command, &snapshotReader{stream: server}); err != nil {
		return err
	}

	return nil
}

func protocolVersionFromMetadata(ctx context.Context) (raft.ProtocolVersion, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, status.Error(codes.InvalidArgument, "metadata does not contain protocol version")
	}

	str := md.Get("X-Protocol-Version")
	v, err := strconv.Atoi(strings.Join(str, ""))
	if err != nil {
		return 0, status.Errorf(codes.InvalidArgument, "invalid value for protocol version: %v", err)
	}

	return raft.ProtocolVersion(v), nil
}

func logsFromProto(logs []*raftpb.Log) ([]*raft.Log, error) {
	out := make([]*raft.Log, len(logs))
	for i, log := range logs {
		l, err := logFromProto(log)
		if err != nil {
			return nil, err
		}

		out[i] = l
	}

	return out, nil
}

func logFromProto(log *raftpb.Log) (*raft.Log, error) {
	logType, err := typeFromProto(log.GetType())
	if err != nil {
		return nil, err
	}

	return &raft.Log{
		Index:      log.GetIndex(),
		Term:       log.GetTerm(),
		Type:       logType,
		Data:       log.GetData(),
		Extensions: log.GetExtensions(),
		AppendedAt: log.GetAppendedAt().AsTime(),
	}, nil
}

func typeFromProto(in raftpb.LogType) (raft.LogType, error) {
	switch in {
	case raftpb.LogType_LOG_TYPE_COMMAND:
		return raft.LogCommand, nil
	case raftpb.LogType_LOG_TYPE_NOOP:
		return raft.LogNoop, nil
	case raftpb.LogType_LOG_TYPE_ADD_PEER_DEPRECATED:
		return raft.LogAddPeerDeprecated, nil
	case raftpb.LogType_LOG_TYPE_REMOVE_PEER_DEPRECATED:
		return raft.LogRemovePeerDeprecated, nil
	case raftpb.LogType_LOG_TYPE_BARRIER:
		return raft.LogBarrier, nil
	case raftpb.LogType_LOG_TYPE_CONFIGURATION:
		return raft.LogConfiguration, nil
	default:
		return 0, status.Errorf(codes.InvalidArgument, "invalid log type: %v", in)
	}
}

func (svr *GRPC) awaitRPCResponse(ctx context.Context, in any, reader io.Reader) (any, error) {
	resp := make(chan raft.RPCResponse, 1)
	defer close(resp)

	request := raft.RPC{
		Command:  in,
		RespChan: resp,
		Reader:   reader,
	}

	if isHeartbeat(in) {
		svr.transport.heartbeatMutex.Lock()
		fn := svr.transport.heartbeatFunc
		svr.transport.heartbeatMutex.Unlock()
		fn(request)
	} else {
		// Write the inbound RPC request to the channel. This will go into the internal raft implementation where it will
		// be acted upon.
		select {
		case <-ctx.Done():
			return nil, status.FromContextError(ctx.Err()).Err()
		case svr.transport.rpcs <- request:
			break
		}
	}

	// We await a response from the internal raft implementation for our particular request.
	select {
	case <-ctx.Done():
		return nil, status.FromContextError(ctx.Err()).Err()
	case result := <-resp:
		if result.Error != nil {
			return nil, status.Error(codes.Internal, result.Error.Error())
		}

		return result.Response, nil
	}
}

type (
	snapshotReader struct {
		stream raftsvc.RaftService_InstallSnapshotServer
		buffer []byte
	}
)

func (sr *snapshotReader) Read(b []byte) (int, error) {
	// If the previous call to read did not take the entire data contents, we return the
	// contents of the buffer.
	if len(sr.buffer) > 0 {
		n := copy(b, sr.buffer)
		sr.buffer = sr.buffer[n:]
		return n, nil
	}

	msg, err := sr.stream.Recv()
	if err != nil {
		return 0, err
	}

	n := copy(b, msg.GetData())

	// If the call to read did not provide a buffer big enough for the entire
	// contents of the snapshot, we store it in the buffer for the next call.
	if n < len(msg.GetData()) {
		sr.buffer = msg.GetData()[n:]
	}

	return n, nil
}

func isHeartbeat(command any) bool {
	req, ok := command.(*raft.AppendEntriesRequest)
	if !ok {
		return false
	}

	return req.Term != 0 &&
		len(req.Leader) != 0 &&
		req.PrevLogEntry == 0 &&
		req.PrevLogTerm == 0 &&
		len(req.Entries) == 0 &&
		req.LeaderCommitIndex == 0
}
