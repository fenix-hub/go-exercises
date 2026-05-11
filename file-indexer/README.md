`file-indexer` is a simple exercise to experiment with the correct patterns regarding multi-threading, parallelization, synchronization, etc.
The script consists of two types of threads (gorutines) that have to different jobs:
- `crawler` walks through a given directory looking for text files
- a swarm of `indexer` threads read the content of text files and index it, counting the occurrence of words

the swarm synchronizes itself over a channel which is written by the crawler.  
the `Index` consists of a shared map accessed through mutex.  

a `WaitGroup` makes sure the application waits for the completition of all threads.
