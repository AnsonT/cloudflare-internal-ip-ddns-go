# Cloudflare Internal IP DDNS

This project provides a simple command-line tool to update a Cloudflare DNS record with the internal IP address of the machine it's run on. It's particularly useful for setting up a dynamic DNS (DDNS) for home networks or any scenario where the internal IP address may change frequently but needs to be accessible via a consistent domain name.

## Features

- Automatically detects the internal IP address of the machine.
- Updates an existing A record in Cloudflare or creates a new one if it doesn't exist.
- Configurable through command-line flags or environment variables.

## Requirements

- A Cloudflare account and API token with permissions to edit DNS records.
- Go 1.15 or later.

## Installation

Clone the repository and build the binary:

```bash
git clone https://github.com/yourusername/cloudflare-internal-ip-ddns.git
cd cloudflare-internal-ip-ddns
go build -o cloudflare-internal-ip-ddns
```

## Usage

The tool can be run directly from the command line. You can specify the Cloudflare API token, zone name (domain), and record name (subdomain) either through flags or environment variables.

### Command-Line Flags

- `--token` or `-t`: Cloudflare API token (required)
- `--domain` or `-d`: Top-level domain name (Zone name), e.g., 'example.com' (required)
- `--subdomain` or `-s`: Subdomain name (Record name), e.g., 'home' for 'home.example.com', default to '*'

### Environment Variables

Alternatively, you can set the following environment variables:

- `CF_API_TOKEN`: Cloudflare API token
- `CF_ZONE_NAME`: Top-level domain name (Zone name)
- `CF_RECORD_NAME`: Subdomain name (Record name)

### Running the Tool

```bash
./cloudflare-internal-ip-ddns --token your_cloudflare_api_token --domain example.com --subdomain home
```

Or using environment variables:

```bash
export CF_API_TOKEN=your_cloudflare_api_token
export CF_ZONE_NAME=example.com
export CF_RECORD_NAME=home
./cloudflare-internal-ip-ddns
```

Or using .env file

## Contributing

Contributions are welcome! Please feel free to submit a pull request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
```