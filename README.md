Delta
=====

Delta is a simple package for calculating differences between variants of text on a single word level. It can output HTML or plain text. It is built off of the pseudocode for the Longest Common Subsequence problem found on [Wikipedia](http://en.wikipedia.org/wiki/Longest_common_subsequence_problem)

Examples
--------

    delta.Calculate("hello world", "hello earth", false)
        // "hello <del>world</del> <ins>earth</ins>"

    delta.Calculate("hello world", "hello earth", true)
        // "hello ---world--- +++earth+++"
