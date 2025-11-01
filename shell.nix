# shell.nix
{pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
    packages = [
        # DEVELOPMENT
        pkgs.gotest
        pkgs.bun
        pkgs.sqlc
        pkgs.nodejs
        pkgs.prettierd
        pkgs.tailwindcss
        pkgs.pre-commit

        # IDE
        pkgs.gotools
        pkgs.air
        pkgs.htmx-lsp
        pkgs.tailwindcss-language-server
        pkgs.vscode-langservers-extracted

        # DEPLOYMENT
        pkgs.ko
    ];

    LD_LIBRARY_PATH = "${pkgs.lib.makeLibraryPath [
        pkgs.stdenv.cc.cc
    ]}";

    shellHook = ''
        set -a; source .env; set +a
        echo "SHELLHOOK LOG: .env loaded to ENV variables"
    '';
}
