{
  description = "Go Flake";
  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, utils, }:
    utils.lib.eachDefaultSystem (system: rec {
      pkgs = import nixpkgs { inherit system; };
      devShells.default =
        pkgs.mkShell { buildInputs = with pkgs; [ go_1_18 golangci-lint ]; };
    });
}
