using System;
using System.Collections;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text.RegularExpressions;
using Microsoft.Extensions.Options;
using Newtonsoft.Json;

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

        public string CreateHash(string productName, string ingredients)
        {
            var product = new Product(productName, ingredients);
            var json = JsonConvert.SerializeObject(product);
            var hash = StringCompression.Compress(json);
            return hash;
        }

        public Product ProductFromHash(string hash)
        {
            var json = StringCompression.Decompress(hash);
            var product = JsonConvert.DeserializeObject<Product>(json);
            return product;
        }

        public MatchResults ResultsFromHash(string hash)
        {
            var product = ProductFromHash(hash);
            var results = ResultsFromProduct(product);
            results.Hash = hash;

            return results;
        }

        internal MatchResults ResultsFromProduct(Product product)
        {
            var results = new MatchResults(product.Parts);
            results.ProductName = product.Name;

            results = _rules.Aggregate(results, (seed, rule) =>
            {
                var lookup = rule.Ingredients;
                var matches1 = seed.Remainder.Intersect(lookup,
                    IngredientComparer.Instance);

                var matches2 = from i in lookup
                               from r in seed.Remainder
                               where r.Contains('/')
                               where r.Split('/').Contains(i)
                               select r;

                var matches3 = matches1.Concat(matches2).ToList();
                if (matches3.Any())
                {
                    seed.Matches.Add(new Rule
                    {
                        Name = rule.Name,
                        Description = rule.Description,
                        BlogUrl = rule.BlogUrl,
                        Result = rule.Result,
                        Rank = rule.Rank,
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
