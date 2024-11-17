# Start with the base coder server image
FROM codercom/code-server:latest

USER root

# Install additional languages and tools as needed
RUN apt-get update && apt-get install -y \
    python3 \
    python3-pip \
    nodejs

USER coder
# You can add more customizations here based on instructor needs

# Set the working directory
WORKDIR /home/coder

# The entry point remains the same as the base image