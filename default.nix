{ buildGoModule, nix-gitignore }:

buildGoModule {
    pname = "ybs";
    version = "0.0.1";
    src = nix-gitignore.gitignoreSource [] ./.;
    goPackagePath = "github.com/esselius/ybs";
    vendorSha256 = "0ddd3dxq4niakzsxx10rpm6lzcrd7cfklswshh1yqi6vld61hyfk";
}