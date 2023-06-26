# cue quest

PARKED FOR NOW

# Struct types
## Restriction
Struct fields may only have type names, not values in the declaration.

**Valid:**

```
#a : {
    a : int
}
```

**Invalid:**
```
#a : {
    a : 5
}
```

# Map types
Should we define map types in the specification language?

How will we implement them in C

# TODO
* Handle null values
* Should we allow union messages variants with non-declared type?