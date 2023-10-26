{
  description = "language/toolchain agnostic monorepo tool";
  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };
  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        buildDeps = with pkgs; [ go ];
      in {
        devShells.default = pkgs.mkShell {
          nativeBuildInputs = buildDeps
            ++ [ pkgs.gopls pkgs.gofumpt pkgs.golangci-lint pkgs.go-task ];
        };

        packages.default = pkgs.stdenv.mkDerivation {
          nativeBuildInputs = buildDeps;
          name = "frankenrepo";
          src = ./.;
          buildPhase = "go build";
          installPhase = "cp ./frankenrepo $out/";
        };
      });
}
