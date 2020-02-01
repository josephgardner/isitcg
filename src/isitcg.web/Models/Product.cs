using System.Collections.Generic;
using System.Linq;
using Newtonsoft.Json;

namespace isitcg
{
    public class Product
    {
        public Product() : this("", "")
        {

        }
        public Product(string productName, string ingredients)
        {
            this.Name = productName;
            this.Ingredients = ingredients;
        }

        [JsonProperty("n")]
        public string Name { get; set; }
        [JsonProperty("i")]
        public string Ingredients { get; set; }
        [JsonIgnore]
        public IEnumerable<string> Parts
        {
            get => Ingredients
                .Split(',')
                .Select(p => p.Trim().Trim('.'));
        }

    }
}