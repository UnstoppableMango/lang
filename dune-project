(lang dune 3.20)

(name mangolang)

(generate_opam_files true)

(source
 (github UnstoppableMango/lang))

(authors "Author Name <erik.rasmussen@unmango.dev>")

(maintainers "Maintainer Name <erik.rasmussen@unmango.dev>")

(license LICENSE)

(documentation https://github.com/UnstoppableMango/lang)

(package
 (name mangolang)
 (synopsis "Fiddling with language design")
 (description "A custom language project for learning and experimentation")
 (depends ocaml)
 (tags
  ("ocaml" "lang" "compiler" "interpreter" "learning")))

; See the complete stanza docs at https://dune.readthedocs.io/en/stable/reference/dune-project/index.html
