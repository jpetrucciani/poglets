{ jacobi ? import
    (fetchTarball {
      name = "jpetrucciani-2022-10-31";
      url = "https://nix.cobi.dev/x/02c19fb4ae64983ec48fd8c536c178ace7270549";
      sha256 = "0gn78c9kvwy6dismcwix829ch98dvhxbci3rgpv0kdkjpax9n290";
    })
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
        export DOCKER_HOST=$(${jacobi.docker-client}/bin/docker context inspect --format '{{.Endpoints.docker.Host}}')
        ${jacobi.act}/bin/act --container-architecture linux/amd64 -r --rm
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
      vendorSha256 = "sha256-Hjdv2Fvl1S52CDs4TAR3Yt9pEFUIvs5N5sVhZY+Edzo=";

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
