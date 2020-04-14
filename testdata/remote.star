# file no extension
load("mcuadros/starlark-load.v0/fixtures/module", "foo", "modinfo")
assert.ne(foo, "")
assert.eq(modinfo.name, "github.com/mcuadros/starlark-load.v0")
assert.eq(modinfo.ref, "master")
assert.ne(modinfo.path, "")
assert.ne(modinfo.commit, "")
assert.ne(modinfo.repository, "")
