using System;
using System.Collections;
using System.Collections.Generic;

namespace isitcg
{
    public class RuleComparer : IComparer<Rule>
    {
        public int Compare(Rule x, Rule y)
        {
            var rank = x.Rank.CompareTo(y.Rank);
            return rank == 0 ? 1 : rank;
        }
    }
}