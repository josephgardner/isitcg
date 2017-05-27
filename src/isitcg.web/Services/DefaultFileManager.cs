using System;
using System.IO;
using System.Threading.Tasks;
using Microsoft.Extensions.Caching.Distributed;
using Newtonsoft.Json;

namespace isitcg
{
    public class DefaultFileManager : IFileManager
    {
        private readonly IDistributedCache _distributedCache;
        public DefaultFileManager(IDistributedCache distributedCache)
        {
            if (distributedCache == null)
                throw new ArgumentNullException(nameof(distributedCache));

            _distributedCache = distributedCache;
        }
        public async Task<string> WriteAsync(MatchResults data)
        {
            var id = generateFileId();
            var json = JsonConvert.SerializeObject(data);
            await _distributedCache.SetStringAsync(id, json);
            return id;
        }

        public async Task<MatchResults> ReadAsync(string id)
        {
            var json = await _distributedCache.GetStringAsync(id);
            var results = JsonConvert.DeserializeObject<MatchResults>(json);
            return results;            
        }

        private string generateFileId()
        {
            var guid = Guid.NewGuid();
            var b64 = Convert.ToBase64String(guid.ToByteArray())
                .Replace("+", "")
                .Replace("=", "")
                .Replace("/", "");
            return b64;
        }
    }
}