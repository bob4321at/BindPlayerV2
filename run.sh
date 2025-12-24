#!/usr/bin/env bash
cd $HOME/golang/BindPlayerV2
nix develop --command bash -c "./bind_player_ui" 2> >(grep -v 'Git tree.*dirty' >&2)

