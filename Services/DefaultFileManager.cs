using System;
using System.IO;
using Newtonsoft.Json;

namespace isitcg
{
    public class DefaultFileManager : IFileManager
    {
        public string Write(MatchResults data)
        {
            var id = generateFileId();
            using (StreamWriter file = File.CreateText($"results/{id}"))
            {
                JsonSerializer serializer = new JsonSerializer();
                serializer.Serialize(file, data);
            }
            return id;
        }

        public MatchResults Read(string id)
        {
            using (StreamReader file = File.OpenText($"results/{id}"))
            {
                JsonSerializer serializer = new JsonSerializer();
                var results = (MatchResults)serializer.Deserialize(file, typeof(MatchResults));
                return results;
            }
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