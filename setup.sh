#!/bin/bash

SESH="core"

tmux has-session -t $SESH 2>/dev/null

if [ $? != 0 ]; then
	tmux new-session -d -s $SESH -n "nvim"

	tmux send-keys -t $SESH:nvim "cd ~/core/" C-m
	tmux send-keys -t $SESH:nvim "nvim" C-m

	tmux new-window -t $SESH -n "sh"
	tmux send-key -t $SESH:sh "cd ~/core/" C-m

	tmux new-window -t $SESH -n "docker"
	tmux send-key -t $SESH:docker "cd ~/core/" C-m
	tmux send-key -t $SESH:docker "docker compose up --no-attach cloudflared" C-m

	tmux select-window -t $SESH:nvim
fi

tmux attach-session -t $SESH
