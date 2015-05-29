# CS464
Projects for CS464 (Distributed Systems) at WSU taught by David Bakken during the Spring 2015 semester.

Projects 1 and 2 were coded in Java, while project 3 was coded in Golang. Projects 1 and 2 were run on virtual servers provided by the university for this course, while project 3 was run on nodes on the Deterlab system.

Projects:
- p1: In this project we setup a simple subscriber-publisher message passing system with a single subsriber and a single publisher. We used RTI dds for this project.
- p2: In this project we extended our project 1 code to be a subscriber filtering out messages that it isn't interested in. In this case, the subscriber was interested in messages about high cpu usage or high memory usage on the publisher computer. Due to problems getting RTI's built-in filtering working, I instead went with a flood approach to sending publisher messages and handled filtering subscriber-side (although code for the built-in filtering is in place and the system will use that instead if it can setup the filtering properly).
- p3: In this project we used 4 virtual server nodes running on the Deterlab system to implement a Chandy-Lamport snapshot algorithm on a basic scenario (4 generals trying to balance out troop numbers, with a single commanding general being in control of where to send troops). This was programemd in Golang (which I was unfamiliar with) and I've provided the report that went along with the project as well as the specifications for the project and the paper for the snapshot algorithm.
