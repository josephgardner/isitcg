package isitcg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompress(t *testing.T) {
	for _, c := range []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple",
			input: "hello world",
			want:  "ykjNyclXKM8vykkBBAAA__8",
		},
		{
			name:  "complex",
			input: `{"n":"Test Ingredients","i":["Aqua (Purified Water)","Sodium Cocoamphoacetate","Lauryl Glucoside","Sodium Cocoyl Glutamate","Sodium Lauryl Glucose Carboxylate","Decyl Glucoside","Cocamidopropyl Betaine","Glycerin","Panthenol","Potassium Sorbate","Citric Acid","Polysorbate 20","Phenoxyethanol","Menthol","Mentha Piperita (Peppermint) Oil","Aloe Barbadensis (Aloe Vera) Leaf Juice","Carthamus Tinctorius (Safflower) Oil","Achillea Millefolium Extract","Chamomilla Recutita (Matricaria) Flower Extract","Equisetum Arvense Extract","Eucalyptus Globulus (Eucalyptus) Oil","Lavandula Angustifolia (Lavender) Leaf Extract","Melaleuca Alternifolia (Tea Tree) Leaf Oil","Rosmarinus Officinalis (Rosemary) Leaf Extract","Salvia Officinalis (Sage) Leaf Extract","Thymus Vulgaris (Thyme) Extract","Tussilago Farfara (Coltsfoot) Flower Extract","Urtica Dioica (Nettle) Extract"]}`,
			want:  "bJBPb9tIDMW_CjEnG_BhscfcHOcPdmE3Qeymh6IHekRZBKihwuG4EYp-92KspHGanjTge--nR_4IKVyEHWWH_9LBqGFKnsMicLj4GpZPBWF2X4xbpga-oJPNwyJsteHSw0qjYj90ipEcncIirLHYKHArJWrmht6bJ8Wxn8wvyrsMwQptr8-jTJ4rin_wVhqx50YH02EUuCRHTlW4lTGScQqLcI_JO0oq9a2OOdcfbdX2E3XFbhxhGbk5OWTMkwb__lMHNfs8knc4MTaUvHt7IdzzQMZer0PDQNZz8jnccbUsRQku0fbYUMqcYXaaPJLhHNaELfxfOJ5qoHmHfcmw4xRdjUuG2RbbVvQ72W9g7FiEEDb106rUZa6f3TB6pXTYa88iCA8Ui59qbbBuiMY4h5sT7Sxx_VQ4k5celnaklOlcKxFlHLxkuBXdF6mV3oavndZ4xNQUQVimQ8nOtRbCbI1HSk3tftr0jbshQaESEZbiZOk1sCOEnRG9BCb6g-YejVPJcNe2HDmh1EM-aKYebfxA36IcGd-bt3igD8ZdN9Z7PxY5oFVXHdD83FFyZsGDwg1ai4YwW6l4blX9L7f8bM4R4Yq1fmafyF3OeN9-_goAAP__",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, Compress(c.input))
		})
	}
}

func TestDecompress(t *testing.T) {
	for _, c := range []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple_legacy",
			input: "y0jNyclXKM8vykkBAA",
			want:  "hello world",
		},
		{
			name:  "simple",
			input: "ykjNyclXKM8vykkBBAAA__8",
			want:  "hello world",
		},
		{
			name:  "complex_legacy",
			input: "ZVLLbttADPyVhU420EPRY26K4xgt7MaI3fRQ9ECvKIkAtVT24Voo-u_lSnGjNCcRHHJmONrfhStuiiOGaD67xmNF6GIoPhRU3PwoyucEZrFPnmrCynyHiH6p4EEqSp1ZiRXo-lbAYlRMkS0kP7DZcLISqMqt2fCEROim4RfkzQ6aFfiTXAaeZu7Q_senqtBRJb2XXqFblSaXhTY8WPTktNyDiy064VxLhBCy34P408S6oujJmtJSNU7wECbMfPqYG3n3MmBsYeLYaSrtyDZWYPbUq1TM6WCvZUcuLs0DZcGSBc2tXgEVukDBLMbOE3pYmi1Cbb4kstmxnqoSXQrmSM5G8aTl4gB1zfJLo74S2paYEcwuf2rhfMz6Ej3YmFmUQjqFwDyiTXG0tYN8IXhSzfuRbbaxfk4UMCpL6c_qEedYssBDH9XJhuWUOFta_2tePW3hDK5Kqlm6JoVI2ZbGoX10VfY-XvrqcocMjMpjStZn5K4LR73r6BFfFqYIHyV06t2p9kNdkyUHnIPUPiowvGM_AJ9V_s3wAZor66uNYzvkvJ8SNyqglLmhY7OJpK-FoRFzD74Gr0ethGOoRfQXv8vym4-as7kjyZ_FV4yRZ3w___wF",
			want:  `{"n":"Test Ingredients","i":["Aqua (Purified Water)","Sodium Cocoamphoacetate","Lauryl Glucoside","Sodium Cocoyl Glutamate","Sodium Lauryl Glucose Carboxylate","Decyl Glucoside","Cocamidopropyl Betaine","Glycerin","Panthenol","Potassium Sorbate","Citric Acid","Polysorbate 20","Phenoxyethanol","Menthol","Mentha Piperita (Peppermint) Oil","Aloe Barbadensis (Aloe Vera) Leaf Juice","Carthamus Tinctorius (Safflower) Oil","Achillea Millefolium Extract","Chamomilla Recutita (Matricaria) Flower Extract","Equisetum Arvense Extract","Eucalyptus Globulus (Eucalyptus) Oil","Lavandula Angustifolia (Lavender) Leaf Extract","Melaleuca Alternifolia (Tea Tree) Leaf Oil","Rosmarinus Officinalis (Rosemary) Leaf Extract","Salvia Officinalis (Sage) Leaf Extract","Thymus Vulgaris (Thyme) Extract","Tussilago Farfara (Coltsfoot) Flower Extract","Urtica Dioica (Nettle) Extract"]}`,
		},
		{
			name:  "complex",
			input: "bJBPb9tIDMW_CjEnG_BhscfcHOcPdmE3Qeymh6IHekRZBKihwuG4EYp-92KspHGanjTge--nR_4IKVyEHWWH_9LBqGFKnsMicLj4GpZPBWF2X4xbpga-oJPNwyJsteHSw0qjYj90ipEcncIirLHYKHArJWrmht6bJ8Wxn8wvyrsMwQptr8-jTJ4rin_wVhqx50YH02EUuCRHTlW4lTGScQqLcI_JO0oq9a2OOdcfbdX2E3XFbhxhGbk5OWTMkwb__lMHNfs8knc4MTaUvHt7IdzzQMZer0PDQNZz8jnccbUsRQku0fbYUMqcYXaaPJLhHNaELfxfOJ5qoHmHfcmw4xRdjUuG2RbbVvQ72W9g7FiEEDb106rUZa6f3TB6pXTYa88iCA8Ui59qbbBuiMY4h5sT7Sxx_VQ4k5celnaklOlcKxFlHLxkuBXdF6mV3oavndZ4xNQUQVimQ8nOtRbCbI1HSk3tftr0jbshQaESEZbiZOk1sCOEnRG9BCb6g-YejVPJcNe2HDmh1EM-aKYebfxA36IcGd-bt3igD8ZdN9Z7PxY5oFVXHdD83FFyZsGDwg1ai4YwW6l4blX9L7f8bM4R4Yq1fmafyF3OeN9-_goAAP__",
			want:  `{"n":"Test Ingredients","i":["Aqua (Purified Water)","Sodium Cocoamphoacetate","Lauryl Glucoside","Sodium Cocoyl Glutamate","Sodium Lauryl Glucose Carboxylate","Decyl Glucoside","Cocamidopropyl Betaine","Glycerin","Panthenol","Potassium Sorbate","Citric Acid","Polysorbate 20","Phenoxyethanol","Menthol","Mentha Piperita (Peppermint) Oil","Aloe Barbadensis (Aloe Vera) Leaf Juice","Carthamus Tinctorius (Safflower) Oil","Achillea Millefolium Extract","Chamomilla Recutita (Matricaria) Flower Extract","Equisetum Arvense Extract","Eucalyptus Globulus (Eucalyptus) Oil","Lavandula Angustifolia (Lavender) Leaf Extract","Melaleuca Alternifolia (Tea Tree) Leaf Oil","Rosmarinus Officinalis (Rosemary) Leaf Extract","Salvia Officinalis (Sage) Leaf Extract","Thymus Vulgaris (Thyme) Extract","Tussilago Farfara (Coltsfoot) Flower Extract","Urtica Dioica (Nettle) Extract"]}`,
		},
		{
			name:  "search console",
			input: "bZA9TwMxDIb_inXzwcDAwHY9UT4EaoGBpYubc3sWvrhyEuiB-O_42oEPsUSxk_f14_ejitVF1RaTBDMptCaz0W-cEjzSDg0zvxLcEfp5E6HV2HFmjWRVXbFrl8V4w9TBM2ayGmbUU8zGg0YuA9xT7jUV2fhrDS1lQhsFGgnaqxw6v0vLR7uFbTFy-EE1t8IZLvfZMORJ-fdrq0FjybBgN76SMZD59JOz8xqe1NacjwP_yhpRghnaGjuKidO07QZuCwcnvh47UxnffcGHwlHRvVheoBm8gCZwl_5F6XHQgYVgLvpG9o09N9waxsl76UnpfnRGjAc03Hk2HsfErnK6slWsPr8A",
			want:  `{"n":"Curls Blueberry Bliss Reparative Leave In Conditioner","i":"Purified Water, Behentrimonium Methosulfate, Cetearyl Alcohol, Cetyl Alcohol, Certified Organic Blueberry Fruit Extract, Certified Organic Coconut Oil, Glycereth-26, Sorbitol, Certified Organic Aloe Barbadensis Leaf Juice, Hydrolyzed Quinoa, Silk Amino Acids, Certified Organic Chamomile Flower Extract, Fragrance, Phenoxyethanol, Caprylyl Glycol.\r\n"}`,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, Decompress(c.input))
		})
	}
}
