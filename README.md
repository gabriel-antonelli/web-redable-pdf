# Web Readable PDF

Generate web-readable PDFs from multiple web pages using Go.

## Overview

`web-readable-pdf` is a command-line tool written in Go that enables the creation of web-readable PDFs from multiple web pages. It leverages Go and various libraries to convert HTML content into PDF files.

## Features

- Combine multiple HTML files into a single PDF.
- Create separate PDFs for each web page.

## Dependencies

The main dependencies for `web-readable-pdf` are:

- [go-wkhtmltopdf](https://github.com/SebastiaanKlippert/go-wkhtmltopdf) v1.9.2: A Go wrapper for the `wkhtmltopdf` command line tool, which is used for converting HTML to PDF.

- [go-readability](https://github.com/go-shiori/go-readability) v0.0.0-20231029095239-6b97d5aba789: A library for extracting article content from HTML pages, used to improve the readability of the input HTML.

- [cobra](https://github.com/spf13/cobra) v1.8.0: A popular Go library for creating powerful command-line applications.

## Installation

Ensure that you have Go installed on your machine. Install `web-readable-pdf` using the following command:

go get -u github.com/your-username/web-readable-pdf

## Usage

### Basic Usage

web-readable-pdf [web-page1] [web-page2] [web-page3] [output.pdf]

### Options

- `-o, --output`: Specify the output PDF file. (Default: article title)
- `-s, --separate`: Create separate PDFs for each web page.

### Examples

1. Combine multiple web pages into a single PDF:

web-readable-pdf https://example.com/page1 https://example.com/page2 

2. Create separate PDFs for each web page:

web-readable-pdf -s https://example.com/page1 https://example.com/page2

## Contributing

Contributions are welcome! Please follow our [contribution guidelines](CONTRIBUTING.md).

## License

This project is licensed under the GNU GENERAL PUBLIC LICENSE - see the [LICENSE](LICENSE) file for details.

Make sure to customize the placeholders like `[web-page1]`, `[web-page2]`, `your-username`, and update the example URLs based on your specific use case. Feel free to add more sections or details as needed.
