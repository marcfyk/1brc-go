# CHANGELOG

This changelog serves to document the iterative progress and optimizations over the 1BRC solution.

---

## 08-04-2024

- Implement naive solution as a starting baseline.
    - I/O on the input file is done by scanning line by line.
    - station names are `strings` with a new type definition, `Station`, and temperatures are `float64` with a new type definition, `Temperature`.
    - Station data is in a struct `Info` keeps track of:
        - counts of that station's temperature in the input file `uint`.
        - sum of temperatures of the station's temperature as `Temperature`
        - min temperature observed as a `Temperature`.
        - max temperature observed as a `Temperature`.
    - Data for all stations to their related `Info` is stored in native go `map`.
    - Mean temperature of stations only when writing results to stdout and after all input data is processed.
    - Single threaded solution.
