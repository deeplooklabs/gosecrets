# gosecrets

**gosecrets** is a Go-based tool that assists in identifying sensitive keys and confidential information that may be exposed in code repositories, files, and directories. This tool is useful for security researchers and developers who want to ensure that sensitive information is not publicly accessible.

## Features

- Search for sensitive keys in directories and files.
- Support for parallel routines to speed up the search.
- Extensibility to add new rules and search patterns.

## Installation

To install **gosecrets**, you can use the following command:

```shell
go install github.com/deeplooklabs/gosecrets
```

## Usage

1.Navigate to the directory you want to check for key exposure and sensitive information.

2.Execute the following command:
```shell
gosecrets
```

The tool will start searching in the current directory and its subdirectories for patterns and keys that may pose a security risk.

## Contributing

Contributions are welcome! If you want to add new rules, improve the code, or report issues, feel free to open an issue or submit a pull request.

## Disclaimer
This script is intended for educational and security research purposes only. Use it responsibly and always respect local laws and regulations when conducting scans on systems and applications that are not your own.

## License
This project is licensed under the MIT License. See the LICENSE file for details.