using System.Threading.Tasks;

namespace isitcg
{
    public interface IFileManager
    {
        Task<string> WriteAsync(MatchResults data);
        Task<MatchResults> ReadAsync(string id);
    }
}