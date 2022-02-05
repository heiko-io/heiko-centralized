with import <nixpkgs> {};
mkShell {
  packages = [
    rustc
    cargo
  ];
  RUST_SRC_PATH = "${pkgs.rust.packages.stable.rustPlatform.rustLibSrc}";
}
