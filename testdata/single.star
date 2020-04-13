load("fixtures/module", "foo")
first = foo

load("fixtures/module", "foo")
second = foo

assert.eq(first, second)