# Distributed_Mutual_Exclusion

## How to run

in seperate terminals, run a server and as many clients you want (less than 10 is recomended)

server:
go run .\Server\

client:
go run .\Client\ -sender [name]



## Description:

You have to implement distributed mutual exclusion between nodes in your distributed system. 

You can choose to implement any of the algorithms, that were discussed in lecture 7.

## System Requirements:

R1: Any node can at any time decide it wants access to the Critical Section

R2: Only one node at the same time is allowed to enter the Critical Section 

R2: Every node that requests access to the Critical Section, will get access to the Critical Section (at some point in time)

## Technical Requirements:

Use Golang to implement the service's nodes
Use gRPC for message passing between nodes
Your nodes need to find each other.  For service discovery, you can choose one of the following options
 supply a file with  ip addresses/ports of other nodes
enter ip adress/ports trough command line
use the Serf package for service discovery
Demonstrate that the system can be started with at least 3 nodes
Demonstrate using logs,  that a node gets access to the Critical Section

# Hand-in requirements:

Hand in a single report in a pdf file
Provide a link to a Git repo with your source code in the report
Include system logs, that document the requirements are met, in the appendix of your report

# Grading notes

Partial implementations may be accepted, if the students can reason what they should have done in the report.
