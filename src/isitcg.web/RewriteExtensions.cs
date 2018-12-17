using System;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Http.Extensions;
using Microsoft.AspNetCore.Rewrite;
using Microsoft.Net.Http.Headers;

namespace isitcg
{
    internal static class RewriteExtensions
    {
        public static void AddDomainRedirect(this RewriteOptions options, string oldDomain, string newDomain)
        {
            options.Rules.Add(new RedirectToDomainRule(oldDomain, newDomain));
        }
    }

    internal class RedirectToDomainRule : IRule
    {
        private readonly string _oldDomain, _newDomain;
        public RedirectToDomainRule(string oldDomain, string newDomain)
        {
            _oldDomain = oldDomain ?? throw new ArgumentNullException(nameof(oldDomain));
            _newDomain = newDomain ?? throw new ArgumentNullException(nameof(newDomain));
        }
        public void ApplyRule(RewriteContext context)
        {
            var req = context.HttpContext.Request;
            if (req.Host.Host.Equals(_newDomain, StringComparison.OrdinalIgnoreCase))
            {
                return;
            }

            var wwwHost = new HostString(_newDomain);
            var newUrl = UriHelper.BuildAbsolute(req.Scheme, wwwHost, req.PathBase, req.Path, req.QueryString);
            var response = context.HttpContext.Response;
            response.StatusCode = 301;
            response.Headers[HeaderNames.Location] = newUrl;
            context.Result = RuleResult.EndResponse;
        }
    }
}