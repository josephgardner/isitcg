using System;
using System.Collections;
using System.Collections.Generic;

namespace isitcg
{
    public class RuleComparer : IComparer<Rule>
    {
        public int Compare(Rule x, Rule y)
        {
            if (x.Result == y.Result)
                return -1;
            if (x.Result == "danger")
                return -1;
            if (y.Result == "danger")
                return 1;
            if (x.Result == "warning")
                return -1;
            return 1;
        }
    }
}