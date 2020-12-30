#!/bin/bash

DAY=$1

curl -b ~/.aoc_cookies https://adventofcode.com/2019/day/$DAY/input -o inputs/$DAY.input

curl -b ~/.aoc_cookies https://adventofcode.com/2019/day/$DAY -o - | \
  awk '
    BEGIN{num=0;}
    /<pre><code>/{gsub("<pre><code>",""); on=1;}
    /<\/code><\/pre>/{on=0;num++;}
    on{file = "inputs/'$DAY'sample" num ".input"; print > file;}
  '
