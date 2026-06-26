{
  description = "Gososialize development shell (nix shell)";

  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    go-overlay.url = "github:purpleclay/go-overlay";
  };

  outputs = inputs:
    inputs.flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import (inputs.nixpkgs) {
          inherit system;
          overlays = [inputs.go-overlay.overlays.default];
        };
      in {
        devShell = pkgs.mkShell {
          buildInputs = with pkgs; [
            # golang
            go-bin.versions."1.26.1"
            air
            sqlc
            goose
            gotestfmt

            # node
            nodejs
            pnpm
            eslint_d
            prettierd
            tailwindcss-language-server
            vscode-langservers-extracted

            # utils
            jq
            bruno
            lazygit
            resterm
          ];

          shellHook = ''
            export GOBIN=$HOME/go/bin
            export PATH=$GOBIN:$PATH

            echo "Gososialize dev shell activated! Happy coding 🚀."
          '';
        };
      }
    );
}
