# Wikipedia Path Finder

This algorithm finds paths in between two wikipedia pages. I'm still trying
to figure out how to improve the algorithm performance, but currently, it can
find a path from the top 1 to the top 25_000 pages in 20 seconds on a Ryzen 7
4800H.

The main bottleneck right now is on the binary search algorithm. The goal of this
project was to help me learn cache-locality, profiling and optimizations in Golang.
