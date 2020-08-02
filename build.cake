var target = Argument("target", "Default");
var tag = Argument("tag", "cake");

Task("Restore")
  .Does(() =>
{
    DotNetCoreRestore(".");
});

Task("Build")
  .IsDependentOn("Restore")
  .Does(() =>
{
    var settings = new DotNetCoreBuildSettings
    {
        Framework = "netcoreapp3.1",
        Configuration = "Release"
    };

    DotNetCoreBuild("./src/isitcg.web/isitcg.csproj", settings);
    DotNetCoreBuild("./tests/isitcg.tests/isitcg.tests.csproj", settings);
});

Task("Test")
  .IsDependentOn("Build")
  .Does(() =>
{
    var files = GetFiles("tests/**/*.csproj");
    foreach(var file in files)
    {
        var settings = new DotNetCoreTestSettings
        {
            Configuration = "Release"
        };

        DotNetCoreTest(file.ToString(), settings);
    }
});

Task("Default")
    .IsDependentOn("Restore")
    .IsDependentOn("Build")
    .IsDependentOn("Test");

RunTarget(target);