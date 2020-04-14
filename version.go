package load

import (
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/mcuadros/go-version"
)

// Versions is a helper to search the suitable tag for a given contraint based
// on a list of git references from a single repository.
type Versions map[string]*plumbing.Reference

// NewVersions returns a Version from a memory.ReferenceStorage.
func NewVersions(refs memory.ReferenceStorage) Versions {
	versions := make(Versions, 0)
	for _, ref := range refs {
		if !ref.Name().IsTag() && !ref.Name().IsBranch() {
			continue
		}

		versions[ref.Name().Short()] = ref
	}

	return versions
}

// Match the needed version againt the references.
func (v Versions) Match(needed string) *plumbing.Reference {
	if version, ok := v[needed]; ok {
		return version
	}

	matched := v.doMatch(needed)
	if len(matched) != 0 {
		return matched[0]
	}

	if needed == "v0" || needed == "0" {
		return v.handleV0()
	}

	return nil
}

func (v Versions) doMatch(needed string) []*plumbing.Reference {
	c := newConstrain(needed)

	var names []string
	for _, ref := range v {
		name := ref.Name().Short()
		if c.Match(version.Normalize(name)) {
			names = append(names, name)
		}
	}

	version.Sort(names)
	var matched []*plumbing.Reference
	for n := len(names) - 1; n >= 0; n-- {
		matched = append(matched, v[names[n]])
	}

	return matched
}

func (v Versions) handleV0() *plumbing.Reference {
	return v.Match("master")
}

func newConstrain(needed string) *version.ConstraintGroup {
	if needed[0] == 'v' && needed[1] >= 28 && needed[1] <= 57 {
		needed = needed[1:]
	}

	return version.NewConstrainGroupFromString(needed + ".*")
}
