* Go Documentation Rules (AGENT.org)
This document defines the =Go documentation standards= for this project,
based on the [[https://go.dev/doc/comment][official Go documentation
guidelines]]. These rules are designed to work with =ZED + Codestral=
and must be followed by all contributors.

** General Principles
- Always use =English= for documentation.
- Write complete sentences with proper punctuation.
- Be concise, precise, and avoid redundancy.
- Use backticks for code references (e.g., =See [storage.Upload]=).

* Table of Contents
- [[#package-documentation][Package Documentation]]
- [[#function-and-method-documentation][Function and Method
  Documentation]]
- [[#type-and-struct-documentation][Type and Struct Documentation]]
- [[#variables-and-constants][Variables and Constants]]
- [[#code-examples][Code Examples]]
- [[#using-go-doc][Using =go doc=]]
- [[#best-practices][Best Practices]]
- [[#enforcement][Enforcement]]

* Package Documentation
Every package /must/ start with a documentation comment.

*** Format
```go
// Package <name> <concise and complete description>.
//
// <Additional explanations if needed.>
package <name>
```

*** Example
```go
// Package storage provides tools for storing and managing images in MinIO.
// It includes functions for uploading, deleting, and retrieving metadata.
package storage
```

* Function and Method Documentation
*** Rule
Every /exported/ function/method /must/ have a descriptive comment.

*** Format
```go
// <Name> <description of what the function/method does>.
// <Details about parameters, return values, and possible errors.>
func <Name>(<parameters>) (<returns>) {
    // ...
}
```

*** Example
```go
// ResizeImage resizes an image to the specified dimensions.
// It returns the resized image or an error if resizing fails.
// Parameters:
//   - src: source image in image.Image format.
//   - width, height: new dimensions in pixels.
// Returns:
//   - image.Image: resized image.
//   - error: potential error (e.g., unsupported format).
func ResizeImage(src image.Image, width, height int) (image.Image, error) {
    // ...
}
```

* Type and Struct Documentation
*** Rule
All /exported/ types and structs /must/ be documented.

*** Format
```go
// <Name> represents <description of the type/struct>.
// <Details about fields if necessary.>
type <Name> struct {
    // <FieldName> <description of the field>.
    <FieldName> <Type>
}
```

*** Example
```go
// ImageMetadata contains metadata extracted from an image.
// Fields include EXIF data and dimensions.
type ImageMetadata struct {
    // Width and Height are the image dimensions in pixels.
    Width  int
    Height int
    // EXIF contains raw EXIF data.
    EXIF map[string]string
}
```

* Variables and Constants
*** Rule
All /exported/ variables and constants /must/ be documented.

*** Format
```go
// <Name> is <description of the variable/constant>.
var <Name> <Type> = <value>
```

*** Example
```go
// DefaultQuality is the default quality for resizing (75%).
const DefaultQuality = 75
```

* Code Examples
*** Rule
Add /executable examples/ for complex functions/methods (optional but
recommended).

*** Format
```go
// Example:
//  // package main
//  //
//  // func main() {
//  //     img, _ := ResizeImage(src, 800, 600)
//  //     _ = img
//  // }
func ResizeImage(...) {
    // ...
}
```

* Using =go doc=
*** Rule
Documentation /must/ be readable and informative via =go doc=.

*** Expected Output
#+begin_example
$ go doc storage.ResizeImage
func ResizeImage(src image.Image, width, height int) (image.Image, error)
    ResizeImage resizes an image to the specified dimensions.
    It returns the resized image or an error if resizing fails.
    Parameters:
      - src: source image in image.Image format.
      - width, height: new dimensions in pixels.
    Returns:
      - image.Image: resized image.
      - error: potential error (e.g., unsupported format).
#+end_example

* Best Practices
- Avoid redundancy: Do not repeat the function signature in the comment.
- Be specific: Document edge cases, errors, and expected behavior.
- Use links: Reference related functions/types with backticks (e.g.,
  =See [storage.Upload]=).

* Enforcement
*** Linting
Use [[https://golangci-lint.run/][golangci-lint]] with the following
linters:

- =godot= (checks for proper sentence formatting)
- =misspell= (checks for spelling errors)
- =revive= (checks for documentation coverage)

*** Editor Setup (ZED)
- Create /snippets/ for documentation templates.
- Configure /Codestral/ to suggest or generate documentation comments.

*** Code Reviews
- Verify that all exported symbols are documented.
- Ensure documentation is clear, complete, and follows the format above.

* Tools and Automation
- Linting: Run =golangci-lint run= before committing.
- Documentation Preview: Use =go doc -all= to preview package
  documentation.
- ZED Integration: Use Codestral to auto-generate or suggest
  documentation.

* Notes for Maintainers
- Update this file if the project's documentation standards evolve.
- Add project-specific examples or exceptions as needed.
