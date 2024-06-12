## General

General object parameters:
- Name (string enumeration name)
- Base Type
  - Belt-like
  - Inserter
  - Storage
  - Converter
- Shape
  - Single hex
  - Diamond
  - Bighex

## Belts

Belt tiers:
| tier | speed(ticks/width) | underground reach |
| :--- | :----------------: | :---------------: |
| T1   |         18         |         5         |
| T2   |         9          |         7         |
| T3   |         6          |         9         |

Belt-likes (1x1):
| type     | function    |
| :------- | :---------- |
| Belt     | normal      |
| Under    | underground |
| Splitter | splitting   |


## Other Objects

Inserters:
| type         | swing speed | reach | stack size | filtering |
| :----------- | :---------: | :---: | :--------: | :-------- |
| Burner       |     30      |   1   |     1      | none      |
| Simple       |     20      |   1   |     1      | none      |
| Fast         |     10      |   1   |     1      | none      |
| Long         |     20      |   2   |     1      | none      |
| Stack        |     10      |   1   |     64     | none      |
| Filter       |     10      |   1   |     1      | yes       |
| Stack Filter |     10      |   1   |     64     | yes       |

Storage:
| type     | capacity |
| :------- | :------: |
| Wooden   |    8     |
| Iron     |    16    |
| Steel    |    24    |
| Logistic |    24    |

Converters:
| type       |  shape  | buildpower | set recipe  |
| :--------- | :-----: | :--------: | :---------- |
| Furnace    | diamond |     10     | auto select |
| Assembler1 | bighex  |     8      | manual      |
| Assembler2 | bighex  |     10     | manual      |
| Assembler3 | bighex  |     12     | manual      |


## Recipes

Recipe parameters:
Name
Total buildpower
Ingredients: [item, count]
Energy (kJ)
Result: [item, count]
Converters: [type]

Recipes:

| Results       |      Ingredients       |   BP | Energy |  Converters   |
| :------------ | :--------------------: | ---: | -----: | :-----------: |
| Iron Plate x1 |      Iron Ore x2       |   10 |      5 |    Furnace    |
| Steel Beam x1 |     Iron Plate x5      |   25 |     25 |    Furnace    |
| Gear x1       |     Iron Plate x1      |    5 |      - | Assembler T1+ |
| Belt T1 x1    | Gear x1, Iron Plate x1 |   10 |      - | Assembler T1+ |
