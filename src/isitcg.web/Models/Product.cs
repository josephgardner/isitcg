using System.Collections.Generic;
using System.Linq;

namespace isitcg
{
    public class Product
    {
        public Product(string productName, string ingredients)
        {
            this.Name = productName;
            this.Parts = ingredients
                .Split(',')
                .Select(p => p.Trim().Trim('.'));
        }
        public string Name { get; set; }
        public IEnumerable<string> Parts { get; set; }
    }
}