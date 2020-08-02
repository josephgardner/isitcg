using System;
using System.Collections.Generic;
using System.IO;
using System.Security.Claims;
using System.Threading;
using System.Threading.Tasks;
using isitcg;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Http.Features;
using Microsoft.AspNetCore.Rewrite;
using Moq;
using Xunit;

namespace isitcg.tests
{
    public class RewriteExtensionsTests
    {
        [Fact]
        public void DoesNotThrow()
        {
            //Act, Assert
            new RedirectToDomainRule("old", "new");
        }

        [Fact]
        public void RequestNewDomain_ContinuesRules()
        {
            //Arrange
            var rule = new RedirectToDomainRule("old", "new");
            var request = new Mock<HttpRequest>();
            var response = new Mock<HttpResponse>();
            var headers = new HeaderDictionary();
            response.SetupGet(r => r.Headers).Returns(headers);
            response.SetupProperty(r => r.StatusCode);
            request.SetupGet(r => r.Host).Returns(new HostString("new"));
            var context = new Mock<HttpContext>();
            context.SetupGet(c => c.Request).Returns(request.Object);
            context.SetupGet(c => c.Response).Returns(response.Object);
            var rewriteContext = new RewriteContext
            {
                HttpContext = context.Object
            };

            //Act
            rule.ApplyRule(rewriteContext);

            //Assert
            Assert.Equal(RuleResult.ContinueRules, rewriteContext.Result);
            Assert.Equal(response.Object.StatusCode, 0);
            Assert.Equal(response.Object.Headers.Count, 0);
        }

        [Theory]
        [InlineData("http", "", "", "", "http://www.isitcg.com/")]
        [InlineData("http", "", "/Home/Results/pglyHUjEaaz9QtRn9sQ", "", "http://www.isitcg.com/Home/Results/pglyHUjEaaz9QtRn9sQ")]
        public void RequestOldDomain_Redirects(string scheme, string pathBase, string path, string queryString, string expectedUrl)
        {
            //Arrange
            var rule = new RedirectToDomainRule("isitcg.herokuapp.com", "www.isitcg.com");
            var request = new Mock<HttpRequest>();
            request.SetupGet(r => r.Host).Returns(new HostString("old"));
            request.SetupGet(r => r.Scheme).Returns(scheme);
            request.SetupGet(r => r.PathBase).Returns(new PathString(pathBase));
            request.SetupGet(r => r.Path).Returns(new PathString(path));
            request.SetupGet(r => r.QueryString).Returns(new QueryString(queryString));
            var response = new Mock<HttpResponse>();
            response.SetupProperty(r => r.StatusCode);
            var headers = new HeaderDictionary();
            response.SetupGet(r => r.Headers).Returns(headers);
            var context = new Mock<HttpContext>();
            context.SetupGet(c => c.Request).Returns(request.Object);
            context.SetupGet(c => c.Response).Returns(response.Object);
            var rewriteContext = new RewriteContext
            {
                HttpContext = context.Object
            };

            //Act
            rule.ApplyRule(rewriteContext);

            //Assert
            Assert.Equal(RuleResult.EndResponse, rewriteContext.Result);
            Assert.Equal(301, response.Object.StatusCode);
            Assert.Equal(1, response.Object.Headers.Count);
            Assert.Equal(expectedUrl, response.Object.Headers["Location"]);
        }
    }
}