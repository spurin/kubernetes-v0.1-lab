#!/bin/bash

# Specify the operating systems and architectures
archs=("amd64" "arm" "arm64")
os=("linux" "darwin")

# Create the base bin directory if it doesn't exist
mkdir -p bin

# Loop through each subdirectory in cmd
for dir in $(find cmd -mindepth 1 -maxdepth 1 -type d); do
    base_dir=$(basename $dir)
    echo "Processing directory: $base_dir"

    for opsys in "${os[@]}"; do
        for arch in "${archs[@]}"; do
            # Specific rule: Only allow amd64 and arm64 for Darwin
            if [[ "$opsys" == "darwin" && "$arch" != "amd64" && "$arch" != "arm64" ]]; then
                echo "   - Skipping build for $base_dir on $opsys-$arch as only amd64 and arm64 are supported on Darwin."
                continue
            fi

            # Set environment variables
            export GOOS=$opsys
            export GOARCH=$arch

            # Create architecture-specific directory under bin
            mkdir -p "bin/${opsys}/${arch}"

            # Define output binary name
            output_name="${base_dir}"

            echo " - Building $output_name..."
            cd $dir

            # Build the binary and place it in the appropriate directory
            if go build -o "../../bin/${opsys}/${arch}/$output_name"; then
                echo "   - Successfully built $output_name and stored in 'bin/${opsys}/${arch}' directory."
            else
                echo "   - Failed to build $output_name. Check the Go files and dependencies."
            fi

            # Return to the root directory
            cd - > /dev/null
        done
    done
done

echo "All builds are completed. Binaries are organized by architecture in the 'bin' directory."
