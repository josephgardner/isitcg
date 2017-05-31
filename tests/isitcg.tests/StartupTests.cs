using System;
using Xunit;
using isitcg;
using Moq;
using Newtonsoft.Json;
using System.IO;

namespace isitcg.tests
{
    public class StartupTests
    {
        [Fact]
        public void CanLoadIngredientsJson()
        {
            var ingredients = JsonConvert.DeserializeObject<IngredientRules>(
                File.ReadAllText("ingredientrules.json"));
            Assert.NotNull(ingredients);
        }
    }
}