{ pkgs ? import <nixpkgs> { } }:
with pkgs;
mkShell {
  name = "polytask-shell";
  buildInputs = [
    go_1_17
    go-rice
    nixpkgs-fmt
  ];
}
