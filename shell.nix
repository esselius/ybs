{ pkgs ? import <nixpkgs> { } }:

let
  ybs = pkgs.callPackage ./default.nix {};
in
  with pkgs; mkShell {
    buildInputs = [
      ybs

      chromedriver
    ];
  }