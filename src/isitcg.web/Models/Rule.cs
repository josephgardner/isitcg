using System.Collections.Generic;

namespace isitcg
{
    public class Rule
    {
        public string Name { get; set; }
        public string Description { get; set; }
        public string Result { get; set; }
        public string BlogUrl { get; set; }
        public IEnumerable<string> Ingredients { get; set; }
    }
}