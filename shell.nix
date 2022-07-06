{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    age
    flutter
    go
    httpie
    nodejs
    zig
  ];
}
