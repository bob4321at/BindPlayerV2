{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShell {
  buildInputs = with pkgs; [
    go
    gcc
    xorg.libXrandr
    xorg.libXinerama
    xorg.libXcursor
    xorg.libXi
    xorg.libXxf86vm
    mesa
    libglvnd
    alsa-lib
  ];
  nativeBuildInputs = with pkgs; [
    pkg-config
  ];

  shellHook = ''
     export LD_LIBRARY_PATH=${pkgs.mesa}/lib:${pkgs.libglvnd}/lib:$LD_LIBRARY_PATH
  '';

    CGO_ENABLED = "1";
}

