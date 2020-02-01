namespace isitcg
{
    public interface IIngredientHandler
    {
        string CreateHash(string product, string ingredients);
        MatchResults ResultsFromHash(string hash);
        Product ProductFromHash(string hash);
    }
}