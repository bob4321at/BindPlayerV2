{
  description = "BindPlayerV2 flake";

  inputs = { nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable"; };

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
      buildInputs = [
        pkgs.pkgconf
        pkgs.gcc
        pkgs.xorg.libX11
        pkgs.xorg.libXrandr
        pkgs.xorg.libXinerama
        pkgs.xorg.libXcursor
        pkgs.xorg.libXi
        pkgs.mesa
        pkgs.libGL
        pkgs.libGLU
        pkgs.xorg.libXxf86vm
      ];
    in {
      packages.${system}.default = pkgs.buildGoModule {
        pname = "BindPlayerV2";
        version = "1.0.0";
        src = ./.;
        vendorHash = null;
        doCheck = false;

        inherit buildInputs;
      };

      devShells.${system}.default = pkgs.mkShell {
        inherit buildInputs;
        nativeBuildInputs = [ pkgs.go ];
        
        shellHook = ''
          export LD_LIBRARY_PATH=${pkgs.lib.makeLibraryPath buildInputs}:$LD_LIBRARY_PATH
        '';
      };
    };
}


