package azure

import (
	"context"
	"fmt"
	servicebus "github.com/Azure/azure-service-bus-go"
	"log"
)

// create a client

type SBClient struct {
	Queue *servicebus.Queue
}

func NewClient() (*SBClient, error) {

	connStr := "Endpoint=sb://twitchmessages.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=E/4f207CQvnAIQMoeHrfoqXLYg/QbcedYGOHVOKzSgw="
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		log.Println("new namespace error")
		return nil, err
	}

	qm := ns.NewQueueManager()
	target, err := ensureQueue(context.Background(), qm, "test-1")
	if err != nil {
		log.Println("ensure Queue error")
		return nil, err
	}

	queue, err := ns.NewQueue(target.Name)
	if err != nil {
		fmt.Println("new Queue error")
		return nil, err
	}

	return &SBClient{Queue: queue}, nil
}

func (s *SBClient) Close() error {
	return s.Queue.Close(context.Background())
}

func ensureQueue(ctx context.Context, qm *servicebus.QueueManager, name string, opts ...servicebus.QueueManagementOption) (*servicebus.QueueEntity, error) {
	_, err := qm.Get(ctx, name)
	if err == nil {
		_ = qm.Delete(ctx, name)
	}

	qe, err := qm.Put(ctx, name, opts...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return qe, nil
}
