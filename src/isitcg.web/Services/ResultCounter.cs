using StackExchange.Redis;
using System;

namespace isitcg
{
    public class ResultCounter
    {
        private const string PRODUCTS_KEY = "products";
        private const string UNKNOWN_INGREDIENTS_KEY = "ingredients:unknown";
        private readonly IDatabase _db;

        public ResultCounter(IDatabase db)
        {
            _db = db ?? throw new System.ArgumentNullException(nameof(db));
        }

        public void Count(MatchResults results)
        {
            try
            {
                if (!string.IsNullOrEmpty(results.ProductName))
                {
                    _db.SortedSetIncrement(
                         PRODUCTS_KEY, results.ProductName, 1, CommandFlags.FireAndForget);
                }

                foreach (var unknown in results.Remainder)
                {
                    if (!string.IsNullOrEmpty(unknown))
                    {
                        _db.SortedSetIncrement(
                            UNKNOWN_INGREDIENTS_KEY, unknown, 1, CommandFlags.FireAndForget);
                    }
                }
            }
            catch
            {
                //swallow error
            }
        }
    }
}