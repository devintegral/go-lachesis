package posnode

import (
	"context"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/Fantom-foundation/go-lachesis/src/posnode/api"
)

type (
	// connection wraps grpc.ClientConn
	connection struct {
		*grpc.ClientConn
		addr    string
		created time.Time
		used    int
	}

	// connPool is connections to peers.
	connPool struct {
		cache map[string]*connection
		size  int

		connectTimeout time.Duration
		opts           []grpc.DialOption

		sync.RWMutex
	}
)

func (n *Node) initClient() {
	if n.connPool.cache != nil {
		return
	}
	n.connPool.size = peersCount * 2
	n.connPool.cache = make(map[string]*connection, n.connPool.size)
	n.connPool.connectTimeout = n.conf.ConnectTimeout

	n.connPool.opts = append(n.connPool.opts,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(api.ClientAuth(n.key)))
}

// ConnectTo connects to other node service.
func (n *Node) ConnectTo(peer *Peer) (client api.NodeClient, free func(), fail func(error), err error) {
	addr := n.NetAddrOf(peer.Host)
	n.log.Debugf("connect to %s", addr)

	c, err := n.connPool.Get(addr)
	if err != nil {
		err = errors.Wrapf(err, "connect to: %s", addr)
		n.log.Warn(err)
		return
	}

	// used should be decremented once
	var released uint32
	free = func() {
		count := atomic.CompareAndSwapUint32(&released, 0, 1)
		n.connPool.Release(c, count, nil)
	}
	fail = func(err error) {
		count := atomic.CompareAndSwapUint32(&released, 0, 1)
		n.connPool.Release(c, count, err)
	}

	client = api.NewNodeClient(c.ClientConn)

	return
}

/*
 * connectionPool utils:
 */

func (cc *connPool) Get(addr string) (*connection, error) {
	cc.Lock()
	defer cc.Unlock()

	conn := cc.cache[addr]
	if conn == nil {
		// make new
		var err error
		conn, err = cc.newConn(addr)
		if err != nil {
			return nil, err
		}
		cc.cache[addr] = conn

		if len(cc.cache) >= cc.size {
			go cc.Clean()
		}
	}

	conn.used += 1

	return conn, nil
}

func (cc *connPool) Release(c *connection, count bool, err error) {
	cc.Lock()
	defer cc.Unlock()

	if count {
		c.used -= 1
	}

	// try to close if error now or before
	if cached := cc.cache[c.addr]; err != nil || c != cached {
		if c == cached {
			delete(cc.cache, c.addr)
		}
		if c.used < 1 {
			c.Close()
		}
	}
}

func (cc *connPool) Clean() {
	cc.Lock()
	defer cc.Unlock()

	if len(cc.cache) < cc.size {
		return
	}

	all := make([]*connection, 0, len(cc.cache))
	for _, c := range cc.cache {
		all = append(all, c)
	}
	sort.Sort(byCreation(all))
	old := all[cc.size/2:]

	for _, c := range old {
		if cached := cc.cache[c.addr]; c == cached {
			delete(cc.cache, c.addr)
			delete(api.ConnsID(), c.addr)
		}
	}
}

func (сс *connPool) newConn(addr string) (*connection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), сс.connectTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr, append(сс.opts, grpc.WithBlock())...)
	if err != nil {
		return nil, err
	}

	return &connection{
		ClientConn: conn,
		addr:       addr,
		created:    time.Now(),
	}, nil
}

/*
 * sorting:
 */

type byCreation []*connection

func (s byCreation) Len() int { return len(s) }

func (s byCreation) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s byCreation) Less(i, j int) bool {
	return s[i].created.Before(s[j].created)
}
