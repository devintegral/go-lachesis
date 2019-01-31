package peer_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/rpc"
	"reflect"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/Fantom-foundation/go-lachesis/src/common"
	lnet "github.com/Fantom-foundation/go-lachesis/src/net"
	"github.com/Fantom-foundation/go-lachesis/src/net/fakenet"
	"github.com/Fantom-foundation/go-lachesis/src/net/peer"
	"github.com/Fantom-foundation/go-lachesis/src/poset"
)

var (
	expEagerSyncRequest = &lnet.EagerSyncRequest{
		FromID: fakeAddress(0),
		Events: []poset.WireEvent{
			{
				Body: poset.WireBody{
					Transactions:         [][]byte(nil),
					SelfParentIndex:      1,
					OtherParentCreatorID: fakeAddress(10).Bytes(),
					OtherParentIndex:     0,
					CreatorID:            fakeAddress(9).Bytes(),
				},
			},
		},
	}
	expEagerSyncResponse  = &lnet.EagerSyncResponse{FromID: fakeAddress(1), Success: true}
	expFastForwardRequest = &lnet.FastForwardRequest{FromID: fakeAddress(0)}
	expSyncRequest        = &lnet.SyncRequest{
		FromID: fakeAddress(0),
		Known: map[common.Address]int64{
			fakeAddress(0): 1,
			fakeAddress(1): 2,
			fakeAddress(2): 3,
		},
	}

	expSyncResponse = &lnet.SyncResponse{
		FromID: fakeAddress(1),
		Events: []poset.WireEvent{
			{
				Body: poset.WireBody{
					Transactions:         [][]byte(nil),
					SelfParentIndex:      1,
					OtherParentCreatorID: fakeAddress(10).Bytes(),
					OtherParentIndex:     0,
					CreatorID:            fakeAddress(9).Bytes(),
				},
			},
		},
		Known: map[common.Address]int64{
			fakeAddress(0): 5,
			fakeAddress(1): 5,
			fakeAddress(2): 6,
		},
	}
	testError = errors.New("error")
)

type mockRpcClient struct {
	t    *testing.T
	err  error
	resp interface{}
}

func newRPCClient(t *testing.T, err error, resp interface{}) *mockRpcClient {
	return &mockRpcClient{t: t, err: err, resp: resp}
}

func (m *mockRpcClient) Go(serviceMethod string, args interface{},
	reply interface{}, done chan *rpc.Call) *rpc.Call {
	raw, err := json.Marshal(m.resp)
	if err != nil {
		m.t.Fatal(err)
	}
	if err := json.Unmarshal(raw, reply); err != nil {
		m.t.Fatal(err)
	}
	call := &rpc.Call{}
	call.Done = make(chan *rpc.Call, 10)
	call.Error = m.err
	call.Reply = reply
	call.Done <- call
	return call
}

func (m *mockRpcClient) Close() error {
	return m.err
}

func newClient(t *testing.T, m *mockRpcClient) *peer.Client {
	cli, err := peer.NewClient(m)
	if err != nil {
		t.Fatal(err)
	}
	return cli
}

func TestClientSync(t *testing.T) {
	ctx := context.Background()
	m := newRPCClient(t, testError, expSyncResponse)
	cli := newClient(t, m)
	defer cli.Close()

	resp := &lnet.SyncResponse{}
	if err := cli.Sync(
		ctx, expSyncRequest, resp); err != testError {
		t.Fatalf("expected error: %s, got: %s", testError, err)
	}

	m.err = nil

	if err := cli.Sync(
		ctx, expSyncRequest, resp); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(resp, expSyncResponse) {
		t.Fatalf("failed to get response, expected: %+v, got: %+v",
			expSyncResponse, resp)
	}
}

func TestClientForceSync(t *testing.T) {
	ctx := context.Background()
	m := newRPCClient(t, testError, expEagerSyncResponse)
	cli := newClient(t, m)
	defer cli.Close()

	resp := &lnet.EagerSyncResponse{}
	if err := cli.ForceSync(
		ctx, expEagerSyncRequest, resp); err != testError {
		t.Fatalf("expected error: %s, got: %s", testError, err)
	}

	m.err = nil

	if err := cli.ForceSync(
		ctx, expEagerSyncRequest, resp); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(resp, expEagerSyncResponse) {
		t.Fatalf("failed to get response, expected: %+v, got: %+v",
			expEagerSyncResponse, resp)
	}
}

func TestClientFastForward(t *testing.T) {
	expResponse := newFastForwardResponse(t)
	ctx := context.Background()
	m := newRPCClient(t, testError, expResponse)
	cli := newClient(t, m)
	defer cli.Close()

	resp := &lnet.FastForwardResponse{}
	if err := cli.FastForward(
		ctx, expFastForwardRequest, resp); err != testError {
		t.Fatalf("expected error: %s, got: %s", testError, err)
	}

	m.err = nil

	if err := cli.FastForward(
		ctx, expFastForwardRequest, resp); err != nil {
		t.Fatal(err)
	}

	checkFastForwardResponse(t, expResponse, resp)
}

func TestNewClient(t *testing.T) {
	timeout := time.Second
	conf := &peer.BackendConfig{
		ReceiveTimeout: timeout,
		ProcessTimeout: timeout,
		IdleTimeout:    timeout,
	}
	done := make(chan struct{})
	defer close(done)

	address := newAddress()
	backend := newBackend(t, conf, logger, address, done,
		expSyncResponse, 0, net.Listen)
	defer backend.Close()

	rpcCli, err := peer.NewRPCClient(
		peer.TCP, address, time.Second, net.DialTimeout)
	if err != nil {
		t.Fatal(err)
	}

	cli, err := peer.NewClient(rpcCli)
	if err != nil {
		t.Fatal(err)
	}

	resp := &lnet.SyncResponse{}
	if err := cli.Sync(
		context.Background(), &lnet.SyncRequest{}, resp); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(resp, expSyncResponse) {
		t.Fatalf("failed to get response, expected: %+v, got: %+v",
			expSyncResponse, resp)
	}
}

func TestFakeNet(t *testing.T) {
	timeout := time.Second
	conf := &peer.BackendConfig{
		ReceiveTimeout: timeout,
		ProcessTimeout: timeout,
		IdleTimeout:    timeout,
	}
	done := make(chan struct{})
	defer close(done)

	// Create fake network
	network := fakenet.NewNetwork()

	address := newAddress()
	backend := newBackend(t, conf, logger, address, done,
		expSyncResponse, 0, network.CreateListener)
	defer backend.Close()

	rpcCli, err := peer.NewRPCClient(peer.TCP, address, time.Second,
		network.CreateNetConn)
	if err != nil {
		t.Fatal(err)
	}

	cli, err := peer.NewClient(rpcCli)
	if err != nil {
		t.Fatal(err)
	}

	resp := &lnet.SyncResponse{}
	if err := cli.Sync(
		context.Background(), &lnet.SyncRequest{}, resp); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(resp, expSyncResponse) {
		t.Fatalf("failed to get response, expected: %+v, got: %+v",
			expSyncResponse, resp)
	}
}

func newFastForwardResponse(t *testing.T) *lnet.FastForwardResponse {
	frame := poset.Frame{}
	block, err := poset.NewBlockFromFrame(1, frame)
	if err != nil {
		t.Fatal(err)
	}

	return &lnet.FastForwardResponse{
		FromID:   fakeAddress(1),
		Block:    block,
		Frame:    frame,
		Snapshot: []byte("snapshot"),
	}
}

func checkFastForwardResponse(t *testing.T, exp, got *lnet.FastForwardResponse) {
	if !got.Block.Equals(&exp.Block) || !got.Frame.Equals(&exp.Frame) ||
		got.FromID != exp.FromID || !bytes.Equal(got.Snapshot, exp.Snapshot) {
		t.Fatalf("bad response, expected: %+v, got: %+v", exp, got)
	}
}
