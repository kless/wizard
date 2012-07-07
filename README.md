gowizard
========

Tired of adding the same files every time you create a new Go project?  
Don't know how to structure it?  
Don't know how to apply the license?

http://go.pkgdoc.org/github.com/kless/gowizard/gowizard


## Installation

	go get github.com/kless/gowizard/gowizard

To only install the package, which could be used by a Go IDE:

	go get github.com/kless/gowizard


## Suggestions

### Maintenance

Copyright notices only need the year when the file was created, so don't add new
years.

### Licenses

My suggestion is to use MPL 2.0 because it allows covered source code to be mixed
with other files under a different, even proprietary license. However, code
files licensed under the MPL must remain under the MPL and freely available in
source form.

The BSD-like licenses have been excluded because they can not mitigate threats
from software patents and LGPL because it has not benefits using into a language
of static linking.

#### GPL/AGPL 3.0

Proprietary software linking: Not allowed.
Distribution with code under another license: Not allowed with software whose license is not GNU GPL compatible.
Redistributing of the code with changes: Only under GNU GPL/AGPL.

#### Apache Public 2.0

Proprietary software linking: Allowed.
Distribution with code under another license: Allowed.
Redistributing of the code with changes: Allowed.
Compatible with GNU GPL: Yes.

#### MPL 2.0

Proprietary software linking: Allowed.
Distribution with code under another license: Allowed.
Redistributing of the code with changes: Only under MPL.
Compatible with GNU GPL: Yes


## Copyright and licensing

*Copyright 2010  The "gowizard" Authors*. See file AUTHORS and CONTRIBUTORS.  
Unless otherwise noted, the source files are distributed under the
*Mozilla Public License, version 2.0* found in the LICENSE file.

