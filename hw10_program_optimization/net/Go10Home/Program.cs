using System.Diagnostics;

namespace StreamingJsonFile;



class Program
{
    
    private const long Mb = 1<<20; // 1 MB 
    private const long MemoryLimit = 30 * Mb;
    private static readonly TimeSpan TimeLimit = TimeSpan.FromMicroseconds(300);
    private static readonly long TimeLimitMilliseconds = 300*1_000;
    
    static async Task Main(string[] args)
    {
        GC.Collect();       // on demand
        GC.WaitForPendingFinalizers(); // wait
        GC.Collect();       //on demand once more time)

        long memoryBefore = GC.GetTotalMemory(true);
        var stopwatch = new Stopwatch();
        var parser = new JsonFileParse();
        var src = new CancellationTokenSource();
        stopwatch.Start();
        
        // await parser.CalculateStat("users.dat", src.Token, false).ConfigureAwait(false);
        await parser.CalculateStat("users.dat.zip", src.Token, true, true).ConfigureAwait(false);

        stopwatch.Stop();
        
        GC.Collect();       // on demand
        GC.WaitForPendingFinalizers(); // wait
        GC.Collect();       //on demand once more time)
        long memoryAfter = GC.GetTotalMemory(true);
        
        var factTime = stopwatch.ElapsedMilliseconds * 1000;
        
        var deltaMemory = (memoryAfter - memoryBefore);
        if (deltaMemory <= MemoryLimit)
        {
            Console.WriteLine($"Success!!! Memory:  {deltaMemory}");
        }
        else
        {
            Console.WriteLine($"Fail!!! Memory:  {deltaMemory}");
        }
        
        if (factTime < TimeLimitMilliseconds)
        {
            Console.WriteLine($"Success!!! Elapsed time: {factTime}");
        }
        else
        {
            Console.WriteLine($"Fail!!! Elapsed time: {factTime}");
        }
        
        // Console.WriteLine(factTime);
        // Console.WriteLine(factTime.TotalMicroseconds);
        // Console.WriteLine(TimeLimit);
        
    }
    
    
    
}