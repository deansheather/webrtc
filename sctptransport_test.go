// +build !js

package webrtc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateDataChannelID(t *testing.T) {
	sctpTransportWithChannels := func(ids []uint16) *SCTPTransport {
		ret := &SCTPTransport{dataChannels: []*DataChannel{}}

		for _, id := range ids {
			id := id
			dc := &DataChannel{
				id: &id,
			}
			dc.readyState.Store(DataChannelStateOpen)

			ret.dataChannels = append(ret.dataChannels, dc)
		}

		return ret
	}

	t.Run("OK", func(t *testing.T) {
		testCases := []struct {
			role   DTLSRole
			s      *SCTPTransport
			result uint16
		}{
			{DTLSRoleClient, sctpTransportWithChannels([]uint16{}), 0},
			{DTLSRoleClient, sctpTransportWithChannels([]uint16{1}), 0},
			{DTLSRoleClient, sctpTransportWithChannels([]uint16{0}), 2},
			{DTLSRoleClient, sctpTransportWithChannels([]uint16{0, 2}), 4},
			{DTLSRoleClient, sctpTransportWithChannels([]uint16{0, 4}), 2},
			{DTLSRoleServer, sctpTransportWithChannels([]uint16{}), 1},
			{DTLSRoleServer, sctpTransportWithChannels([]uint16{0}), 1},
			{DTLSRoleServer, sctpTransportWithChannels([]uint16{1}), 3},
			{DTLSRoleServer, sctpTransportWithChannels([]uint16{1, 3}), 5},
			{DTLSRoleServer, sctpTransportWithChannels([]uint16{1, 5}), 3},
		}
		for _, testCase := range testCases {
			idPtr := new(uint16)
			err := testCase.s.generateAndSetDataChannelID(testCase.role, &idPtr)
			if err != nil {
				t.Errorf("failed to generate id: %v", err)
				return
			}
			if *idPtr != testCase.result {
				t.Errorf("Wrong id: %d expected %d", *idPtr, testCase.result)
			}
		}
	})

	t.Run("IgnoresClosed", func(t *testing.T) {
		s := sctpTransportWithChannels([]uint16{0})
		for _, dc := range s.dataChannels {
			dc.readyState.Store(DataChannelStateClosed)
		}

		idPtr := new(uint16)
		err := s.generateAndSetDataChannelID(DTLSRoleClient, &idPtr)
		require.NoError(t, err)
		require.NotNil(t, idPtr)
		assert.EqualValues(t, 0, *idPtr)
	})
}
