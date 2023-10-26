{
  description = "language/toolchain agnostic monorepo tool";
  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };
  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        version = "0.0.1";
      in {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [ go gopls gofumpt golangci-lint go-task ];
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
