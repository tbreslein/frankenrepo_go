{
  description = "language/toolchain agnostic monorepo tool";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    rust-overlay.url = "github:oxalica/rust-overlay";
    flake-parts.url = "github:hercules-ci/flake-parts";
  };

  outputs = inputs@{ self, ... }:
    inputs.flake-parts.lib.mkFlake { inherit inputs; } {
      systems =
        [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];

      perSystem = { system, pkgs, ... }:
        let version = "0.0.1";
        in {
          _module.args.pkgs = import inputs.nixpkgs {
            inherit system;
            overlays = [ (import inputs.rust-overlay) ];
          };
          devShells.default = pkgs.mkShell {
            buildInputs = with pkgs; [
              go
              gopls
              gofumpt
              golangci-lint
              just

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
        };
    };
}
