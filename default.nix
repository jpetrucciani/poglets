{ jacobi ? import
    (
      fetchTarball {
        name = "jpetrucciani-2022-09-20";
        url = "https://github.com/jpetrucciani/nix/archive/845e6c83a608145991297f0a2ff0dd959f1679a6.tar.gz";
        sha256 = "1840v62z68dnfvlcjs90v0f3xg8bdps7pz8bkwxa1gi87rdrsjv2";
      }
    )
    { }
}:
let
  name = "poglets";
  tools = with jacobi; {
    cli = [
      nixpkgs-fmt
    ];
    go = [
      go_1_19
      go-tools
    ];
    scripts = [
      (writeShellScriptBin "test_actions" ''
        ${jacobi.act}/bin/act --artifact-server-path ./.cache/ -r --rm
      '')
    ];
  };

  env = jacobi.enviro {
    inherit name tools;
  };

  poglets =
    let
      version = "0.0.1";
      commit = "";
    in
    jacobi.buildGo119Module rec {
      pname = name;
      version = "0.0.1";
      src = ./.;
      vendorSha256 = "sha256-Pra/yuEnkrsJNgEYpIX7fsf1vhQKFXTSb9bHVu8/3UM=";

      ldflags = [
        "-s"
        "-w"
        "-X main.Version=${version}"
        "-X main.GitCommit=${commit}"
      ];

      meta = with jacobi.lib; {
        description = "";
      };
    };

in
env // { inherit poglets; }
