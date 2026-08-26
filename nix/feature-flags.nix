{
  lib,
  enabled ? { },
}:
let
  inherit (lib) fileset;

  root = ../src/features;

  requested = lib.filterAttrs (_: value: value != false) enabled;

  featureDir =
    name:
    let
      dir = root + "/${name}";
    in
    lib.throwIfNot (builtins.pathExists dir) "feature-flags: no feature directory at ${toString dir}"
      dir;

  sources = lib.mapAttrsToList (name: value: if value == true then featureDir name else value) (
    lib.filterAttrs (_: value: !lib.isDerivation value) requested
  );
in
{
  # Repo-local features, unioned into one fileset rooted in the source tree.
  fileset = fileset.unions sources;

  # Features built outside the tree; a fileset cannot span the store and the
  # repo, so the caller places these itself.
  drvs = lib.filterAttrs (_: lib.isDerivation) requested;
}
