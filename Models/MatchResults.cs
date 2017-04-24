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
        }
        public IList<string> Remainder { get; private set; }
        public IDictionary<string, IList<string>> Matches { get; } =
             new Dictionary<string, IList<string>>();
    }
}