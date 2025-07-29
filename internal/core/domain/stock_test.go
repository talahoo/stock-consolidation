package domain_test

import (
	"encoding/json"
	"testing"
	"time"

	"stock-consolidation/internal/core/domain"
)

func TestStockUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    domain.Stock
		wantErr bool
	}{
		{
			name: "valid json with RFC3339 timestamp",
			json: `{
				"id": "123e4567-e89b-12d3-a456-426614174000",
				"product_id": 1,
				"branch_id": 2,
				"quantity": 100,
				"reserved": 10,
				"created_at": "2025-07-29T05:17:55.443242",
				"updated_at": "2025-07-29T05:17:55.443242"
			}`,
			want: domain.Stock{
				ID:        "123e4567-e89b-12d3-a456-426614174000",
				ProductID: 1,
				BranchID:  2,
				Quantity:  100,
				Reserved:  10,
				CreatedAt: time.Date(2025, 7, 29, 5, 17, 55, 443242000, time.UTC),
				UpdatedAt: time.Date(2025, 7, 29, 5, 17, 55, 443242000, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "valid json with postgres timestamp",
			json: `{
				"id": "123e4567-e89b-12d3-a456-426614174000",
				"product_id": 1,
				"branch_id": 2,
				"quantity": 100,
				"reserved": 10,
				"created_at": "2025-07-29T05:17:55.443242",
				"updated_at": "2025-07-29T05:17:55.443242"
			}`,
			want: domain.Stock{
				ID:        "123e4567-e89b-12d3-a456-426614174000",
				ProductID: 1,
				BranchID:  2,
				Quantity:  100,
				Reserved:  10,
				CreatedAt: time.Date(2025, 7, 29, 5, 17, 55, 443242000, time.UTC),
				UpdatedAt: time.Date(2025, 7, 29, 5, 17, 55, 443242000, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "invalid json",
			json: `{
				"id": "123e4567-e89b-12d3-a456-426614174000",
				"product_id": 1,
				"branch_id": 2,
				"quantity": 100,
				"reserved": 10,
				"created_at": "invalid",
				"updated_at": "invalid"
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got domain.Stock
			err := json.Unmarshal([]byte(tt.json), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("Stock.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.ID != tt.want.ID {
					t.Errorf("Stock.UnmarshalJSON() got = %v, want %v", got.ID, tt.want.ID)
				}
				if got.ProductID != tt.want.ProductID {
					t.Errorf("Stock.UnmarshalJSON() got = %v, want %v", got.ProductID, tt.want.ProductID)
				}
				if got.BranchID != tt.want.BranchID {
					t.Errorf("Stock.UnmarshalJSON() got = %v, want %v", got.BranchID, tt.want.BranchID)
				}
				if got.Quantity != tt.want.Quantity {
					t.Errorf("Stock.UnmarshalJSON() got = %v, want %v", got.Quantity, tt.want.Quantity)
				}
				if got.Reserved != tt.want.Reserved {
					t.Errorf("Stock.UnmarshalJSON() got = %v, want %v", got.Reserved, tt.want.Reserved)
				}
				// Compare timestamps with some tolerance for timezone differences
				if got.CreatedAt.Unix() != tt.want.CreatedAt.Unix() {
					t.Errorf("Stock.UnmarshalJSON() got = %v, want %v", got.CreatedAt, tt.want.CreatedAt)
				}
				if got.UpdatedAt.Unix() != tt.want.UpdatedAt.Unix() {
					t.Errorf("Stock.UnmarshalJSON() got = %v, want %v", got.UpdatedAt, tt.want.UpdatedAt)
				}
			}
		})
	}
}
