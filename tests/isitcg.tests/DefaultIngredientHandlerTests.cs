using System;
using Xunit;
using isitcg;
using Moq;
using Microsoft.Extensions.Options;
using System.Collections.Generic;
using System.Linq;
using Newtonsoft.Json;
using System.IO;
using Newtonsoft.Json.Linq;
using YamlDotNet.Serialization;
using YamlDotNet.RepresentationModel;
using YamlDotNet.Serialization.NamingConventions;

namespace isitcg.tests
{
    public class RuleTests
    {
        public IEnumerable<RuleTest> Tests { get; set; }
    }
    public class RuleTest
    {
        public string TestName { get; set; }
        public string Ingredients { get; set; }
        public Product Product { get => new Product("test", Ingredients); }
        public IEnumerable<Rule> Rules { get; set; }
        public string ExpectedResult { get; set; }
        public IEnumerable<Rule> ExpectedMatches { get; set; }

        public IEnumerable<string> ExpectedRemainder { get; set; }
        public override string ToString() => TestName;
    }
    public class DefaultIngredientHandlerTests
    {
        public static IEnumerable<object[]> RuleData;
        static DefaultIngredientHandlerTests()
        {
            var input = new StringReader(File.ReadAllText("rule-tests.yml"));
            var deserializer = new DeserializerBuilder().Build();
            var tests = deserializer.Deserialize<RuleTests>(input);

            RuleData = tests.Tests.Select(t => new object[] { t });
        }

        [Fact]
        public void DoesNotThrow()
        {
            //Arrange
            var mockRules = new Mock<IOptions<IngredientRules>>();
            mockRules.Setup(r => r.Value).Returns(new IngredientRules());

            //Act, Assert
            new DefaultIngredientHandler(mockRules.Object);
        }

        [Theory]
        [MemberData(nameof(RuleData))]
        public void RuleMatchesIngredient(RuleTest test)
        {
            //Arrange
            var mockRules = new Mock<IOptions<IngredientRules>>();
            mockRules.Setup(r => r.Value).Returns(new IngredientRules
            {
                Rules = test.Rules
            });
            var handler = new DefaultIngredientHandler(mockRules.Object);

            //Act
            var actual = handler.ResultsFromProduct(test.Product);

            //Assert
            Assert.NotNull(actual);
            Assert.Equal(test.ExpectedResult, actual.Result);
            Assert.Equal(test.ExpectedMatches.Count(), actual.Matches.Count);
            for (int i = 0; i < test.ExpectedMatches.Count(); i++)
            {
                var expectMatch = test.ExpectedMatches.ElementAt(i);
                var actualMatch = actual.Matches.ElementAt(i);
                Assert.Equal(expectMatch.Result, actualMatch.Result);
                Assert.Equal(expectMatch.Ingredients, actualMatch.Ingredients);
            }

            Assert.Equal<string>(test.ExpectedRemainder, actual.Remainder);
        }
    }
}
