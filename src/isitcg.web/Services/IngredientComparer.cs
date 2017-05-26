using System;
using System.Collections.Generic;
using System.Text.RegularExpressions;

namespace isitcg
{
    public class IngredientComparer : IEqualityComparer<string>
    {
        public static IngredientComparer Instance => new IngredientComparer();
        private static readonly Regex _regex = new Regex(@"(\[.*?\])|(\(.*?\))|-|\*|\s");

        public bool Equals(string strx, string stry)
        {
            if (strx == null)
                return string.IsNullOrWhiteSpace(stry);
            else if (stry == null)
                return string.IsNullOrWhiteSpace(strx);

            string a = _regex.Replace(strx, "");
            string b = _regex.Replace(stry, "");
            return String.Compare(a, b, true) == 0;
        }

        public int GetHashCode(string obj)
        {
            if (obj == null)
                return 0;

            string a = _regex.Replace(obj, "");
            return a.ToLower().GetHashCode();
        }
    }
}