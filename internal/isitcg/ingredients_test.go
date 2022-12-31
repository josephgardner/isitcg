package isitcg

import (
	"reflect"
	"testing"
)

func TestProductFromHash(t *testing.T) {
	for _, tt := range []struct {
		name  string
		rules []Rule
		hash  string
		want  Product
	}{
		{
			name: "search console",
			hash: "bZA9TwMxDIb_inXzwcDAwHY9UT4EaoGBpYubc3sWvrhyEuiB-O_42oEPsUSxk_f14_ejitVF1RaTBDMptCaz0W-cEjzSDg0zvxLcEfp5E6HV2HFmjWRVXbFrl8V4w9TBM2ayGmbUU8zGg0YuA9xT7jUV2fhrDS1lQhsFGgnaqxw6v0vLR7uFbTFy-EE1t8IZLvfZMORJ-fdrq0FjybBgN76SMZD59JOz8xqe1NacjwP_yhpRghnaGjuKidO07QZuCwcnvh47UxnffcGHwlHRvVheoBm8gCZwl_5F6XHQgYVgLvpG9o09N9waxsl76UnpfnRGjAc03Hk2HsfErnK6slWsPr8A",
			want: Product{
				"Curls Blueberry Bliss Reparative Leave In Conditioner",
				"Purified Water, Behentrimonium Methosulfate, Cetearyl Alcohol, Cetyl Alcohol, Certified Organic Blueberry Fruit Extract, Certified Organic Coconut Oil, Glycereth-26, Sorbitol, Certified Organic Aloe Barbadensis Leaf Juice, Hydrolyzed Quinoa, Silk Amino Acids, Certified Organic Chamomile Flower Extract, Fragrance, Phenoxyethanol, Caprylyl Glycol.\r\n",
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			h := defaultIngredientHandler{}
			if got := h.ProductFromHash(tt.hash); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("defaultIngredientHandler.ProductFromHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
