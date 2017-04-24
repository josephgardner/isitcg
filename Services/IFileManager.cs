namespace isitcg
{
    public interface IFileManager
    {
        string Write(MatchResults data);
        MatchResults Read(string id);
    }
}