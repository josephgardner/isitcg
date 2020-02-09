using StackExchange.Redis;

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
            _db.SortedSetIncrement(
                 PRODUCTS_KEY, results.ProductName, 1, CommandFlags.FireAndForget);

            foreach (var unknown in results.Remainder)
            {
                _db.SortedSetIncrement(
                    UNKNOWN_INGREDIENTS_KEY, unknown, 1, CommandFlags.FireAndForget);
            }
        }
    }
}