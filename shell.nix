{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = [
    pkgs.age
    pkgs.httpie
    pkgs.nodejs
  ];
}
