#!/usr/bin/env bash
V=$(cat ./VERSION)
case "$(uname -s)" in

   Darwin)
     OS=secrets-darwin-amd64
     ;;

   Linux)
     OS=secrets-linux-amd64
     ;;

   CYGWIN*|MINGW32*|MSYS*)
     OS=secrets-windows-amd64.exe
     ;;
   *)
     echo 'OS not supported' 
     ;;
esac
mv ./target/${V}/${OS} /usr/local/bin/secrets
