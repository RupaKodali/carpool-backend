#!/bin/bash
go build -o carpool-backend &&
sudo mv carpool-backend /usr/local/bin/ &&
sudo systemctl restart carpool.service &&
echo "🚀 Deployed and restarted"