package aperture_test

import (
	"sync"
	"testing"
	"time"

	"github.com/hnlq715/go-loadbalance/aperture"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/balancer"
)

func TestAperture(t *testing.T) {
	t.Run("0 item", func(t *testing.T) {
		ll := aperture.New()
		item, done := ll.Next()
		done(balancer.DoneInfo{})
		assert.Nil(t, item)
	})

	t.Run("1 item", func(t *testing.T) {
		ll := aperture.New()
		ll.SetLocalPeers([]string{"1"})
		ll.SetRemotePeers([]interface{}{"8"})
		ll.SetLocalPeerID("1")

		item, done := ll.Next()
		done(balancer.DoneInfo{})
		assert.Equal(t, "8", item)
	})

	t.Run("3 items", func(t *testing.T) {
		ll := aperture.New()
		ll.SetLocalPeers([]string{"1", "2", "3"})
		ll.SetRemotePeers([]interface{}{"8", "9", "10"})
		ll.SetLocalPeerID("1")

		item, done := ll.Next()
		done(balancer.DoneInfo{})
		assert.Equal(t, "8", item)

		ll.SetLocalPeerID("2")

		item, done = ll.Next()
		done(balancer.DoneInfo{})
		assert.Equal(t, "9", item)

		ll.SetLocalPeerID("3")

		item, done = ll.Next()
		done(balancer.DoneInfo{})
		assert.Equal(t, "10", item)
	})

	t.Run("count", func(t *testing.T) {
		ll := aperture.New()
		ll.SetLocalPeers([]string{"1", "2", "3"})
		ll.SetRemotePeers([]interface{}{"8", "9", "10", "11", "12"})
		ll.SetLocalPeerID("1")
		ll.SetLogicalWidth(2)

		countMap := make(map[interface{}]int)

		totalCount := 10000
		wg := sync.WaitGroup{}
		wg.Add(totalCount)

		mu := sync.Mutex{}
		for i := 0; i < totalCount; i++ {
			go func() {
				defer wg.Done()
				item, done := ll.Next()
				time.Sleep(1 * time.Second)
				done(balancer.DoneInfo{})

				mu.Lock()
				countMap[item]++
				mu.Unlock()
			}()
		}

		wg.Wait()

		total := 0
		for _, count := range countMap {
			total += count
		}
		assert.Less(t, 2990, countMap["8"])
		assert.Less(t, 2990, countMap["9"])
		assert.Less(t, 2990, countMap["10"])
		assert.Less(t, 990, countMap["11"])

		assert.Equal(t, totalCount, total)
	})
}