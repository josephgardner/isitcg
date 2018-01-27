
using System;
using System.Collections;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text.RegularExpressions;
using Microsoft.Extensions.Options;

namespace isitcg
{
    internal class DefaultIngredientHandler : IIngredientHandler
    {
        private readonly IEnumerable<Rule> _rules;
        public DefaultIngredientHandler(IOptions<IngredientRules> rules)
        {
            if (rules == null)
                throw new ArgumentNullException(nameof(rules));

            _rules = rules.Value.Rules;
        }
        public MatchResults CreateResults(string ingredients)
        {
            var parts = ingredients
                .Split(',')
                .Select(p => p.Trim().Trim('.'));
            var results = new MatchResults(parts);

            results = _rules.Aggregate(results, (seed, rule) =>
            {
                var lookup = rule.Ingredients;
                var matches1 = seed.Remainder.Intersect(lookup, 
                    IngredientComparer.Instance);
                
                var matches2 =  from i in lookup
                                from r in seed.Remainder
                                where r.Contains('/')
                                where r.Split('/').Contains(i)
                                select r;

                var matches3 = matches1.Concat(matches2).ToList();
                if (matches3.Any())
                {
                    seed.Matches.Add(new Rule{
                        Name = rule.Name,
                        Description = rule.Description,
                        Result = rule.Result,
                        Ingredients = matches3
                    });
                    if (rule.Result == "danger")
                        results.Result = "danger";
                    else if (rule.Result == "warning" && results.Result == "good")
                        results.Result = "warning";
                    
                    foreach (var match in matches3)
                    {
                        seed.Remainder.Remove(match);
                    }
                }
                return seed;
            });

            return results;
        }
    }
}
