# Hex Factory

Technical name: hextopdown

This project is aimed to explore game mechanics and implementation details of transport belt driven system with two-lane belts at hexagonal grid.

Popular game Factorio was one of the first to incorporate two-lane transport belts as one of the core game mechanics. Since then many games thied to replicate Factorio's success but most of them avoided complications of implementing two-lane belts.

It turns out those complications are fairly significant, especially at hexagonal grid (Factorio is in orthogonal grid).

This project is a result of several months of work on implementation of two-lane transport belt mechanics on hexagonal grid.

The current state of the project could be described as "early prototype". Mechanics of moving items on two-lane transport belts are fully implemented and tested. Some parts of content is included: inserters can move items from one hex to the other on opposite side of the inserter, furnaces and assemblers can convert one set on items to another, window-based interface works to some extent.
