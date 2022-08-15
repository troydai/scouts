package entry

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-zookeeper/zk"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Options(
	fx.Provide(ProvideZookeeperClient),
)

type ZkClient struct {
	conn *zk.Conn
}

func ProvideZookeeperClient(lc fx.Lifecycle, logger *zap.Logger) *ZkClient {
	c := &ZkClient{}

	var eventCh <-chan zk.Event
	eventLoopCtx, cancelEventLoop := context.WithCancel(context.Background())

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				var err error
				c.conn, eventCh, err = zk.Connect([]string{"zk-client"}, 10*time.Second)
				if err != nil {
					return fmt.Errorf("fail to create zk connection: %w", err)
				}

				go func() {
					for {
						select {
						case <-eventLoopCtx.Done():
						case event := <-eventCh:
							logger.Debug("recieved zk event", zap.Any("event", event))
						}
					}
				}()

				acl := []zk.ACL{
					{
						Perms:  zk.PermAll,
						Scheme: "world",
						ID:     "anyone",
					},
				}
				_, err = c.conn.CreateContainer("/scouts", nil, zk.FlagTTL, acl)
				if err != nil {
					if !errors.Is(err, zk.ErrNodeExists) {
						return fmt.Errorf("fail to create scouts container: %w", err)
					}
				}

				_, err = c.conn.Create(fmt.Sprintf("/scouts/%s", os.Getenv("HOSTNAME")), nil, zk.FlagEphemeral, acl)
				if err != nil {
					return fmt.Errorf("fail to create scout node: %w", err)
				}

				return nil
			},
			OnStop: func(ctx context.Context) error {
				cancelEventLoop()

				if c.conn != nil {
					c.conn.Close()
				}
				return nil
			},
		},
	)

	return c
}
