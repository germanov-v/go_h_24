using System.Collections;
using System.IO.Compression;
using System.Runtime.CompilerServices;
using System.Text.Json;
using System.Text.RegularExpressions;

namespace StreamingJsonFile;

public record Person(string Email);

public class JsonFileParse
{
    private Dictionary<string, int> _stats = new();

    private readonly Regex DomainRegex = new(@"\.(.+)", RegexOptions.Compiled | RegexOptions.IgnoreCase);


    public async IAsyncEnumerable<Person?> ParseFile(string pathFile
        , [EnumeratorCancellation] CancellationToken cancellationToken)
    {
        using var streamReader = new StreamReader(pathFile);

        // foreach (var item in (await streamReader.ReadLineAsync(cancellationToken))!)
        // while(streamReader.EndOfStream == false)
        while (true)
        {
            var line = await streamReader.ReadLineAsync(cancellationToken);
            if (line == null) yield break;
            yield return JsonSerializer.Deserialize<Person>(line!);
        }

        // await foreach (var item in   streamReader.ReadLineAsync( cancellationToken))
        // {
        //     yield return JsonSerializer.Deserialize<Person>(item);
        // }
    }

    public async Task CalculateStat(string pathFile
        , CancellationToken cancellationToken, bool forceLarge = false, bool zipFile = false)
    {
        var parser = forceLarge
            ? ParseFileLarge(pathFile, cancellationToken,zipFile)
            : ParseFile(pathFile, cancellationToken);

        await foreach (var item in parser.WithCancellation(cancellationToken))
        {
            if (item != null)
            {
                var parts = item.Email.Split("@");

                if (parts.Length != 2) continue;
                //  var domainRegex = new Regex(@"\."+Regex.Escape( item.Email), RegexOptions.Compiled);
                //  if (domainRegex.IsMatch(parts[1]))

                if (!DomainRegex.IsMatch(parts[1])) continue;

                var domainStr = parts[1];

                if (!_stats.TryAdd(domainStr, 1))
                {
                    _stats[domainStr]++;
                }
            }
        }
    }

    public async IAsyncEnumerable<Person?> ParseFileLarge(string pathFile
        , [EnumeratorCancellation] CancellationToken cancellationToken, bool zipFile = false)
    {

        if (zipFile)
        {
            using var archive = ZipFile.OpenRead(pathFile);
            var entry = archive.Entries.FirstOrDefault();
            if (entry == null) throw new Exception("No files in archive.");

            await using var stream = entry.Open();
            await using var bufferedStream = new BufferedStream(stream, 8192);
            using var streamReader = new StreamReader(bufferedStream);

            while (true)
            {
                var line = await streamReader.ReadLineAsync(cancellationToken);
                if (line is null) yield break;

                yield return JsonSerializer.Deserialize<Person>(line);
            }
        }
        else
        {
            await using var fileStream = new FileStream(pathFile, FileMode.Open, FileAccess.Read, FileShare.Read, 4096,
                FileOptions.SequentialScan);
            await using var bufferedStream = new BufferedStream(fileStream, 8192);
            using var streamReader = new StreamReader(bufferedStream);

            while (true)
            {
                var line = await streamReader.ReadLineAsync(cancellationToken);
                if (line is null) yield break;

                yield return JsonSerializer.Deserialize<Person>(line);
            }
        }
        
       
    }
}