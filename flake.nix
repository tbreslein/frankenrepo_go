{
  description = "language/toolchain agnostic monorepo tool";
  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    rust-overlay.url = "github:oxalica/rust-overlay";
  };
  outputs = { self, nixpkgs, flake-utils, rust-overlay, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        overlays = [ (import rust-overlay) ];
        #pkgs = nixpkgs.legacyPackages.${system};
        pkgs = import nixpkgs { inherit system overlays; };
        version = "0.0.1";
      in {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gofumpt
            golangci-lint
            go-task

            # for rust test projects
            openssl
            pkg-config
            (rust-bin.stable.latest.default.override {
              extensions = [ "rust-analyzer" ];
            })

            # for c and c++ test projects
            cmake
            gnumake
            ninja
            gcc
          ];
        };

        packages.default = pkgs.buildGoModule {
          pname = "frankenrepo";
          inherit version;
          src = ./.;
          # vendorSha256 = "0000000000000000000000000000000000000000000000000000";
          vendorSha256 = "3tO/+Mnvl/wpS7Ro3XDIVrlYTGVM680mcC15/7ON6qM=";
        };
      });
}
