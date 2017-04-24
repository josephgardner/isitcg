
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
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
            var parts = ingredients.Split(',').Select(p => p.Trim().Trim('.'));
            var results = new MatchResults(parts);

            results = _rules.Aggregate(results, (seed, rule) =>
            {
                var lookup = rule.Ingredients;
                var matches = seed.Remainder.Intersect(lookup, StringComparer.OrdinalIgnoreCase).ToList();
                if (matches.Any())
                {
                    seed.Matches.Add(rule.Name, matches);
                    foreach (var match in matches)
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