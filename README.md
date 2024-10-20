# bitmap

`bitmap` is an open-source CLI tool for reading, manipulating, and transforming bitmap image files. It provides a set of powerful features for image processing, allowing users to perform various operations on BMP files.

## Features

- **Header Information**: Display detailed metadata about bitmap files.
- **Mirror**: Flip images horizontally or vertically.
- **Filter**: Apply various color filters and effects to images.
- **Rotate**: Rotate images by specified angles.
- **Crop**: Trim images to desired dimensions.
- **Combine**: Apply multiple transformations in a single command.

## Installation

### Prerequisites

- Go 1.16 or higher

### Cloning the Repository

```bash
git clone https://github.com/yourusername/bitmap.git
cd bitmap
```

### Building the Project

To build the project, use the provided Makefile:

```bash
make build
```

This will create an executable named `bitmap` in the current directory.

## Usage

### General Syntax

```bash
./bitmap <command> [options] <input_file> [output_file]
```

### Commands

1. **Header Information**

   Display bitmap file header information:

   ```bash
   ./bitmap header sample.bmp
   ```

2. **Apply Transformations**

   Apply various transformations to an image:

   ```bash
   ./bitmap apply [options] <input_file> <output_file>
   ```

   Options:
   - `--mirror=<direction>`: Mirror the image (horizontal/vertical)
   - `--filter=<type>`: Apply a color filter (blue/red/green/grayscale/negative/pixelate/blur)
   - `--rotate=<angle>`: Rotate the image (right/left/90/180/270/-90/-180/-270)
   - `--crop=<x>-<y>-<width>-<height>`: Crop the image

   Example:
   ```bash
   ./bitmap apply --mirror=horizontal --filter=grayscale --rotate=90 --crop=20-20-100-100 input.bmp output.bmp
   ```

3. **Help**

   Display help information:

   ```bash
   ./bitmap --help
   ./bitmap <command> --help
   ```

## Development

To format the code and run the project during development:

```bash
make run
```

To only format the code:

```bash
make format
```

## License

This project is open source and available under the [MIT License](LICENSE).

## Contributing

Contributions to the bitmap project are welcome! Please feel free to submit a Pull Request.
