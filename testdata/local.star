load("assert.star", "assert")

# file no extension
load("fixtures/module", "foo", "modinfo")
assert.ne(foo, "")
assert.eq(modinfo.name, "fixtures/module")

# file with extension
load("fixtures/module.star", "foo", "modinfo")
assert.ne(foo, "")
assert.eq(modinfo.name, "fixtures/module")

# directory
load("fixtures/directory", "foo", "modinfo")
assert.ne(foo, "")
assert.eq(modinfo.name, "fixtures/directory")

# directory with main
load("fixtures/directory/main.star", "foo", "modinfo")
assert.ne(foo, "")
assert.eq(modinfo.name, "fixtures/directory")

# file in path
load("module", "foo", "modinfo")
assert.ne(foo, "")
assert.eq(modinfo.name, "module")
