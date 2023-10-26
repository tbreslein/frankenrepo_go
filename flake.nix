{
  description = "language/toolchain agnostic monorepo tool";
  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };
  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = nixpkgs.legacyPackages.${system};
      in {
        devShells.default =
          pkgs.mkShell { nativeBuildInputs = with pkgs; [ go ]; };

        packages.default = pkgs.stdenv.mkDerivation {
          nativeBuildInputs = with pkgs; [ go ];
          name = "frankenrepo";
          src = ./.;

          buildPhase = ''
            go build
          '';

          installPhase = ''
            cp ./frankenrepo $out/
          '';
        };
      });
}
