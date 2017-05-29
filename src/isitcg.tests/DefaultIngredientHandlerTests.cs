using System;
using Xunit;
using isitcg;
using Moq;
using Microsoft.Extensions.Options;
using Microsoft.AspNetCore.Hosting;

namespace isitcg.tests
{
    public class DefaultIngredientHandlerTests
    {
        [Fact]
        public void DoesNotThrow()
        {
            //Arrange
            var mockRules = new Mock<IOptions<IngredientRules>>();
            mockRules.Setup(r => r.Value).Returns(new IngredientRules());

            //Act, Assert
            new DefaultIngredientHandler(mockRules.Object);
        }
    }
}
