{ sources ? import ./nix/sources.nix }:

let
    pkgs = import sources.nixpkgs {};
in with pkgs; rec {
    application = pkgs.callPackage ./default.nix {};
    docker = dockerTools.buildImage {
        name = application.name;
        contents = application;
        config = {
            Cmd = "ybs";
        };
    };
}