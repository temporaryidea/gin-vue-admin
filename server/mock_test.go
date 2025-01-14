package __mock

import (
	"code.byted.org/overpass/bytedance_bits_meta/kitex_gen/bytedance/bits/meta"
	"code.byted.org/overpass/bytedance_bits_meta/rpc/bytedance_bits_meta"
	"context"
	"github.com/bytedance/mockey"
	"testing"
	"training/examples/some_depends/domain"
	"training/examples/some_depends/infra/entity"
)

func TestUnitTestServiceDomainRule_BitsToOneSite(t *testing.T) {
	tests := []struct {
		name           string
		bitsID        int64
		mockDB        func() *entity.BitsInfo
		mockRPC       func() (*meta.QueryAppSimpleInfoByIdsResponse, error)
		expectID      int64
		expectErr     bool
		expectDBWrite bool
	}{
		{
			name:    "cache hit",
			bitsID:  123,
			mockDB:  func() *entity.BitsInfo { return &entity.BitsInfo{BitsID: 123, OneSiteID: 456} },
			expectID: 456,
		},
		{
			name:    "cache miss - rpc success",
			bitsID:  123,
			mockDB:  func() *entity.BitsInfo { return nil },
			mockRPC: func() (*meta.QueryAppSimpleInfoByIdsResponse, error) {
				spaceID := int64(456)
				return &meta.QueryAppSimpleInfoByIdsResponse{
					AppInfos: []*meta.AppSimpleInfo{{SpaceId: &spaceID}},
				}, nil
			},
			expectID:      456,
			expectDBWrite: true,
		},
		{
			name:    "cache miss - rpc error",
			bitsID:  123,
			mockDB:  func() *entity.BitsInfo { return nil },
			mockRPC: func() (*meta.QueryAppSimpleInfoByIdsResponse, error) {
				return nil, fmt.Errorf("rpc error")
			},
			expectErr: true,
		},
		{
			name:    "cache miss - app not found",
			bitsID:  123,
			mockDB:  func() *entity.BitsInfo { return nil },
			mockRPC: func() (*meta.QueryAppSimpleInfoByIdsResponse, error) {
				return &meta.QueryAppSimpleInfoByIdsResponse{
					AppInfos: []*meta.AppSimpleInfo{},
				}, nil
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock DB operations
			mockRepo := &mockBitsInfoReader{}
			mockey.Mock(mockRepo.FindBitsInfoFromDB).Return(tt.mockDB()).Build()
			
			if tt.expectDBWrite {
				mockey.Mock(mockRepo.UpdateBitsInfoToDB).To(func(info *entity.BitsInfo) {
					if info.BitsID != tt.bitsID || info.OneSiteID != tt.expectID {
						t.Errorf("unexpected DB write with BitsID=%d, OneSiteID=%d", info.BitsID, info.OneSiteID)
					}
				}).Build()
			}

			// Mock RPC call if needed
			if tt.mockRPC != nil {
				mockey.Mock(bytedance_bits_meta.RawCall.QueryAppSimpleInfoByIds).Return(tt.mockRPC()).Build()
			}

			// Create service instance
			svc := &UnitTestServiceDomainRule{repo: mockRepo}

			// Execute test
			gotID, err := svc.BitsToOneSite(context.Background(), tt.bitsID)

			// Verify results
			if tt.expectErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if gotID != tt.expectID {
					t.Errorf("expected ID %d but got %d", tt.expectID, gotID)
				}
			}
		})
	}
}

// Mock implementation of BitsInfoReader interface
type mockBitsInfoReader struct {
	domain.BitsInfoReader
}

func (m *mockBitsInfoReader) FindBitsInfoFromDB(bitsID int64) *entity.BitsInfo {
	return nil // Will be mocked by mockey
}

func (m *mockBitsInfoReader) UpdateBitsInfoToDB(info *entity.BitsInfo) {
	// Will be mocked by mockey
}
