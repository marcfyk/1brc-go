# CHANGELOG

This changelog serves to document the iterative progress and optimizations over the 1BRC solution.

---

## 14-04-2024

- The contexts where we split strings on a delimiter guarantees that there aways strictly 1 delimiter in the string and that splits only into 2 substrings.
  Therefore, splitting of strings on a delimiter now use `strings.Cut` instead of `strings.Split`.
  This is because `strings.Cut` outperforms `strings.Split` due on lesser runtime overhead on lower memory allocations.
- The challenge's constraints specify that temperatures are in the range of `[-99.9, 99.9]`, and always with 1 fractional digit,
  therefore temperatures can be represented precisely as integers by multiplying the input temperature by 10.
  Hence, `Temperature` is now an `int` instead of `float64`.
  This improves parsing performance as parsing integers is faster than floating point numbers due to the data format.
  Additionally, integer arithmetic is faster than floating point arithmetic.

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
