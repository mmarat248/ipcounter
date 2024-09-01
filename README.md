# Counting IP Addresses

## Notes

For this task, using a bitmap counter would be sufficient. 
The size of the bitmap counter depends on the number of possible IP addresses. 
In IPv4, there are 2^32 (around 4.3 billion) possible IP addresses. Therefore, to represent all possible IPv4 addresses, 
a bitmap counter of size 2^32 bits, or 512 MB, would be required.

However, for my own interest, I implemented the HyperLogLogPlus algorithm based on the article:
https://static.googleusercontent.com/media/research.google.com/en//pubs/archive/40671.pdf.
In general, I tried to follow the concepts from the article, and I also implemented several benchmarks for different 
counter variants. But there's probably room for further optimization.

For reading files, I decided to use syscall.Mmap instead of the standard file reading approach. 
Note that this imposes some limitations. Using mmap can lead to data corruption or data loss if the mapped file is 
modified externally while being accessed through the memory mapping, as the changes may not be reflected in the mapped 
memory region, and it's not supported consistently across all platforms. 
In production, I would think twice to use mmap as it imposes significant limitations and complexities when processing 
files in chunks, and additional complex testing is required. 
However, within this task, I was interested in trying out this mechanism since I had not worked with it before.

Also, please note that concurrency could have been implemented more optimally, as we use buckets in each counter, 
which easily allows for concurrent counting. 
Alternatively, we could have created separate counters for each thread, but this would require implementing a merge 
function for each counter. However, the behavior of the counters is deterministic, making the merge operation reliable.

Additionally, I intentionally omitted some checks for IP address correctness and similar validations to keep the 
code concise. For a production version, it would be desirable to refactor the code slightly, add more checks, 
and include more tests.


## How to Run
This program is a command-line utility that counts the number of unique IP addresses in a file 
using either the HyperLogLogPlus or Bitmap algorithm. 

Follow these steps to run the program:

Run the command: 
```
go run main.go -file /path/to/file.txt -counter bitmap
```
Replace /path/to/file.txt with the actual path to the file containing the IP addresses you want to count. 
You can also use the -counter flag to specify the algorithm to use. The available options are bitmap, hyperloglog and 
hyperloglogplus.

Example: 
```
go run main.go -file ./ipcounter/ipsbig -counter bitmap

go run main.go -file ./ip_addresses -counter bitmap
go run main.go -file ./ip_addresses -counter hyperloglog
go run main.go -file ./ip_addresses -counter hyperloglogplus
```
This command will count the unique IP addresses in the ipsbig file located in the testdata directory using the Bitmap algorithm.
View the Output: The program will print the count of unique IP addresses to the console using the selected algorithm.


## Self-Reflection

Upon reflection, there are several areas where this implementation could be improved:

- Performance optimization: There might be room for further improvements, especially in terms of memory usage and
  processing speed.
- Concurrency: As mentioned in the notes, the concurrency implementation could be optimized to take full advantage
  of the bucket structure in the counters.
- Testing: More comprehensive unit tests could be added to ensure the correctness of both the HyperLogLogPlus and Bitmap
  implementations.
- Error handling: The current implementation could benefit from more robust error handling,
  especially when dealing with file operations and memory mapping.
- Validation: Adding more thorough input validation, especially for IP address correctness, would make the program
  more robust and reliable.
- Add more control: For calculating IP via bitmap HashFunc is not required, but for calculating via an approximate
  option - it is desirable to have it. etc...
