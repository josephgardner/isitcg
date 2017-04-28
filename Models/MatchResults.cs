using System;
using System.Collections.Generic;
using System.Linq;
using Microsoft.Extensions.Options;

namespace isitcg
{
    public class MatchResults
    {
        public MatchResults(IEnumerable<string> remainder)
        {
            if (remainder == null)
                throw new ArgumentNullException(nameof(remainder));

            Remainder = remainder.ToList();
            Result = "good";
        }
        public string Result { get; set; }
        public IList<string> Remainder { get; private set; }
        public IList<Rule> Matches { get; set; } = new List<Rule>();
    }
}