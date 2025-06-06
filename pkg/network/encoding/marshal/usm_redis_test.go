// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build linux && linux_bpf

package marshal

import (
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	model "github.com/DataDog/agent-payload/v5/process"

	"github.com/DataDog/datadog-agent/pkg/network"
	"github.com/DataDog/datadog-agent/pkg/network/protocols/redis"
	"github.com/DataDog/datadog-agent/pkg/process/util"
)

const (
	redisClientPort = uint16(2345)
	redisServerPort = uint16(5432)
)

var (
	redisDefaultConnection = network.ConnectionStats{ConnectionTuple: network.ConnectionTuple{
		Source: localhost,
		Dest:   localhost,
		SPort:  redisClientPort,
		DPort:  redisServerPort,
	}}
)

type RedisSuite struct {
	suite.Suite
}

func TestRedisStats(t *testing.T) {
	suite.Run(t, &RedisSuite{})
}

func (s *RedisSuite) TestFormatRedisStats() {
	t := s.T()

	dummyKey := redis.NewKey(
		localhost,
		localhost,
		redisClientPort,
		redisServerPort,
		redis.GetCommand,
		"dummyKey",
		false,
	)

	in := &network.Connections{
		BufferedData: network.BufferedData{
			Conns: []network.ConnectionStats{
				redisDefaultConnection,
			},
		},
		USMData: network.USMProtocolsData{
			Redis: map[redis.Key]*redis.RequestStats{
				dummyKey: {ErrorToStats: map[bool]*redis.RequestStat{
					false: {
						FirstLatencySample: 1,
						Count:              2,
					},
					true: {
						FirstLatencySample: 1,
						Count:              2,
					},
				}},
			},
		},
	}

	out := &model.DatabaseAggregations{
		Aggregations: []*model.DatabaseStats{
			{
				DbStats: &model.DatabaseStats_Redis{
					Redis: &model.RedisStats{
						Command:   model.RedisCommand_RedisGetCommand,
						KeyName:   dummyKey.KeyName.Get(),
						Truncated: dummyKey.Truncated,
						ErrorToStats: map[int32]*model.RedisStatsEntry{
							0: {
								FirstLatencySample: in.USMData.Redis[dummyKey].ErrorToStats[false].FirstLatencySample,
								Count:              uint32(in.USMData.Redis[dummyKey].ErrorToStats[false].Count),
							},
							1: {
								FirstLatencySample: in.USMData.Redis[dummyKey].ErrorToStats[false].FirstLatencySample,
								Count:              uint32(in.USMData.Redis[dummyKey].ErrorToStats[false].Count),
							},
						},
					},
				},
			},
		},
	}

	encoder := newRedisEncoder(in.USMData.Redis)
	t.Cleanup(encoder.Close)

	aggregations := getRedisAggregations(t, encoder, in.Conns[0])

	require.NotNil(t, aggregations)
	assert.ElementsMatch(t, out.Aggregations, aggregations.Aggregations)
}

func (s *RedisSuite) TestRedisIDCollisionRegression() {
	t := s.T()
	assert := assert.New(t)
	connections := []network.ConnectionStats{
		{ConnectionTuple: network.ConnectionTuple{
			Source: localhost,
			SPort:  redisClientPort,
			Dest:   localhost,
			DPort:  redisServerPort,
			Pid:    1,
		}},
		{ConnectionTuple: network.ConnectionTuple{
			Source: localhost,
			SPort:  redisClientPort,
			Dest:   localhost,
			DPort:  redisServerPort,
			Pid:    2,
		}},
	}

	redisKey := redis.NewKey(
		localhost,
		localhost,
		redisClientPort,
		redisServerPort,
		redis.GetCommand,
		"dummyKey",
		false,
	)

	in := &network.Connections{
		BufferedData: network.BufferedData{
			Conns: connections,
		},
		USMData: network.USMProtocolsData{
			Redis: map[redis.Key]*redis.RequestStats{
				redisKey: {
					ErrorToStats: map[bool]*redis.RequestStat{false: {Count: 10}},
				},
			},
		},
	}

	encoder := newRedisEncoder(in.USMData.Redis)
	t.Cleanup(encoder.Close)
	aggregations := getRedisAggregations(t, encoder, in.Conns[0])
	assert.NotNil(aggregations)
	assert.Equal(10, int(aggregations.Aggregations[0].GetRedis().ErrorToStats[0].Count))

	// assert that the other connections sharing the same (source,destination)
	// addresses but different PIDs *won't* be associated with the Redis stats
	// object
	streamer := NewProtoTestStreamer[*model.Connection]()
	encoder.EncodeConnection(in.Conns[1], model.NewConnectionBuilder(streamer))
	var conn model.Connection
	streamer.Unwrap(t, &conn)
	assert.Empty(conn.DataStreamsAggregations)
}

func (s *RedisSuite) TestRedisLocalhostScenario() {
	t := s.T()
	assert := assert.New(t)
	connections := []network.ConnectionStats{
		{ConnectionTuple: network.ConnectionTuple{
			Source: localhost,
			SPort:  redisClientPort,
			Dest:   localhost,
			DPort:  redisServerPort,
			Pid:    1,
		}},
		{ConnectionTuple: network.ConnectionTuple{
			Source: localhost,
			SPort:  redisServerPort,
			Dest:   localhost,
			DPort:  redisClientPort,
			Pid:    2,
		}},
	}

	redisKey := redis.NewKey(
		localhost,
		localhost,
		redisClientPort,
		redisServerPort,
		redis.GetCommand,
		"dummyKey",
		false,
	)

	in := &network.Connections{
		BufferedData: network.BufferedData{
			Conns: connections,
		},
		USMData: network.USMProtocolsData{
			Redis: map[redis.Key]*redis.RequestStats{
				redisKey: {
					ErrorToStats: map[bool]*redis.RequestStat{false: {Count: 10}},
				},
			},
		},
	}

	encoder := newRedisEncoder(in.USMData.Redis)
	t.Cleanup(encoder.Close)

	// assert that both ends (client:server, server:client) of the connection
	// will have Redis stats
	for _, conn := range in.Conns {
		aggregations := getRedisAggregations(t, encoder, conn)
		assert.NotNil(aggregations.Aggregations[0].GetRedis())
		assert.Equal(10, int(aggregations.Aggregations[0].GetRedis().ErrorToStats[0].Count))
	}
}

func getRedisAggregations(t *testing.T, encoder *redisEncoder, c network.ConnectionStats) *model.DatabaseAggregations {
	streamer := NewProtoTestStreamer[*model.Connection]()
	encoder.EncodeConnection(c, model.NewConnectionBuilder(streamer))

	var conn model.Connection
	streamer.Unwrap(t, &conn)

	var aggregations model.DatabaseAggregations
	err := proto.Unmarshal(conn.DatabaseAggregations, &aggregations)
	require.NoError(t, err)

	return &aggregations
}

func generateBenchMarkPayloadRedis(sourcePortsMax, destPortsMax uint16) network.Connections {
	localhost := util.AddressFromString("127.0.0.1")

	payload := network.Connections{
		BufferedData: network.BufferedData{
			Conns: make([]network.ConnectionStats, sourcePortsMax*destPortsMax),
		},
		USMData: network.USMProtocolsData{
			Redis: make(map[redis.Key]*redis.RequestStats),
		},
	}

	for sport := uint16(0); sport < sourcePortsMax; sport++ {
		for dport := uint16(0); dport < destPortsMax; dport++ {
			index := sport*sourcePortsMax + dport

			payload.Conns[index].Dest = localhost
			payload.Conns[index].Source = localhost
			payload.Conns[index].DPort = dport + 1
			payload.Conns[index].SPort = sport + 1
			if index%2 == 0 {
				payload.Conns[index].IPTranslation = &network.IPTranslation{
					ReplSrcIP:   localhost,
					ReplDstIP:   localhost,
					ReplSrcPort: dport + 1,
					ReplDstPort: sport + 1,
				}
			}

			payload.USMData.Redis[redis.NewKey(
				localhost,
				localhost,
				sport+1,
				dport+1,
				redis.GetCommand,
				"dummyKey",
				false,
			)] = &redis.RequestStats{ErrorToStats: map[bool]*redis.RequestStat{false: {Count: 10}}}
		}
	}

	return payload
}

func commonBenchmarkRedisEncoder(b *testing.B, numberOfPorts uint16) {
	payload := generateBenchMarkPayloadRedis(numberOfPorts, numberOfPorts)
	streamer := NewProtoTestStreamer[*model.Connection]()
	a := model.NewConnectionBuilder(streamer)
	b.ResetTimer()
	b.ReportAllocs()
	var h *redisEncoder
	for i := 0; i < b.N; i++ {
		h = newRedisEncoder(payload.USMData.Redis)
		streamer.Reset()
		for _, conn := range payload.Conns {
			h.EncodeConnection(conn, a)
		}
		h.Close()
	}
}

func BenchmarkRedisEncoder100Requests(b *testing.B) {
	commonBenchmarkRedisEncoder(b, 10)
}

func BenchmarkRedisEncoder10000Requests(b *testing.B) {
	commonBenchmarkRedisEncoder(b, 100)
}
