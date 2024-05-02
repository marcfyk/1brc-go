# CHANGELOG

This changelog serves to document the iterative progress and optimizations over the 1BRC solution.

---

## 02-05-2024

- Optimizations
    - Temperature parsing is directly computed based on the offset of the delimiter, `;`, from the end of the line.
      This is more efficient than manually iterating through the bytes representing the temperature after `;` to calculate the temperature.

## 21-04-2024

- Bug fixes
    - Calculation of temperature mean was not rounding correctly due to integer division of sum and count.
      For example, an average of `1.57...` should be rounded to `1.6` but was incorrectly rounded to `1.5`.
      This issue has been fixed with temperature mean rounding up/down correctly based on the rounding digit.
- Optimizations:
    - `Temperature` is changed from an `int` to `int16` to occupy less memory due to `Temperature` being bounded to `[-999, 999]`.
    - `TemperatureSum` type is added as an `int64` to track the sum of `Temperature` as it has accept a range of `[1_000_000_000 * 999, 1_000_000_000 * -999]`.
    - `Count` is changed from `uint` to `uint32` to prevent unnecessary allocation of `uint64` on 64bit systems as only `uint32` is needed to accept a range of `[0, 1_000_000_000]`.
    - Measurement parsing is now done directly on the byte slice line from the file's buffered scanner to reduce unnecessary memory allocations.
      Previously, the line is read as a string before being subsequent string operations are performed to parse each segment of the measurement, which performs many unnecessary allocations.
      Now, all parsing of each segment of the measurement is done directly on the byte slice by exploiting the constraints of the input.
    - Writing to output now uses a single allocated buffer using a computed size to reduce unnecessary memory allocations.
      Previously, output was written by first allocating strings for each (station, min, mean, max) set before joining these strings into another allocated string joining them with newlines.
      Now, a single allocated buffer is used to write all output before allocating once into a string.
      The size of the buffer is sufficient for the entire output by estimating the size of the output.
      The estimation of the output is done by computing the sum of station names, the necessary newlines and `;` delimiters, and assuming each temperature is 5 digits long, which is the max length of a temperature.
      The max length of the temperature is 5 when it is a negative temperature less than `-9.9`, such as `-10.7` and `-99.9`.

## 14-04-2024

- Optimizations
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
