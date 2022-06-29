{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    age
    flutter
    httpie
    nodejs
    zig
  ];
}
