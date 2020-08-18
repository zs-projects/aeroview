# Introduction

Dear reader,   
  
Thank you for granting us some of your precious time and attention.  
  
The following work you are about to read was motivated by a simple observation. Let us consider the common data structures that we use in our day to day jobs. How would their design change if we were to consider a different set of trade-off's and if we optimised for different things.  
  
For example, let's consider array of`integers` for a moment. It is one of the simplest data structure you could consider. What if we wanted to compress it ? could we still retain the random access property ? what time complexity could we achieve ?  If you want to fit a giant inverted index for a full text search database in memory, the answer could be of interest.  

Throughout the book we will look at alternative data layouts and algorithms for common and well known data structures that optimize for different constraints : GC Friendlieness , Fast deserialisation, Lower memory usage. In the first few chapters, we will consider common data structures and alternative designs. We will start with arrays, the we will explores maps and we will finish with trees. In the following chapters, we will use what we learned to  to further improve the designs and explore some intresting applications.  


