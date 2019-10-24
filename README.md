# Generate character sequences from POSIX Bracket Expressions

> a bracket expression matches any character among those listed between the
> opening and closing square brackets.  Within a bracket expression, a range
> expression consists of two characters separated by a hyphen. It matches any
> single character that sorts between the two characters, based upon the
> system’s native character set. For example, ‘[0-9]’ is equivalent to
> ‘[0123456789]’ ~ [GNU Awk
> manual](https://www.gnu.org/software/gawk/manual/html_node/Bracket-Expressions.html)

Bracket expressions are often used to limit searches, regular expressions and
filters. But sometimes you need to have the range of characters expanded and
explicitly listed; this library helps do that.

```go
gen, _ := bechars.New()
rng, _ := gen.Generate("[:print:]")
fmt.Println(rng) // => " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"

rng, _ = gen.Generate("[\u0e010-2]")
fmt.Println(rng) // => "ก012"

rng, _ = gen.Generate("[:punct:]")
fmt.Println(rng) // => "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~]"
```

The generate can also be configured:

```go
gen, _ := bechars.New(MinRune('a'), MaxRune('z'))
rng, _ := gen.Generate("[^:cntrl::punct:]")
fmt.Println(rng) // => "abcdefghijklmnopqrstuvwxyz"
```
