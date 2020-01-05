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
        public string ProductName { get; set; }
        public string Result { get; set; }
        public IList<string> Remainder { get; }
        private readonly SortedSet<Rule> _matches = new SortedSet<Rule>(new RuleComparer());
        public ICollection<Rule> Matches
        {
            get { return _matches; }
        }
    }
}