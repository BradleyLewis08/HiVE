#!/bin/bash

# Read student list from a config map or environment variable
STUDENTS="student1 student2 student3"

for student in $STUDENTS
do
    # Create user
    useradd -m -d /home/coder/project/students/$student $student
    
    # Set permissions
    chown -R $student:$student /home/coder/project/students/$student
    chmod 700 /home/coder/project/students/$student

    # Create assignments directory
    mkdir -p /home/coder/project/students/$student/assignments
    chown $student:$student /home/coder/project/students/$student/assignments
    chmod 755 /home/coder/project/students/$student/assignments
done

# Set up shared directories
mkdir -p /home/coder/project/shared/assignments
mkdir -p /home/coder/project/shared/resources
chmod 755 /home/coder/project/shared/assignments
chmod 755 /home/coder/project/shared/resources

# Ensure code-server can access all directories
chown -R coder:coder /home/coder/project
