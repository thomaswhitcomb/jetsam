# jetsam

Jetsam is a simple, concurrent pipeline that loads multiple ASCII files and places the individual lines from the files onto an internal channel.  The lines from the internal channel are "reduced" or aggregated by a separate function provided by a pipeline client.  The separate function is referred to as the "reducer".

The "reducer" function pulls the individual lines from the source line channel and performs client specific aggreations such as averages, summations and median calculations.

After draining the source line channel, the "reducer" function writes its results to the output channel.

## API

1. Instantiate the jetsam.Pipeline object and specify:
- Array of URLs that point to the source files to be loaded.
-  Buffersize of the sourceLine queue.  This is the queue the reducer reads.  The buffersize is the number of sourceLines that can be read by the reduce without synchronizing.
- The number of Processor threads.  These threads concurrently read from the source URLs.
- The Reducing function.  This function reads from the sourceLine channel, performs the reduction processing and then outputs results on the output channel.
1. Invoke the Provision method on the pipeline.
1. Invoke the Run method on the pipeline.  The Run method performes the concurrent source URL reads, writes the lines to the source line channel and invokes the reducer function to process the source lines.


## Build

    make
